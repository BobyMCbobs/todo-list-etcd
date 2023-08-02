package todolist

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/BobyMCbobs/todo-list-etcd/pkg/etcd"
	"github.com/BobyMCbobs/todo-list-etcd/pkg/types"
)

type Manager struct {
	clientset *etcd.Client
}

func NewManager(clientset *etcd.Client) *Manager {
	mgr := &Manager{
		clientset: clientset,
	}
	return mgr
}

func (m *Manager) ValidateList(input *types.List) error {
	if len(input.Name) == 0 {
		return fmt.Errorf("error: name must not be empty")
	}
	if len(input.Name) > 35 {
		return fmt.Errorf("error: name must be under 35 character in length")
	}
	if len(input.Description) > 45 {
		return fmt.Errorf("error: name must be under 45 character in length")
	}
	if _, err := uuid.Parse(input.ID); input.ID != "" && err != nil {
		return fmt.Errorf("error: failed to parse id uuid")
	}
	if _, err := uuid.Parse(input.AuthorID); input.AuthorID != "" && err != nil {
		return fmt.Errorf("error: failed to parse author id uuid")
	}
	return nil
}

func (m *Manager) ValidateItem(input *types.Item) error {
	if len(input.Name) == 0 {
		return fmt.Errorf("error: name must not be empty")
	}
	if len(input.Name) > 35 {
		return fmt.Errorf("error: name must be under 35 character in length")
	}
	if len(input.Description) > 45 {
		return fmt.Errorf("error: name must be under 45 character in length")
	}
	if _, err := uuid.Parse(input.ID); input.ID != "" && err != nil {
		return fmt.Errorf("error: failed to parse id uuid")
	}
	if _, err := uuid.Parse(input.AuthorID); input.AuthorID != "" && err != nil {
		return fmt.Errorf("error: failed to parse author id uuid")
	}
	if _, err := uuid.Parse(input.ListID); input.ListID != "" && err != nil {
		return fmt.Errorf("error: failed to parse list id uuid")
	}
	if input.ListID == "" {
		return fmt.Errorf("error: missing list id")
	}
	if list, err := m.Lists().Get(context.TODO(), input.ListID); err != nil || list == nil {
		return fmt.Errorf("error: list not found")
	}
	return nil
}

type listManager struct {
	clientset *etcd.Client
	manager   *Manager
}

func (m *Manager) Lists() *listManager {
	return &listManager{
		clientset: m.clientset,
		manager:   m,
	}
}

func (m *listManager) Get(ctx context.Context, id string) (*types.List, error) {
	if id == "" {
		return nil, fmt.Errorf("error: id is empty")
	}
	val, err := m.clientset.Get("/list/" + id)
	if err != nil {
		return nil, err
	}
	var list types.List
	if err := json.Unmarshal(val.Value, &list); err != nil {
		return nil, err
	}
	list.Revision = val.Version
	items, err := m.manager.Items(id).List(context.TODO())
	if err != nil {
		return nil, err
	}
	list.Items = items
	return &list, nil
}

func (m *listManager) List(ctx context.Context) ([]*types.List, error) {
	vals, err := m.clientset.ListWithPrefix("/list/")
	if err != nil {
		return nil, err
	}
	var lists []*types.List
	for _, val := range vals {
		var list *types.List
		if err := json.Unmarshal(val.Value, &list); err != nil {
			return nil, err
		}
		list.Revision = val.Version
		items, err := m.manager.Items(list.ID).List(context.TODO())
		if err != nil {
			return nil, err
		}
		list.Items = items
		lists = append(lists, list)
	}
	return lists, nil
}

func (m *listManager) Put(ctx context.Context, item *types.List) (*types.List, error) {
	if err := m.manager.ValidateList(item); err != nil {
		return nil, err
	}
	currentTimestamp := time.Now().Unix()
	var id string
	if item.ID != "" {
		existingList, err := m.Get(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		id = existingList.ID
		item.CreationTimestamp = existingList.CreationTimestamp
	} else {
		id = uuid.New().String()
		item.CreationTimestamp = fmt.Sprintf("%v", currentTimestamp)
	}

	item.ID = id
	item.ModificationTimestamp = fmt.Sprintf("%v", currentTimestamp)
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	resp, err := m.clientset.Put("/list/"+id, string(itemBytes))
	if err != nil {
		return nil, err
	}
	item.Revision = resp.Header.Revision
	return item, nil
}

func (m *listManager) Delete(ctx context.Context, id string) error {
	if err := m.manager.Items(id).DeleteAll(ctx); err != nil {
		return err
	}
	_, err := m.clientset.DeleteWithPrefix("/list/" + id)
	if err != nil {
		return err
	}
	return nil
}

type itemManager struct {
	listid string

	clientset *etcd.Client
	manager   *Manager
}

func (m *Manager) Items(listid string) *itemManager {
	return &itemManager{
		listid:    listid,
		clientset: m.clientset,
		manager:   m,
	}
}

func (m *itemManager) Get(ctx context.Context, id string) (*types.Item, error) {
	val, err := m.clientset.Get("/item/" + m.listid + "/" + id)
	if err != nil {
		return nil, err
	}
	var item types.Item
	if err := json.Unmarshal(val.Value, &item); err != nil {
		return nil, err
	}
	item.Revision = val.Version
	return &item, nil
}

func (m *itemManager) List(ctx context.Context) ([]*types.Item, error) {
	vals, err := m.clientset.ListWithPrefix("/item/" + m.listid + "/")
	if err != nil {
		return nil, err
	}
	var items []*types.Item
	for _, val := range vals {
		var item *types.Item
		if err := json.Unmarshal(val.Value, &item); err != nil {
			return nil, err
		}
		item.Revision = val.Version
		items = append(items, item)
	}
	return items, nil
}

func (m *itemManager) Put(ctx context.Context, item *types.Item) (*types.Item, error) {
	if err := m.manager.ValidateItem(item); err != nil {
		return nil, err
	}
	currentTimestamp := time.Now().Unix()
	var id string
	if item.ID != "" {
		existingList, err := m.Get(ctx, item.ID)
		if err != nil {
			return nil, err
		}
		id = existingList.ID
		item.CreationTimestamp = existingList.CreationTimestamp
	} else {
		id = uuid.New().String()
		item.CreationTimestamp = fmt.Sprintf("%v", currentTimestamp)
	}

	item.ID = id
	item.ListID = m.listid
	item.ModificationTimestamp = fmt.Sprintf("%v", currentTimestamp)
	itemBytes, err := json.Marshal(item)
	if err != nil {
		return nil, err
	}
	resp, err := m.clientset.Put("/item/"+m.listid+"/"+id, string(itemBytes))
	if err != nil {
		return nil, err
	}
	item.Revision = resp.Header.Revision
	return item, nil
}

func (m *itemManager) Delete(ctx context.Context, id string) error {
	_, err := m.clientset.Delete("/item/" + m.listid + "/" + id)
	if err != nil {
		return err
	}
	return nil
}

func (m *itemManager) DeleteAll(ctx context.Context) error {
	_, err := m.clientset.DeleteWithPrefix("/item/" + m.listid + "/")
	if err != nil {
		return err
	}
	return nil
}
