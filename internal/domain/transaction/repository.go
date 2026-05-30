package domain_transaction

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"ExpenseTracker-Backend/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type Repository interface {
	Create(ctx context.Context, tx *Transaction) error
	GetByID(ctx context.Context, id, userID int64) (*Transaction, error)
	Update(ctx context.Context, tx *Transaction) error
	UpdateStatus(ctx context.Context, id, userID int64, status int) error
	List(ctx context.Context, userID int64, req TransactionListRequest) ([]Transaction, error)
	GetInfo(ctx context.Context, userID int64, req TransactionInfoRequest) (*TransactionInfoResponse, error)
	GetAnalytics(ctx context.Context, userID int64) (*TransactionAnalyticsResponse, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, tx *Transaction) error {
	query := `
		INSERT INTO transactions (user_id, amount, category, type, status, description, date)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		tx.UserID, tx.Amount, tx.Category, tx.Type, tx.Status, tx.Description, tx.Date,
	).Scan(&tx.ID, &tx.CreatedAt, &tx.UpdatedAt)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id, userID int64) (*Transaction, error) {
	query := `
		SELECT id, user_id, amount, category, type, status, description, date, created_at, updated_at
		FROM transactions
		WHERE id = $1 AND user_id = $2
	`
	var tx Transaction
	err := r.db.QueryRow(ctx, query, id, userID).Scan(
		&tx.ID, &tx.UserID, &tx.Amount, &tx.Category, &tx.Type, &tx.Status, &tx.Description, &tx.Date, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("transaction not found or access denied")
		}
		return nil, err
	}
	return &tx, nil
}

func (r *postgresRepository) Update(ctx context.Context, tx *Transaction) error {
	query := `
		UPDATE transactions
		SET amount = $1, category = $2, type = $3, status = $4, description = $5, date = $6, updated_at = NOW()
		WHERE id = $7 AND user_id = $8
		RETURNING updated_at
	`
	err := r.db.QueryRow(ctx, query,
		tx.Amount, tx.Category, tx.Type, tx.Status, tx.Description, tx.Date, tx.ID, tx.UserID,
	).Scan(&tx.UpdatedAt)
	return err
}

func (r *postgresRepository) UpdateStatus(ctx context.Context, id, userID int64, status int) error {
	query := `
		UPDATE transactions
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND user_id = $3
	`
	_, err := r.db.Exec(ctx, query, status, id, userID)
	return err
}

func (r *postgresRepository) List(ctx context.Context, userID int64, req TransactionListRequest) ([]Transaction, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = $1")
	args = append(args, userID)

	paramCount := 2

	if req.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, *req.Status)
		paramCount++
	} else {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, int(utils.Active))
		paramCount++
	}

	if len(req.CategoryIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("category = ANY($%d)", paramCount))
		args = append(args, req.CategoryIDs)
		paramCount++
	} else if req.Category != nil {
		conditions = append(conditions, fmt.Sprintf("category = $%d", paramCount))
		args = append(args, *req.Category)
		paramCount++
	}

	if req.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", paramCount))
		args = append(args, *req.Type)
		paramCount++
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		matchingCatIDs := utils.GetCategoryIDsByNameSearch(search)
		if len(matchingCatIDs) > 0 {
			conditions = append(conditions, fmt.Sprintf("(description ILIKE $%d OR category = ANY($%d))", paramCount, paramCount+1))
			args = append(args, "%"+search+"%", matchingCatIDs)
			paramCount += 2
		} else {
			conditions = append(conditions, fmt.Sprintf("description ILIKE $%d", paramCount))
			args = append(args, "%"+search+"%")
			paramCount++
		}
	}

	if req.FromDate != "" {
		if t, err := time.Parse(time.RFC3339, req.FromDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("date >= $%d", paramCount))
			args = append(args, t)
			paramCount++
		}
	}

	if req.ToDate != "" {
		if t, err := time.Parse(time.RFC3339, req.ToDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("date <= $%d", paramCount))
			args = append(args, t)
			paramCount++
		}
	}

	// Dynamic cursor logic
	var orderClause string
	if req.AfterID > 0 {
		var cursorAmount decimal.Decimal
		var cursorErr error

		// Determine sorting cursor conditions
		switch req.Sort {
		case "amount_asc":
			cursorErr = r.db.QueryRow(ctx, "SELECT amount FROM transactions WHERE id = $1", req.AfterID).Scan(&cursorAmount)
			if cursorErr == nil {
				conditions = append(conditions, fmt.Sprintf("(amount > $%d OR (amount = $%d AND id > $%d))", paramCount, paramCount, paramCount+1))
				args = append(args, cursorAmount, req.AfterID)
				paramCount += 2
			}
		case "amount_desc":
			cursorErr = r.db.QueryRow(ctx, "SELECT amount FROM transactions WHERE id = $1", req.AfterID).Scan(&cursorAmount)
			if cursorErr == nil {
				conditions = append(conditions, fmt.Sprintf("(amount < $%d OR (amount = $%d AND id < $%d))", paramCount, paramCount, paramCount+1))
				args = append(args, cursorAmount, req.AfterID)
				paramCount += 2
			}
		default: // default/newest (id desc)
			conditions = append(conditions, fmt.Sprintf("id < $%d", paramCount))
			args = append(args, req.AfterID)
			paramCount++
		}
	}

	switch req.Sort {
	case "amount_asc":
		orderClause = "ORDER BY amount ASC, id ASC"
	case "amount_desc":
		orderClause = "ORDER BY amount DESC, id DESC"
	default:
		orderClause = "ORDER BY id DESC"
	}

	limit := req.Limit
	if limit <= 0 {
		limit = 20
	} else if limit > 100 {
		limit = 100
	}

	query := fmt.Sprintf(`
		SELECT id, user_id, amount, category, type, status, description, date, created_at, updated_at
		FROM transactions
		WHERE %s
		%s
		LIMIT $%d
	`, strings.Join(conditions, " AND "), orderClause, paramCount)

	args = append(args, limit)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var tx Transaction
		err = rows.Scan(
			&tx.ID, &tx.UserID, &tx.Amount, &tx.Category, &tx.Type, &tx.Status, &tx.Description, &tx.Date, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *postgresRepository) GetInfo(ctx context.Context, userID int64, req TransactionInfoRequest) (*TransactionInfoResponse, error) {
	var conditions []string
	var args []interface{}

	conditions = append(conditions, "user_id = $1")
	args = append(args, userID)

	paramCount := 2

	if req.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, *req.Status)
		paramCount++
	} else {
		conditions = append(conditions, fmt.Sprintf("status = $%d", paramCount))
		args = append(args, int(utils.Active))
		paramCount++
	}

	if len(req.CategoryIDs) > 0 {
		conditions = append(conditions, fmt.Sprintf("category = ANY($%d)", paramCount))
		args = append(args, req.CategoryIDs)
		paramCount++
	} else if req.Category != nil {
		conditions = append(conditions, fmt.Sprintf("category = $%d", paramCount))
		args = append(args, *req.Category)
		paramCount++
	}

	if req.Type != nil {
		conditions = append(conditions, fmt.Sprintf("type = $%d", paramCount))
		args = append(args, *req.Type)
		paramCount++
	}

	search := strings.TrimSpace(req.Search)
	if search != "" {
		matchingCatIDs := utils.GetCategoryIDsByNameSearch(search)
		if len(matchingCatIDs) > 0 {
			conditions = append(conditions, fmt.Sprintf("(description ILIKE $%d OR category = ANY($%d))", paramCount, paramCount+1))
			args = append(args, "%"+search+"%", matchingCatIDs)
			paramCount += 2
		} else {
			conditions = append(conditions, fmt.Sprintf("description ILIKE $%d", paramCount))
			args = append(args, "%"+search+"%")
			paramCount++
		}
	}

	if req.FromDate != "" {
		if t, err := time.Parse(time.RFC3339, req.FromDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("date >= $%d", paramCount))
			args = append(args, t)
			paramCount++
		}
	}

	if req.ToDate != "" {
		if t, err := time.Parse(time.RFC3339, req.ToDate); err == nil {
			conditions = append(conditions, fmt.Sprintf("date <= $%d", paramCount))
			args = append(args, t)
			paramCount++
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	// Query 1: Total stats
	statsQuery := fmt.Sprintf(`
		SELECT 
			COUNT(*) AS total_transactions,
			COALESCE(SUM(CASE WHEN type = 1 THEN amount ELSE 0 END), 0) AS total_expense,
			COALESCE(SUM(CASE WHEN type = 2 THEN amount ELSE 0 END), 0) AS total_income
		FROM transactions
		WHERE %s
	`, whereClause)
	var stats TransactionInfoResponse
	var totalExpenseDec, totalIncomeDec decimal.Decimal
	err := r.db.QueryRow(ctx, statsQuery, args...).Scan(
		&stats.TotalTransactions, &totalExpenseDec, &totalIncomeDec,
	)
	if err != nil {
		return nil, err
	}
	stats.TotalExpense, _ = totalExpenseDec.Float64()
	stats.TotalIncome, _ = totalIncomeDec.Float64()

	// Query 2: Category breakdown using SQL Aggregation
	var catConditions = append([]string{}, conditions...)
	if req.Type == nil {
		catConditions = append(catConditions, "type = 1")
	}

	catWhereClause := strings.Join(catConditions, " AND ")
	categoryQuery := fmt.Sprintf(`
		SELECT 
			category,
			COUNT(*) AS count,
			COALESCE(SUM(amount), 0) AS amount
		FROM transactions
		WHERE %s
		GROUP BY category
		ORDER BY amount DESC
	`, catWhereClause)
	rows, err := r.db.Query(ctx, categoryQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var summaries []CategorySummary
	for rows.Next() {
		var cs CategorySummary
		var categoryAmount decimal.Decimal
		err = rows.Scan(&cs.Category, &cs.Count, &categoryAmount)
		if err != nil {
			return nil, err
		}
		cs.Amount, _ = categoryAmount.Float64()
		summaries = append(summaries, cs)
	}

	stats.CategorySummary = summaries
	if stats.CategorySummary == nil {
		stats.CategorySummary = []CategorySummary{}
	}
	return &stats, nil
}

func (r *postgresRepository) GetAnalytics(ctx context.Context, userID int64) (*TransactionAnalyticsResponse, error) {
	now := time.Now()

	var res TransactionAnalyticsResponse
	var errs [3]error
	var wg sync.WaitGroup
	wg.Add(3)

	// 1. Weekly
	go func() {
		defer wg.Done()
		offset := int(now.Weekday()) - 1
		if offset < 0 {
			offset = 6
		}
		startOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset)
		endOfWeek := startOfWeek.AddDate(0, 0, 7)

		labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
		expense := make([]float64, 7)
		income := make([]float64, 7)

		query := `
			SELECT 
				EXTRACT(isodow FROM date)::int as day_num,
				COALESCE(SUM(CASE WHEN type = 1 THEN amount ELSE 0 END), 0)::numeric as expense,
				COALESCE(SUM(CASE WHEN type = 2 THEN amount ELSE 0 END), 0)::numeric as income
			FROM transactions
			WHERE user_id = $1 AND status = 1 AND date >= $2 AND date < $3
			GROUP BY day_num
		`
		rows, err := r.db.Query(ctx, query, userID, startOfWeek, endOfWeek)
		if err != nil {
			errs[0] = err
			return
		}
		defer rows.Close()

		var totalExp, totalInc float64
		for rows.Next() {
			var dayNum int
			var expDecimal, incDecimal decimal.Decimal
			if err := rows.Scan(&dayNum, &expDecimal, &incDecimal); err != nil {
				errs[0] = err
				return
			}
			exp, _ := expDecimal.Float64()
			inc, _ := incDecimal.Float64()
			if dayNum >= 1 && dayNum <= 7 {
				expense[dayNum-1] = exp
				income[dayNum-1] = inc
				totalExp += exp
				totalInc += inc
			}
		}

		res.Weekly = AnalyticsData{
			Labels:       labels,
			Expense:      expense,
			Income:       income,
			TotalExpense: totalExp,
			TotalIncome:  totalInc,
			Net:          totalInc - totalExp,
		}
	}()

	// 2. Monthly
	go func() {
		defer wg.Done()
		startMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, -11, 0)
		endMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).AddDate(0, 1, 0)

		labels := make([]string, 12)
		expense := make([]float64, 12)
		income := make([]float64, 12)

		for i := 0; i < 12; i++ {
			m := startMonth.AddDate(0, i, 0)
			labels[i] = m.Format("Jan-06")
		}

		query := `
			SELECT 
				DATE_TRUNC('month', date) as month_start,
				COALESCE(SUM(CASE WHEN type = 1 THEN amount ELSE 0 END), 0)::numeric as expense,
				COALESCE(SUM(CASE WHEN type = 2 THEN amount ELSE 0 END), 0)::numeric as income
			FROM transactions
			WHERE user_id = $1 AND status = 1 AND date >= $2 AND date < $3
			GROUP BY month_start
		`
		rows, err := r.db.Query(ctx, query, userID, startMonth, endMonth)
		if err != nil {
			errs[1] = err
			return
		}
		defer rows.Close()

		var totalExp, totalInc float64
		for rows.Next() {
			var monthStart time.Time
			var expDecimal, incDecimal decimal.Decimal
			if err := rows.Scan(&monthStart, &expDecimal, &incDecimal); err != nil {
				errs[1] = err
				return
			}
			exp, _ := expDecimal.Float64()
			inc, _ := incDecimal.Float64()
			monthsDiff := (monthStart.Year() - startMonth.Year()) * 12 + int(monthStart.Month() - startMonth.Month())
			if monthsDiff >= 0 && monthsDiff < 12 {
				expense[monthsDiff] = exp
				income[monthsDiff] = inc
				totalExp += exp
				totalInc += inc
			}
		}

		res.Monthly = AnalyticsData{
			Labels:       labels,
			Expense:      expense,
			Income:       income,
			TotalExpense: totalExp,
			TotalIncome:  totalInc,
			Net:          totalInc - totalExp,
		}
	}()

	// 3. Yearly
	go func() {
		defer wg.Done()
		startYear := time.Date(now.Year()-4, 1, 1, 0, 0, 0, 0, now.Location())
		endYear := time.Date(now.Year()+1, 1, 1, 0, 0, 0, 0, now.Location())

		labels := make([]string, 5)
		expense := make([]float64, 5)
		income := make([]float64, 5)

		for i := 0; i < 5; i++ {
			labels[i] = strconv.Itoa(startYear.Year() + i)
		}

		query := `
			SELECT 
				EXTRACT(year FROM date)::int as yr,
				COALESCE(SUM(CASE WHEN type = 1 THEN amount ELSE 0 END), 0)::numeric as expense,
				COALESCE(SUM(CASE WHEN type = 2 THEN amount ELSE 0 END), 0)::numeric as income
			FROM transactions
			WHERE user_id = $1 AND status = 1 AND date >= $2 AND date < $3
			GROUP BY yr
		`
		rows, err := r.db.Query(ctx, query, userID, startYear, endYear)
		if err != nil {
			errs[2] = err
			return
		}
		defer rows.Close()

		var totalExp, totalInc float64
		for rows.Next() {
			var yr int
			var expDecimal, incDecimal decimal.Decimal
			if err := rows.Scan(&yr, &expDecimal, &incDecimal); err != nil {
				errs[2] = err
				return
			}
			exp, _ := expDecimal.Float64()
			inc, _ := incDecimal.Float64()
			yearDiff := yr - startYear.Year()
			if yearDiff >= 0 && yearDiff < 5 {
				expense[yearDiff] = exp
				income[yearDiff] = inc
				totalExp += exp
				totalInc += inc
			}
		}

		res.Yearly = AnalyticsData{
			Labels:       labels,
			Expense:      expense,
			Income:       income,
			TotalExpense: totalExp,
			TotalIncome:  totalInc,
			Net:          totalInc - totalExp,
		}
	}()

	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, err
		}
	}

	return &res, nil
}