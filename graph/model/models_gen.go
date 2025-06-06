// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type Board struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Position  int32   `json:"position"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Lists     []*List `json:"lists"`
}

type Card struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Content   *string `json:"content,omitempty"`
	ListID    string  `json:"listId"`
	BoardID   string  `json:"boardId"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Position  int32   `json:"position"`
}

type CreateBoardInput struct {
	Name     string `json:"name"`
	Position *int32 `json:"position,omitempty"`
}

type CreateCardInput struct {
	ListID  string  `json:"listId"`
	Title   string  `json:"title"`
	Content *string `json:"content,omitempty"`
	BoardID string  `json:"boardId"`
}

type CreateListInput struct {
	BoardID string `json:"boardId"`
	Name    string `json:"name"`
}

type List struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	BoardID   string  `json:"boardId"`
	CreatedAt string  `json:"createdAt"`
	UpdatedAt string  `json:"updatedAt"`
	Position  int32   `json:"position"`
	Cards     []*Card `json:"cards"`
}

type MoveBoardInput struct {
	ID          string `json:"id"`
	NewPosition int32  `json:"newPosition"`
}

type MoveCardInput struct {
	ID           string `json:"id"`
	TargetListID string `json:"targetListId"`
	NewPosition  int32  `json:"newPosition"`
}

type MoveListInput struct {
	ID          string `json:"id"`
	NewPosition int32  `json:"newPosition"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateBoardInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UpdateCardInput struct {
	ID      string  `json:"id"`
	Title   string  `json:"title"`
	Content *string `json:"content,omitempty"`
}

type UpdateListInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
