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

package cache

import (
	"codeberg.org/gruf/go-cache/v3/result"
	"codeberg.org/gruf/go-cache/v3/ttl"
	"github.com/superseriousbusiness/gotosocial/internal/cache/domain"
	"github.com/superseriousbusiness/gotosocial/internal/config"
	"github.com/superseriousbusiness/gotosocial/internal/gtsmodel"
)

type GTSCaches struct {
	account          *result.Cache[*gtsmodel.Account]
	accountNote      *result.Cache[*gtsmodel.AccountNote]
	block            *result.Cache[*gtsmodel.Block]
	blockIDs         *SliceCache[string]
	domainBlock      *domain.BlockCache
	emoji            *result.Cache[*gtsmodel.Emoji]
	emojiCategory    *result.Cache[*gtsmodel.EmojiCategory]
	follow           *result.Cache[*gtsmodel.Follow]
	followIDs        *SliceCache[string]
	followRequest    *result.Cache[*gtsmodel.FollowRequest]
	followRequestIDs *SliceCache[string]
	instance         *result.Cache[*gtsmodel.Instance]
	list             *result.Cache[*gtsmodel.List]
	listEntry        *result.Cache[*gtsmodel.ListEntry]
	marker           *result.Cache[*gtsmodel.Marker]
	media            *result.Cache[*gtsmodel.MediaAttachment]
	mention          *result.Cache[*gtsmodel.Mention]
	notification     *result.Cache[*gtsmodel.Notification]
	report           *result.Cache[*gtsmodel.Report]
	status           *result.Cache[*gtsmodel.Status]
	statusFave       *result.Cache[*gtsmodel.StatusFave]
	tag              *result.Cache[*gtsmodel.Tag]
	tombstone        *result.Cache[*gtsmodel.Tombstone]
	user             *result.Cache[*gtsmodel.User]

	// TODO: move out of GTS caches since unrelated to DB.
	webfinger *ttl.Cache[string, string]
}

// Init will initialize all the gtsmodel caches in this collection.
// NOTE: the cache MUST NOT be in use anywhere, this is not thread-safe.
func (c *GTSCaches) Init() {
	c.initAccount()
	c.initAccountNote()
	c.initBlock()
	c.initBlockIDs()
	c.initDomainBlock()
	c.initEmoji()
	c.initEmojiCategory()
	c.initFollow()
	c.initFollowIDs()
	c.initFollowRequest()
	c.initFollowRequestIDs()
	c.initInstance()
	c.initList()
	c.initListEntry()
	c.initMarker()
	c.initMedia()
	c.initMention()
	c.initNotification()
	c.initReport()
	c.initStatus()
	c.initStatusFave()
	c.initTag()
	c.initTombstone()
	c.initUser()
	c.initWebfinger()
}

// Start will attempt to start all of the gtsmodel caches, or panic.
func (c *GTSCaches) Start() {
	tryStart(c.account, config.GetCacheGTSAccountSweepFreq())
	tryStart(c.accountNote, config.GetCacheGTSAccountNoteSweepFreq())
	tryStart(c.block, config.GetCacheGTSBlockSweepFreq())
	tryUntil("starting block IDs cache", 5, func() bool {
		if sweep := config.GetCacheGTSBlockIDsSweepFreq(); sweep > 0 {
			return c.blockIDs.Start(sweep)
		}
		return true
	})
	tryStart(c.emoji, config.GetCacheGTSEmojiSweepFreq())
	tryStart(c.emojiCategory, config.GetCacheGTSEmojiCategorySweepFreq())
	tryStart(c.follow, config.GetCacheGTSFollowSweepFreq())
	tryUntil("starting follow IDs cache", 5, func() bool {
		if sweep := config.GetCacheGTSFollowIDsSweepFreq(); sweep > 0 {
			return c.followIDs.Start(sweep)
		}
		return true
	})
	tryStart(c.followRequest, config.GetCacheGTSFollowRequestSweepFreq())
	tryUntil("starting follow request IDs cache", 5, func() bool {
		if sweep := config.GetCacheGTSFollowRequestIDsSweepFreq(); sweep > 0 {
			return c.followRequestIDs.Start(sweep)
		}
		return true
	})
	tryStart(c.instance, config.GetCacheGTSInstanceSweepFreq())
	tryStart(c.list, config.GetCacheGTSListSweepFreq())
	tryStart(c.listEntry, config.GetCacheGTSListEntrySweepFreq())
	tryStart(c.marker, config.GetCacheGTSMarkerSweepFreq())
	tryStart(c.media, config.GetCacheGTSMediaSweepFreq())
	tryStart(c.mention, config.GetCacheGTSMentionSweepFreq())
	tryStart(c.notification, config.GetCacheGTSNotificationSweepFreq())
	tryStart(c.report, config.GetCacheGTSReportSweepFreq())
	tryStart(c.status, config.GetCacheGTSStatusSweepFreq())
	tryStart(c.statusFave, config.GetCacheGTSStatusFaveSweepFreq())
	tryStart(c.tag, config.GetCacheGTSTagSweepFreq())
	tryStart(c.tombstone, config.GetCacheGTSTombstoneSweepFreq())
	tryStart(c.user, config.GetCacheGTSUserSweepFreq())
	tryUntil("starting *gtsmodel.Webfinger cache", 5, func() bool {
		if sweep := config.GetCacheGTSWebfingerSweepFreq(); sweep > 0 {
			return c.webfinger.Start(sweep)
		}
		return true
	})
}

// Stop will attempt to stop all of the gtsmodel caches, or panic.
func (c *GTSCaches) Stop() {
	tryStop(c.account, config.GetCacheGTSAccountSweepFreq())
	tryStop(c.accountNote, config.GetCacheGTSAccountNoteSweepFreq())
	tryStop(c.block, config.GetCacheGTSBlockSweepFreq())
	tryUntil("stopping block IDs cache", 5, func() bool {
		if config.GetCacheGTSBlockIDsSweepFreq() > 0 {
			return c.blockIDs.Stop()
		}
		return true
	})
	tryStop(c.emoji, config.GetCacheGTSEmojiSweepFreq())
	tryStop(c.emojiCategory, config.GetCacheGTSEmojiCategorySweepFreq())
	tryStop(c.follow, config.GetCacheGTSFollowSweepFreq())
	tryUntil("stopping follow IDs cache", 5, func() bool {
		if config.GetCacheGTSFollowIDsSweepFreq() > 0 {
			return c.followIDs.Stop()
		}
		return true
	})
	tryStop(c.followRequest, config.GetCacheGTSFollowRequestSweepFreq())
	tryUntil("stopping follow request IDs cache", 5, func() bool {
		if config.GetCacheGTSFollowRequestIDsSweepFreq() > 0 {
			return c.followRequestIDs.Stop()
		}
		return true
	})
	tryStop(c.instance, config.GetCacheGTSInstanceSweepFreq())
	tryStop(c.list, config.GetCacheGTSListSweepFreq())
	tryStop(c.listEntry, config.GetCacheGTSListEntrySweepFreq())
	tryStop(c.marker, config.GetCacheGTSMarkerSweepFreq())
	tryStop(c.media, config.GetCacheGTSMediaSweepFreq())
	tryStop(c.mention, config.GetCacheGTSNotificationSweepFreq())
	tryStop(c.notification, config.GetCacheGTSNotificationSweepFreq())
	tryStop(c.report, config.GetCacheGTSReportSweepFreq())
	tryStop(c.status, config.GetCacheGTSStatusSweepFreq())
	tryStop(c.statusFave, config.GetCacheGTSStatusFaveSweepFreq())
	tryStop(c.tag, config.GetCacheGTSTagSweepFreq())
	tryStop(c.tombstone, config.GetCacheGTSTombstoneSweepFreq())
	tryStop(c.user, config.GetCacheGTSUserSweepFreq())
	tryUntil("stopping *gtsmodel.Webfinger cache", 5, func() bool {
		if config.GetCacheGTSWebfingerSweepFreq() > 0 {
			return c.webfinger.Stop()
		}
		return true
	})
}

// Account provides access to the gtsmodel Account database cache.
func (c *GTSCaches) Account() *result.Cache[*gtsmodel.Account] {
	return c.account
}

// AccountNote provides access to the gtsmodel Note database cache.
func (c *GTSCaches) AccountNote() *result.Cache[*gtsmodel.AccountNote] {
	return c.accountNote
}

// Block provides access to the gtsmodel Block (account) database cache.
func (c *GTSCaches) Block() *result.Cache[*gtsmodel.Block] {
	return c.block
}

// FollowIDs provides access to the block IDs database cache.
func (c *GTSCaches) BlockIDs() *SliceCache[string] {
	return c.blockIDs
}

// DomainBlock provides access to the domain block database cache.
func (c *GTSCaches) DomainBlock() *domain.BlockCache {
	return c.domainBlock
}

// Emoji provides access to the gtsmodel Emoji database cache.
func (c *GTSCaches) Emoji() *result.Cache[*gtsmodel.Emoji] {
	return c.emoji
}

// EmojiCategory provides access to the gtsmodel EmojiCategory database cache.
func (c *GTSCaches) EmojiCategory() *result.Cache[*gtsmodel.EmojiCategory] {
	return c.emojiCategory
}

// Follow provides access to the gtsmodel Follow database cache.
func (c *GTSCaches) Follow() *result.Cache[*gtsmodel.Follow] {
	return c.follow
}

// FollowIDs provides access to the follower / following IDs database cache.
// THIS CACHE IS KEYED AS THE FOLLOWING {prefix}{accountID} WHERE PREFIX IS:
// - '>'  for following IDs
// - 'l>' for local following IDs
// - '<'  for follower IDs
// - 'l<' for local follower IDs
func (c *GTSCaches) FollowIDs() *SliceCache[string] {
	return c.followIDs
}

// FollowRequest provides access to the gtsmodel FollowRequest database cache.
func (c *GTSCaches) FollowRequest() *result.Cache[*gtsmodel.FollowRequest] {
	return c.followRequest
}

// FollowRequestIDs provides access to the follow requester / requesting IDs database
// cache. THIS CACHE IS KEYED AS THE FOLLOWING {prefix}{accountID} WHERE PREFIX IS:
// - '>'  for following IDs
// - '<'  for follower IDs
func (c *GTSCaches) FollowRequestIDs() *SliceCache[string] {
	return c.followRequestIDs
}

// Instance provides access to the gtsmodel Instance database cache.
func (c *GTSCaches) Instance() *result.Cache[*gtsmodel.Instance] {
	return c.instance
}

// List provides access to the gtsmodel List database cache.
func (c *GTSCaches) List() *result.Cache[*gtsmodel.List] {
	return c.list
}

// ListEntry provides access to the gtsmodel ListEntry database cache.
func (c *GTSCaches) ListEntry() *result.Cache[*gtsmodel.ListEntry] {
	return c.listEntry
}

// Marker provides access to the gtsmodel Marker database cache.
func (c *GTSCaches) Marker() *result.Cache[*gtsmodel.Marker] {
	return c.marker
}

// Media provides access to the gtsmodel Media database cache.
func (c *GTSCaches) Media() *result.Cache[*gtsmodel.MediaAttachment] {
	return c.media
}

// Mention provides access to the gtsmodel Mention database cache.
func (c *GTSCaches) Mention() *result.Cache[*gtsmodel.Mention] {
	return c.mention
}

// Notification provides access to the gtsmodel Notification database cache.
func (c *GTSCaches) Notification() *result.Cache[*gtsmodel.Notification] {
	return c.notification
}

// Report provides access to the gtsmodel Report database cache.
func (c *GTSCaches) Report() *result.Cache[*gtsmodel.Report] {
	return c.report
}

// Status provides access to the gtsmodel Status database cache.
func (c *GTSCaches) Status() *result.Cache[*gtsmodel.Status] {
	return c.status
}

// StatusFave provides access to the gtsmodel StatusFave database cache.
func (c *GTSCaches) StatusFave() *result.Cache[*gtsmodel.StatusFave] {
	return c.statusFave
}

// Tag provides access to the gtsmodel Tag database cache.
func (c *GTSCaches) Tag() *result.Cache[*gtsmodel.Tag] {
	return c.tag
}

// Tombstone provides access to the gtsmodel Tombstone database cache.
func (c *GTSCaches) Tombstone() *result.Cache[*gtsmodel.Tombstone] {
	return c.tombstone
}

// User provides access to the gtsmodel User database cache.
func (c *GTSCaches) User() *result.Cache[*gtsmodel.User] {
	return c.user
}

// Webfinger provides access to the webfinger URL cache.
func (c *GTSCaches) Webfinger() *ttl.Cache[string, string] {
	return c.webfinger
}

func (c *GTSCaches) initAccount() {
	c.account = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "URL"},
		{Name: "Username.Domain"},
		{Name: "PublicKeyURI"},
		{Name: "InboxURI"},
		{Name: "OutboxURI"},
		{Name: "FollowersURI"},
		{Name: "FollowingURI"},
	}, func(a1 *gtsmodel.Account) *gtsmodel.Account {
		a2 := new(gtsmodel.Account)
		*a2 = *a1
		return a2
	}, config.GetCacheGTSAccountMaxSize())
	c.account.SetTTL(config.GetCacheGTSAccountTTL(), true)
	c.account.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initAccountNote() {
	c.accountNote = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "AccountID.TargetAccountID"},
	}, func(n1 *gtsmodel.AccountNote) *gtsmodel.AccountNote {
		n2 := new(gtsmodel.AccountNote)
		*n2 = *n1
		return n2
	}, config.GetCacheGTSAccountNoteMaxSize())
	c.accountNote.SetTTL(config.GetCacheGTSAccountNoteTTL(), true)
	c.accountNote.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initBlock() {
	c.block = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "AccountID.TargetAccountID"},
		{Name: "AccountID", Multi: true},
		{Name: "TargetAccountID", Multi: true},
	}, func(b1 *gtsmodel.Block) *gtsmodel.Block {
		b2 := new(gtsmodel.Block)
		*b2 = *b1
		return b2
	}, config.GetCacheGTSBlockMaxSize())
	c.block.SetTTL(config.GetCacheGTSBlockTTL(), true)
	c.block.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initBlockIDs() {
	c.blockIDs = &SliceCache[string]{Cache: ttl.New[string, []string](
		0,
		config.GetCacheGTSBlockIDsMaxSize(),
		config.GetCacheGTSBlockIDsTTL(),
	)}
}

func (c *GTSCaches) initDomainBlock() {
	c.domainBlock = new(domain.BlockCache)
}

func (c *GTSCaches) initEmoji() {
	c.emoji = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "Shortcode.Domain"},
		{Name: "ImageStaticURL"},
		{Name: "CategoryID", Multi: true},
	}, func(e1 *gtsmodel.Emoji) *gtsmodel.Emoji {
		e2 := new(gtsmodel.Emoji)
		*e2 = *e1
		return e2
	}, config.GetCacheGTSEmojiMaxSize())
	c.emoji.SetTTL(config.GetCacheGTSEmojiTTL(), true)
	c.emoji.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initEmojiCategory() {
	c.emojiCategory = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "Name"},
	}, func(c1 *gtsmodel.EmojiCategory) *gtsmodel.EmojiCategory {
		c2 := new(gtsmodel.EmojiCategory)
		*c2 = *c1
		return c2
	}, config.GetCacheGTSEmojiCategoryMaxSize())
	c.emojiCategory.SetTTL(config.GetCacheGTSEmojiCategoryTTL(), true)
	c.emojiCategory.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initFollow() {
	c.follow = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "AccountID.TargetAccountID"},
		{Name: "AccountID", Multi: true},
		{Name: "TargetAccountID", Multi: true},
	}, func(f1 *gtsmodel.Follow) *gtsmodel.Follow {
		f2 := new(gtsmodel.Follow)
		*f2 = *f1
		return f2
	}, config.GetCacheGTSFollowMaxSize())
	c.follow.SetTTL(config.GetCacheGTSFollowTTL(), true)
}

func (c *GTSCaches) initFollowIDs() {
	c.followIDs = &SliceCache[string]{Cache: ttl.New[string, []string](
		0,
		config.GetCacheGTSFollowIDsMaxSize(),
		config.GetCacheGTSFollowIDsTTL(),
	)}
}

func (c *GTSCaches) initFollowRequest() {
	c.followRequest = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "AccountID.TargetAccountID"},
		{Name: "AccountID", Multi: true},
		{Name: "TargetAccountID", Multi: true},
	}, func(f1 *gtsmodel.FollowRequest) *gtsmodel.FollowRequest {
		f2 := new(gtsmodel.FollowRequest)
		*f2 = *f1
		return f2
	}, config.GetCacheGTSFollowRequestMaxSize())
	c.followRequest.SetTTL(config.GetCacheGTSFollowRequestTTL(), true)
}

func (c *GTSCaches) initFollowRequestIDs() {
	c.followRequestIDs = &SliceCache[string]{Cache: ttl.New[string, []string](
		0,
		config.GetCacheGTSFollowRequestIDsMaxSize(),
		config.GetCacheGTSFollowRequestIDsTTL(),
	)}
}

func (c *GTSCaches) initInstance() {
	c.instance = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "Domain"},
	}, func(i1 *gtsmodel.Instance) *gtsmodel.Instance {
		i2 := new(gtsmodel.Instance)
		*i2 = *i1
		return i1
	}, config.GetCacheGTSInstanceMaxSize())
	c.instance.SetTTL(config.GetCacheGTSInstanceTTL(), true)
	c.emojiCategory.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initList() {
	c.list = result.New([]result.Lookup{
		{Name: "ID"},
	}, func(l1 *gtsmodel.List) *gtsmodel.List {
		l2 := new(gtsmodel.List)
		*l2 = *l1
		return l2
	}, config.GetCacheGTSListMaxSize())
	c.list.SetTTL(config.GetCacheGTSListTTL(), true)
	c.list.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initListEntry() {
	c.listEntry = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "ListID", Multi: true},
		{Name: "FollowID", Multi: true},
	}, func(l1 *gtsmodel.ListEntry) *gtsmodel.ListEntry {
		l2 := new(gtsmodel.ListEntry)
		*l2 = *l1
		return l2
	}, config.GetCacheGTSListEntryMaxSize())
	c.list.SetTTL(config.GetCacheGTSListEntryTTL(), true)
	c.list.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initMarker() {
	c.marker = result.New([]result.Lookup{
		{Name: "AccountID.Name"},
	}, func(m1 *gtsmodel.Marker) *gtsmodel.Marker {
		m2 := new(gtsmodel.Marker)
		*m2 = *m1
		return m2
	}, config.GetCacheGTSMarkerMaxSize())
	c.marker.SetTTL(config.GetCacheGTSMarkerTTL(), true)
	c.marker.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initMedia() {
	c.media = result.New([]result.Lookup{
		{Name: "ID"},
	}, func(m1 *gtsmodel.MediaAttachment) *gtsmodel.MediaAttachment {
		m2 := new(gtsmodel.MediaAttachment)
		*m2 = *m1
		return m2
	}, config.GetCacheGTSMediaMaxSize())
	c.media.SetTTL(config.GetCacheGTSMediaTTL(), true)
	c.media.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initMention() {
	c.mention = result.New([]result.Lookup{
		{Name: "ID"},
	}, func(m1 *gtsmodel.Mention) *gtsmodel.Mention {
		m2 := new(gtsmodel.Mention)
		*m2 = *m1
		return m2
	}, config.GetCacheGTSMentionMaxSize())
	c.mention.SetTTL(config.GetCacheGTSMentionTTL(), true)
	c.mention.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initNotification() {
	c.notification = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "NotificationType.TargetAccountID.OriginAccountID.StatusID"},
	}, func(n1 *gtsmodel.Notification) *gtsmodel.Notification {
		n2 := new(gtsmodel.Notification)
		*n2 = *n1
		return n2
	}, config.GetCacheGTSNotificationMaxSize())
	c.notification.SetTTL(config.GetCacheGTSNotificationTTL(), true)
	c.notification.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initReport() {
	c.report = result.New([]result.Lookup{
		{Name: "ID"},
	}, func(r1 *gtsmodel.Report) *gtsmodel.Report {
		r2 := new(gtsmodel.Report)
		*r2 = *r1
		return r2
	}, config.GetCacheGTSReportMaxSize())
	c.report.SetTTL(config.GetCacheGTSReportTTL(), true)
	c.report.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initStatus() {
	c.status = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
		{Name: "URL"},
	}, func(s1 *gtsmodel.Status) *gtsmodel.Status {
		s2 := new(gtsmodel.Status)
		*s2 = *s1
		return s2
	}, config.GetCacheGTSStatusMaxSize())
	c.status.SetTTL(config.GetCacheGTSStatusTTL(), true)
	c.status.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initStatusFave() {
	c.statusFave = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "AccountID.StatusID"},
	}, func(f1 *gtsmodel.StatusFave) *gtsmodel.StatusFave {
		f2 := new(gtsmodel.StatusFave)
		*f2 = *f1
		return f2
	}, config.GetCacheGTSStatusFaveMaxSize())
	c.status.SetTTL(config.GetCacheGTSStatusFaveTTL(), true)
	c.status.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initTag() {
	c.tag = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "Name"},
	}, func(m1 *gtsmodel.Tag) *gtsmodel.Tag {
		m2 := new(gtsmodel.Tag)
		*m2 = *m1
		return m2
	}, config.GetCacheGTSTagMaxSize())
	c.tag.SetTTL(config.GetCacheGTSTagTTL(), true)
	c.tag.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initTombstone() {
	c.tombstone = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "URI"},
	}, func(t1 *gtsmodel.Tombstone) *gtsmodel.Tombstone {
		t2 := new(gtsmodel.Tombstone)
		*t2 = *t1
		return t2
	}, config.GetCacheGTSTombstoneMaxSize())
	c.tombstone.SetTTL(config.GetCacheGTSTombstoneTTL(), true)
	c.tombstone.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initUser() {
	c.user = result.New([]result.Lookup{
		{Name: "ID"},
		{Name: "AccountID"},
		{Name: "Email"},
		{Name: "ConfirmationToken"},
		{Name: "ExternalID"},
	}, func(u1 *gtsmodel.User) *gtsmodel.User {
		u2 := new(gtsmodel.User)
		*u2 = *u1
		return u2
	}, config.GetCacheGTSUserMaxSize())
	c.user.SetTTL(config.GetCacheGTSUserTTL(), true)
	c.user.IgnoreErrors(ignoreErrors)
}

func (c *GTSCaches) initWebfinger() {
	c.webfinger = ttl.New[string, string](
		0,
		config.GetCacheGTSWebfingerMaxSize(),
		config.GetCacheGTSWebfingerTTL(),
	)
}
