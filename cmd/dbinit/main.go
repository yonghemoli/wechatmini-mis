package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	withSeed := flag.Bool("seed", false, "run sql/init-seed.sql after schema")
	flag.Parse()

	dsn := os.Getenv("MIS_DB_DSN")
	if dsn == "" {
		log.Fatal("MIS_DB_DSN 未配置")
	}

	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("打开数据库失败: %v", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}

	if err := runSQLFile(conn, "sql/init-schema.sql"); err != nil {
		log.Fatalf("执行表结构 SQL 失败: %v", err)
	}
	fmt.Println("表结构已初始化或补齐")

	if *withSeed {
		if err := runSQLFile(conn, "sql/init-seed.sql"); err != nil {
			log.Fatalf("执行初始化数据 SQL 失败: %v", err)
		}
		fmt.Println("默认管理员和开发数据已初始化")
	}
}

func runSQLFile(conn *sql.DB, path string) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	for _, stmt := range splitSQLStatements(string(raw)) {
		if _, err := conn.Exec(stmt); err != nil {
			return fmt.Errorf("%s: %w", compactSQL(stmt), err)
		}
	}
	return nil
}

func splitSQLStatements(input string) []string {
	var out []string
	var current strings.Builder
	inSingle := false
	inDouble := false
	escaped := false
	lineComment := false

	for i := 0; i < len(input); i++ {
		ch := input[i]
		next := byte(0)
		if i+1 < len(input) {
			next = input[i+1]
		}

		if lineComment {
			if ch == '\n' {
				lineComment = false
			}
			continue
		}

		if !inSingle && !inDouble && ch == '-' && next == '-' {
			lineComment = true
			i++
			continue
		}

		if escaped {
			current.WriteByte(ch)
			escaped = false
			continue
		}
		if ch == '\\' && (inSingle || inDouble) {
			current.WriteByte(ch)
			escaped = true
			continue
		}
		if ch == '\'' && !inDouble {
			inSingle = !inSingle
		}
		if ch == '"' && !inSingle {
			inDouble = !inDouble
		}
		if ch == ';' && !inSingle && !inDouble {
			appendSQL(&out, current.String())
			current.Reset()
			continue
		}
		current.WriteByte(ch)
	}
	appendSQL(&out, current.String())
	return out
}

func appendSQL(out *[]string, stmt string) {
	stmt = strings.TrimSpace(stmt)
	if stmt != "" {
		*out = append(*out, stmt)
	}
}

func compactSQL(stmt string) string {
	fields := strings.Fields(stmt)
	if len(fields) == 0 {
		return ""
	}
	text := strings.Join(fields, " ")
	if len(text) > 120 {
		return text[:120] + "..."
	}
	return text
}
