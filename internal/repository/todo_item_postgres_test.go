package repository

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
	"todo/internal/models"
)

func TestTodoItemPostgres_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		// wrap err ili logger?
	}
	defer db.Close()

	r := NewTodoItemPostgres(db)

	type args struct {
		listId int
		item   models.TodoItem
	}
	type mockBehavior func(args args, id int)

	testTable := []struct {
		name    string
		mock    mockBehavior
		args    args
		id      int
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				listId: 1,
				item: models.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			args: args{
				listId: 1,
				item: models.TodoItem{
					Title:       "",
					Description: "test description",
				},
			},
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(1, errors.New("some error"))

				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "2nd INSERT error",
			args: args{
				listId: 1,
				item: models.TodoItem{
					Title:       "test title",
					Description: "test description",
				},
			},
			id: 2,
			mockBehavior: func(args args, id int) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO todo_items").WithArgs(args.item.Title, args.item.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO lists_items").WithArgs(args.listId, id).WillReturnResult(errors.New("some error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			testCase.mockBehavior(testCase.args, testCase.id)

			got, err := r.Create(testCase.args.listId, testCase.args.item)
			if testCase.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, got, testCase.id)
			}
		})
	}
}
