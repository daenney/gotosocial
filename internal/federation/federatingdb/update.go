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

package federatingdb

import (
	"context"

	"codeberg.org/gruf/go-logger/v2/level"
	"github.com/superseriousbusiness/activity/streams/vocab"
	"github.com/superseriousbusiness/gotosocial/internal/ap"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtserror"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
	"github.com/superseriousbusiness/gotosocial/internal/log"
	"github.com/superseriousbusiness/gotosocial/internal/messages"
)

// Update sets an existing entry to the database based on the value's
// id.
//
// Note that Activity values received from federated peers may also be
// updated in the database this way if the Federating Protocol is
// enabled. The client may freely decide to store only the id instead of
// the entire value.
//
// The library makes this call only after acquiring a lock first.
func (f *federatingDB) Update(ctx context.Context, asType vocab.Type) error {
	l := log.WithContext(ctx)

	if log.Level() >= level.DEBUG {
		i, err := marshalItem(asType)
		if err != nil {
			return err
		}
		l = l.WithField("update", i)
		l.Debug("entering Update")
	}

	receivingAccount, requestingAccount, internal := extractFromCtx(ctx)
	if internal {
		return nil // Already processed.
	}

	switch asType.GetTypeName() {
	case ap.ActorApplication, ap.ActorGroup, ap.ActorOrganization, ap.ActorPerson, ap.ActorService:
		return f.updateAccountable(ctx, receivingAccount, requestingAccount, asType)
	}

	return nil
}

func (f *federatingDB) updateAccountable(ctx context.Context, receivingAcct *gtsmodel.Account, requestingAcct *gtsmodel.Account, asType vocab.Type) error {
	// Ensure delivered asType is a valid Accountable model.
	accountable, ok := asType.(ap.Accountable)
	if !ok {
		return gtserror.Newf("could not convert vocab.Type %T to Accountable", asType)
	}

	// Extract AP URI of the updated Accountable model.
	idProp := accountable.GetJSONLDId()
	if idProp == nil || !idProp.IsIRI() {
		return gtserror.New("Accountable id prop was nil or not IRI")
	}
	updatedAcctURI := idProp.GetIRI()

	// Don't try to update local accounts, it will break things.
	if updatedAcctURI.Host == config.GetHost() {
		return nil
	}

	// Ensure Accountable and requesting account are one and the same.
	if updatedAcctURIStr := updatedAcctURI.String(); requestingAcct.URI != updatedAcctURIStr {
		return gtserror.Newf("update for %s was requested by %s, this is not valid", updatedAcctURIStr, requestingAcct.URI)
	}

	// Pass in to the processor the existing version of the requesting
	// account that we have, plus the Accountable representation that
	// was delivered along with the Update, for further asynchronous
	// updating of eg., avatar/header, emojis, etc. The actual db
	// inserts/updates will take place there.
	f.state.Workers.EnqueueFederator(ctx, messages.FromFederator{
		APObjectType:     ap.ObjectProfile,
		APActivityType:   ap.ActivityUpdate,
		GTSModel:         requestingAcct,
		APObjectModel:    accountable,
		ReceivingAccount: receivingAcct,
	})

	return nil
}
