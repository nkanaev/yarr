package nostr

const (
	KindProfileMetadata          int = 0
	KindTextNote                 int = 1
	KindRecommendServer          int = 2
	KindFollowList               int = 3
	KindEncryptedDirectMessage   int = 4
	KindDeletion                 int = 5
	KindRepost                   int = 6
	KindReaction                 int = 7
	KindBadgeAward               int = 8
	KindSimpleGroupChatMessage   int = 9
	KindSimpleGroupThreadedReply int = 10
	KindSimpleGroupThread        int = 11
	KindSimpleGroupReply         int = 12
	KindSeal                     int = 13
	KindDirectMessage            int = 14
	KindGenericRepost            int = 16
	KindReactionToWebsite        int = 17
	KindChannelCreation          int = 40
	KindChannelMetadata          int = 41
	KindChannelMessage           int = 42
	KindChannelHideMessage       int = 43
	KindChannelMuteUser          int = 44
	KindChess                    int = 64
	KindMergeRequests            int = 818
	KindBid                      int = 1021
	KIndBidConfirmation          int = 1022
	KindOpenTimestamps           int = 1040
	KindGiftWrap                 int = 1059
	KindFileMetadata             int = 1063
	KindLiveChatMessage          int = 1311
	KindPatch                    int = 1617
	KindIssue                    int = 1621
	KindReply                    int = 1622
	KindStatusOpen               int = 1630
	KindStatusApplied            int = 1631
	KindStatusClosed             int = 1632
	KindStatusDraft              int = 1633
	KindProblemTracker           int = 1971
	KindReporting                int = 1984
	KindLabel                    int = 1985
	KindRelayReviews             int = 1986
	KindAIEmbeddings             int = 1987
	KindTorrent                  int = 2003
	KindTorrentComment           int = 2004
	KindCoinjoinPool             int = 2022
	KindCommunityPostApproval    int = 4550
	KindJobFeedback              int = 7000
	KindSimpleGroupPutUser       int = 9000
	KindSimpleGroupRemoveUser    int = 9001
	KindSimpleGroupEditMetadata  int = 9002
	KindSimpleGroupDeleteEvent   int = 9005
	KindSimpleGroupCreateGroup   int = 9007
	KindSimpleGroupDeleteGroup   int = 9008
	KindSimpleGroupCreateInvite  int = 9009
	KindSimpleGroupJoinRequest   int = 9021
	KindSimpleGroupLeaveRequest  int = 9022
	KindZapGoal                  int = 9041
	KindTidalLogin               int = 9467
	KindZapRequest               int = 9734
	KindZap                      int = 9735
	KindHighlights               int = 9802
	KindMuteList                 int = 10000
	KindPinList                  int = 10001
	KindRelayListMetadata        int = 10002
	KindBookmarkList             int = 10003
	KindCommunityList            int = 10004
	KindPublicChatList           int = 10005
	KindBlockedRelayList         int = 10006
	KindSearchRelayList          int = 10007
	KindSimpleGroupList          int = 10009
	KindInterestList             int = 10015
	KindEmojiList                int = 10030
	KindDMRelayList              int = 10050
	KindUserServerList           int = 10063
	KindFileStorageServerList    int = 10096
	KindGoodWikiAuthorList       int = 10101
	KindGoodWikiRelayList        int = 10102
	KindNWCWalletInfo            int = 13194
	KindLightningPubRPC          int = 21000
	KindClientAuthentication     int = 22242
	KindNWCWalletRequest         int = 23194
	KindNWCWalletResponse        int = 23195
	KindNostrConnect             int = 24133
	KindBlobs                    int = 24242
	KindHTTPAuth                 int = 27235
	KindCategorizedPeopleList    int = 30000
	KindCategorizedBookmarksList int = 30001
	KindRelaySets                int = 30002
	KindBookmarkSets             int = 30003
	KindCuratedSets              int = 30004
	KindCuratedVideoSets         int = 30005
	KindMuteSets                 int = 30007
	KindProfileBadges            int = 30008
	KindBadgeDefinition          int = 30009
	KindInterestSets             int = 30015
	KindStallDefinition          int = 30017
	KindProductDefinition        int = 30018
	KindMarketplaceUI            int = 30019
	KindProductSoldAsAuction     int = 30020
	KindArticle                  int = 30023
	KindDraftArticle             int = 30024
	KindEmojiSets                int = 30030
	KindModularArticleHeader     int = 30040
	KindModularArticleContent    int = 30041
	KindReleaseArtifactSets      int = 30063
	KindApplicationSpecificData  int = 30078
	KindLiveEvent                int = 30311
	KindUserStatuses             int = 30315
	KindClassifiedListing        int = 30402
	KindDraftClassifiedListing   int = 30403
	KindRepositoryAnnouncement   int = 30617
	KindRepositoryState          int = 30618
	KindSimpleGroupMetadata      int = 39000
	KindSimpleGroupAdmins        int = 39001
	KindSimpleGroupMembers       int = 39002
	KindSimpleGroupRoles         int = 39003
	KindWikiArticle              int = 30818
	KindRedirects                int = 30819
	KindFeed                     int = 31890
	KindDateCalendarEvent        int = 31922
	KindTimeCalendarEvent        int = 31923
	KindCalendar                 int = 31924
	KindCalendarEventRSVP        int = 31925
	KindHandlerRecommendation    int = 31989
	KindHandlerInformation       int = 31990
	KindVideoEvent               int = 34235
	KindShortVideoEvent          int = 34236
	KindVideoViewEvent           int = 34237
	KindCommunityDefinition      int = 34550
)

func IsRegularKind(kind int) bool {
	return kind < 10000 && kind != 0 && kind != 3
}

func IsReplaceableKind(kind int) bool {
	return kind == 0 || kind == 3 || (10000 <= kind && kind < 20000)
}

func IsEphemeralKind(kind int) bool {
	return 20000 <= kind && kind < 30000
}

func IsAddressableKind(kind int) bool {
	return 30000 <= kind && kind < 40000
}
