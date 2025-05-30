package services

import (
	"trello-backend/internal/models"
	"trello-backend/internal/repositories"
)

type CardService interface {
	CreateCard(listID uint, boardID uint, title, content string) (*models.Card, error)
	GetCards(listID uint) ([]models.Card, error)
	GetCardByID(id uint) (*models.Card, error)
	UpdateCard(id uint, title, content string) error
	DeleteCard(id uint) error
	MoveCard(id, targetListID uint, newPosition int) error
	GetCardsByBoardID(boardID uint) ([]models.Card, error) // 新增
	GetCardsByListIDs(listIDs []uint) (map[uint][]models.Card, error)
}

type cardService struct {
	cardRepo repositories.CardRepository
}

func NewCardService(repo repositories.CardRepository) CardService {
	return &cardService{cardRepo: repo}
}

func (s *cardService) CreateCard(listID uint, boardID uint, title, content string) (*models.Card, error) {
	// 取得該 list 目前所有卡片
	cards, err := s.cardRepo.GetCardsByListID(listID)
	if err != nil {
		return nil, err
	}
	// 所有卡片 position +1
	for i := range cards {
		cards[i].Position++
		if err := s.cardRepo.UpdateCard(&cards[i]); err != nil {
			return nil, err
		}
	}
	// 新卡片 position = 0
	card := &models.Card{ListID: listID, BoardID: boardID, Title: title, Content: content, Position: 0}
	if err := s.cardRepo.CreateCard(card); err != nil {
		return nil, err
	}
	return card, nil
}

func (s *cardService) GetCards(listID uint) ([]models.Card, error) {
	return s.cardRepo.GetCardsByListID(listID)
}

func (s *cardService) GetCardByID(id uint) (*models.Card, error) {
	return s.cardRepo.GetCardByID(id)
}

func (s *cardService) UpdateCard(id uint, title, content string) error {
	card, err := s.cardRepo.GetCardByID(id)
	if err != nil {
		return err
	}
	card.Title = title
	card.Content = content
	return s.cardRepo.UpdateCard(card)
}

func (s *cardService) DeleteCard(id uint) error {
	return s.cardRepo.DeleteCard(id)
}

func (s *cardService) MoveCard(id, targetListID uint, newPosition int) error {
	card, err := s.cardRepo.GetCardByID(id)
	if err != nil {
		return err
	}
	oldListID := card.ListID
	oldPos := card.Position
	// 調整原清單中的卡片位置
	oldCards, err := s.cardRepo.GetCardsByListID(oldListID)
	if err != nil {
		return err
	}
	for _, c := range oldCards {
		if c.ID == id {
			continue
		}
		if c.Position > oldPos {
			c.Position--
			if err := s.cardRepo.UpdateCard(&c); err != nil {
				return err
			}
		}
	}
	// 插入到新清單
	card.ListID = targetListID
	newCards, err := s.cardRepo.GetCardsByListID(targetListID)
	if err != nil {
		return err
	}
	if newPosition > len(newCards) {
		newPosition = len(newCards)
	}
	for _, c := range newCards {
		if c.Position >= newPosition {
			c.Position++
			if err := s.cardRepo.UpdateCard(&c); err != nil {
				return err
			}
		}
	}
	card.Position = newPosition
	return s.cardRepo.UpdateCard(card)
}

func (s *cardService) GetCardsByBoardID(boardID uint) ([]models.Card, error) {
	return s.cardRepo.GetCardsByBoardID(boardID)
}

func (s *cardService) GetCardsByListIDs(listIDs []uint) (map[uint][]models.Card, error) {
	return s.cardRepo.GetCardsByListIDs(listIDs)
}
