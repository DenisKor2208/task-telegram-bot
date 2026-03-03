// Package dbutils Хелпер-обёртка для выполнения запросов на базе sqlx и для функций подключения к БД (pgx).
package dbutils

// Хелпер-обёртка для выполнения запросов на базе sqlx

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"go.uber.org/multierr"
)

// sqlErr Форматирование текстов ошибок.
// - err: исходная ошибка.
// - query: SQL-запрос, который вызвал ошибку.
// - args: аргументы, переданные в запрос.
func sqlErr(err error, query string, args ...any) error {
	return fmt.Errorf(`run query "%s" with args %+v: %w`, query, args, err)
}

// namedQuery Заполнение запросов именованными параметрами.
// - query: SQL-запрос с именованными параметрами.
// - arg: структура или map с значениями для именованных параметров.
// Возвращает: nq - запрос с заменёнными параметрами, args - срез аргументов, err - ошибка.
func namedQuery(query string, arg any) (nq string, args []any, err error) {
	nq, args, err = sqlx.Named(query, arg)
	if err != nil {
		return "", nil, sqlErr(err, query, args...)
	}
	return nq, args, nil
}

// Exec Выполнение запросов с параметрами (неименованные, в виде \$1...$n).
// - ctx: контекст выполнения.
// - db: интерфейс для выполнения запросов.
// - query: SQL-запрос.
// - args: аргументы для запроса.
// Параметры передаются в виде среза args в порядке плейсхолдеров.
// Подходит для простых запросов с фиксированным порядком аргументов.
// Возвращает: результат выполнения и ошибка.
func Exec(ctx context.Context, db sqlx.ExtContext, query string, args ...any) (sql.Result, error) {
	res, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return res, sqlErr(err, query, args...)
	}

	return res, nil
}

// NamedExec выполняет SQL-запросы с именованными параметрами (например, :name, :age).
// Параметры передаются в виде структуры или map, автоматически преобразуются в позиционные.
// Удобен для сложных запросов с множеством параметров, улучшает читаемость.
// - ctx: контекст выполнения.
// - db: интерфейс для выполнения запросов с поддержкой именованных параметров.
// - query: SQL-запрос с именованными параметрами.
// - arg: структура или map с значениями для именованных параметров.
// Возвращает: результат выполнения и ошибка.
func NamedExec(ctx context.Context, db sqlx.ExtContext, query string, arg any) (sql.Result, error) {
	nq, args, err := namedQuery(query, arg)
	if err != nil {
		return nil, err
	}

	return Exec(ctx, db, db.Rebind(nq), args...)
}

// Select Выборка по запросу с параметрами (неименованные, в виде \$1...$n).
// - ctx: контекст выполнения.
// - db: интерфейс для выполнения запросов.
// - dest: указатель на срез или структуру, куда будут записаны результаты.
// - query: SQL-запрос.
// - args: аргументы для запроса.
//
// Результаты сканируются в предоставленный dest (срез структур, срез map или структура).
// Подходит для множественных строк; если строк нет, dest остаётся пустым.
//
// Возвращает: ошибка выполнения.
func Select(ctx context.Context, db sqlx.ExtContext, dest any, query string, args ...any) error {
	if err := sqlx.SelectContext(ctx, db, dest, query, args...); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// GetMap Выборка по запросу с параметрами (неименованные, в виде \$1...$n).
// Возвращаемое значение - map - map[string]any.
// - ctx: контекст выполнения.
// - db: интерфейс для выполнения запросов.
// - query: SQL-запрос.
// - args: аргументы для запроса.
//
// Возвращает одну строку в виде map[string]any (ключи — имена столбцов).
// Подходит для единственной строки; если строк больше одной, возвращает ошибку.
//
// Возвращает: map с результатами и ошибка.
func GetMap(ctx context.Context, db sqlx.ExtContext, query string, args ...any) (ret map[string]any, err error) {
	row := db.QueryRowxContext(ctx, query, args...)
	if row.Err() != nil {
		return nil, sqlErr(row.Err(), query, args...)
	}

	ret = map[string]any{}
	if err := row.MapScan(ret); err != nil {
		return nil, sqlErr(err, query, args...)
	}

	return ret, nil
}

// Get выполняет запрос и сканирует результат в структуру dest.
// dest должен быть указателем на структуру с тегами db, соответствующими столбцам запроса.
// Возвращает ошибку, если запрос неудачен или строка не найдена (sql.ErrNoRows).
func Get(ctx context.Context, db sqlx.ExtContext, dest interface{}, query string, args ...any) error {
	row := db.QueryRowxContext(ctx, query, args...)
	if row.Err() != nil {
		return sqlErr(row.Err(), query, args...)
	}

	if err := row.StructScan(dest); err != nil {
		return sqlErr(err, query, args...)
	}

	return nil
}

// TxFunc Описание типа вложенной функции для выполнения в транзакции.
// - tx: транзакция sqlx, в которой выполняется функция.
// Возвращает: ошибка, если требуется откат транзакции.
type TxFunc func(tx *sqlx.Tx) error

// TxRunner Интерфейс для запуска транзакции (sqlx).
type TxRunner interface {
	BeginTxx(context.Context, *sql.TxOptions) (*sqlx.Tx, error)
}

// RunTx
//
// Запуск транзакции (в случае ошибки выполнения вложенной функции вызовет откат транзакции).
// Вложенная функция (f TxFunc) должна возвращать ошибку в случае присутствия условий, требущих откат транзакции.
// - ctx: контекст выполнения.
// - db: интерфейс для запуска транзакции.
// - f: функция, выполняемая в транзакции.
// Возвращает: ошибка выполнения.
func RunTx(ctx context.Context, db TxRunner, f TxFunc) (err error) {
	var tx *sqlx.Tx // Объект транзакции sqlx, через который выполняются запросы.

	// Настраиваем параметры транзакции.
	opts := &sql.TxOptions{
		Isolation: sql.LevelReadCommitted, // Уровень изоляции: видим только подтверждённые изменения от других транзакций.
	}

	// Начинаем новую транзакцию.
	tx, err = db.BeginTxx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err) // Если не удалось начать, возвращаем ошибку с деталями.
	}

	// Это отложенный блок: выполнится в конце функции.
	// Он гарантирует, что транзакция завершится (коммит или откат), даже при панике.
	defer func() {
		if err != nil {
			// Есть ошибка — откатываем все изменения в транзакции.
			err = multierr.Combine(err, tx.Rollback())
		} else {
			// Ошибок нет — сохраняем изменения (коммит).
			err = tx.Commit()
		}
	}()

	// Выполняем вашу функцию внутри транзакции и возвращаем её результат.
	return f(tx)
}
