package pg

import "github.com/varunamachi/idx/core"

type CredentialStorage struct {
	hasher core.Hasher
}

func (pcs *CredentialStorage) SetPassword(itemType, id, password string) error {
	return nil
}
func (pcs *CredentialStorage) UpdatePassword(
	itemType, id, oldPw, newPw string) error {
	return nil
}
func (pcs *CredentialStorage) Verify(itemType, id, password string) error {
	return nil
}
