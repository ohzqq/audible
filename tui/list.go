package tui

import (
	"fmt"
	"log"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ohzqq/audible"
	"github.com/ohzqq/bubbles/list"
	"golang.org/x/term"
)

type List struct {
	list.Model
	*audible.ProductsResponse
}

type Item struct {
	audible.Product
}

func New(products *audible.ProductsResponse) List {
	var items []list.Item
	for _, product := range products.Products {
		items = append(items, NewItem(product))
	}

	w, h := TermSize()

	m := list.New(items, list.NewDefaultDelegate(), w, h)
	m.SetNoLimit()

	return List{
		Model:            m,
		ProductsResponse: products,
	}
}

func (l List) Run() []audible.Product {
	p := tea.NewProgram(l)
	_, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	sel := l.ToggledItems()
	books := make([]audible.Product, len(sel))
	for _, idx := range sel {
		books[idx] = l.Products[idx]
	}

	return books
}

func (l List) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			return l, tea.Quit
		}
	}

	var cmd tea.Cmd
	l.Model, cmd = l.Model.Update(msg)
	return l, cmd
}

func NewItem(p audible.Product) Item {
	return Item{
		Product: p,
	}
}

func (i Item) authors() string {
	var auths []string
	for _, a := range i.Authors {
		auths = append(auths, a["name"])
	}
	return strings.Join(auths, " & ")
}

func (i Item) FilterValue() string {
	return fmt.Sprintf("%s %s", i.Product.Title, i.authors())
}

func (i Item) Title() string {
	return fmt.Sprintf("%s & %s", i.Product.Title, i.authors())
}

func (i Item) Description() string {
	return i.PublisherSummary
}

func TermSize() (int, int) {
	w, h, _ := term.GetSize(int(os.Stdin.Fd()))
	println(w)
	println(h)
	return w, h
}
