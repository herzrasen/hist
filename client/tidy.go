package client

import "fmt"

func (c *Client) Tidy() error {
	records, err := c.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("client:Tidy: list: %w", err)
	}
	var idsToDelete []int64
	for _, r := range records {
		if c.Config.IsExcluded(r.Command) {
			idsToDelete = append(idsToDelete, r.Id)
			fmt.Printf("Deleting: %s\n", r.Command)
		}
	}
	return c.Delete(DeleteOptions{Ids: idsToDelete})
}
