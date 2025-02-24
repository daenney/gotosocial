// GoToSocial
// Copyright (C) GoToSocial Authors admin@gotosocial.org
// SPDX-License-Identifier: AGPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package gtsmodel

import "time"

// AccountNote stores a private note from a local account related to any account.
type AccountNote struct {
	ID              string    `validate:"required,ulid" bun:"type:CHAR(26),pk,nullzero,notnull,unique"`                                              // id of this item in the database
	CreatedAt       time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`                                       // when was item created
	UpdatedAt       time.Time `validate:"-" bun:"type:timestamptz,nullzero,notnull,default:current_timestamp"`                                       // when was item last updated
	AccountID       string    `validate:"required,ulid" bun:"type:CHAR(26),unique:account_notes_account_id_target_account_id_uniq,notnull,nullzero"` // ID of the local account that created the note
	Account         *Account  `validate:"-" bun:"rel:belongs-to"`                                                                                    // Account corresponding to accountID
	TargetAccountID string    `validate:"required,ulid" bun:"type:CHAR(26),unique:account_notes_account_id_target_account_id_uniq,notnull,nullzero"` // Who is the target of this note?
	TargetAccount   *Account  `validate:"-" bun:"rel:belongs-to"`                                                                                    // Account corresponding to targetAccountID
	Comment         string    `validate:"-" bun:""`                                                                                                  // The text of the note.
}
