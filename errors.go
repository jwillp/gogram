// Copyright (c) 2023 RoseLoverX

package gogram

import (
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"

	"github.com/jwillp/gogram/internal/mtproto/objects"
)

type ErrResponseCode struct {
	Code           int
	Message        string
	Description    string
	AdditionalInfo any // some errors has additional data like timeout seconds, dc id etc.
}

func RpcErrorToNative(r *objects.RpcError) error {
	nativeErrorName, additionalData := TryExpandError(r.ErrorMessage)

	desc, ok := errorMessages[nativeErrorName]
	if !ok {
		desc = nativeErrorName
	}

	if additionalData != nil {
		desc = fmt.Sprintf(desc, additionalData)
	}

	return &ErrResponseCode{
		Code:           int(r.ErrorCode),
		Message:        nativeErrorName,
		Description:    desc,
		AdditionalInfo: additionalData,
	}
}

type prefixSuffix struct {
	prefix string
	suffix string
	kind   reflect.Kind // int string bool etc.
}

var specificErrors = []prefixSuffix{
	{"EMAIL_UNCONFIRMED_", "", reflect.Int},
	{"FILE_MIGRATE_", "", reflect.Int},
	{"FILE_PART_", "_MISSING", reflect.Int},
	{"FLOOD_TEST_PHONE_WAIT_", "", reflect.Int},
	{"FLOOD_WAIT_", "", reflect.Int},
	{"INTERDC_", "_CALL_ERROR", reflect.Int},
	{"INTERDC_", "_CALL_RICH_ERROR", reflect.Int},
	{"NETWORK_MIGRATE_", "", reflect.Int},
	{"PASSWORD_TOO_FRESH_", "", reflect.Int},
	{"PHONE_MIGRATE_", "", reflect.Int},
	{"SESSION_TOO_FRESH_", "", reflect.Int},
	{"SLOWMODE_WAIT_", "", reflect.Int},
	{"STATS_MIGRATE_", "", reflect.Int},
	{"TAKEOUT_INIT_DELAY_", "", reflect.Int},
	{"USER_MIGRATE_", "", reflect.Int},
	{"PREVIOUS_CHAT_IMPORT_ACTIVE_WAIT_", "MIN", reflect.Int},
}

func TryExpandError(errStr string) (nativeErrorName string, additionalData any) {
	var choosedPrefixSuffix *prefixSuffix

	for _, errCase := range specificErrors {
		if strings.HasPrefix(errStr, errCase.prefix) && strings.HasSuffix(errStr, errCase.suffix) {
			choosedPrefixSuffix = &errCase //nolint:gosec cause we need nil if not found
			break
		}
	}

	if choosedPrefixSuffix == nil {
		return errStr, nil // common error, returning
	}

	nativeErrorName = choosedPrefixSuffix.prefix + "X" + choosedPrefixSuffix.suffix
	trimmedData := strings.TrimSuffix(strings.TrimPrefix(errStr, choosedPrefixSuffix.prefix), choosedPrefixSuffix.suffix)

	switch v := choosedPrefixSuffix.kind; v {
	case reflect.Int:
		var err error
		additionalData, err = strconv.Atoi(trimmedData)
		if err != nil {
			log.Printf("failed to parse %s as int: %s", trimmedData, err.Error())
		}

	case reflect.String:
		additionalData = trimmedData

	default:
		panic("couldn't parse this type: " + v.String())
	}

	return nativeErrorName, additionalData
}

func (e *ErrResponseCode) Error() string {
	return fmt.Sprintf("[%s] %s (code %d)", e.Message, e.Description, e.Code)
}

// gathered all errors from all methods. don't have reference in docs at all
var errorMessages = map[string]string{
	"ABOUT_TOO_LONG":                      "The provided bio is too long",
	"ACCESS_TOKEN_EXPIRED":                "Bot token expired",
	"ACCESS_TOKEN_INVALID":                "The provided token is not valid",
	"ACTIVE_USER_REQUIRED":                "The method is only available to already activated users",
	"ADMINS_TOO_MUCH":                     "Too many admins",
	"ADMIN_RANK_EMOJI_NOT_ALLOWED":        "Emoji are not allowed in admin titles or ranks",
	"ADMIN_RANK_INVALID":                  "The given admin title or rank was invalid (possibly larger than 16 characters)",
	"ALBUM_PHOTOS_TOO_MANY":               "Too many photos were included in the album",
	"API_ID_INVALID":                      "The api_id/api_hash combination is invalid",
	"API_ID_PUBLISHED_FLOOD":              "This API id was published somewhere, you can't use it now",
	"ARTICLE_TITLE_EMPTY":                 "The title of the article is empty",
	"AUDIO_TITLE_EMPTY":                   "The title attribute of the audio must be non-empty",
	"AUDIO_CONTENT_URL_EMPTY":             "The content URL of the audio must be non-empty",
	"AUTH_BYTES_INVALID":                  "The provided authorization is invalid",
	"AUTH_KEY_DUPLICATED":                 "The authorization key (session file) was used under two different IP addresses simultaneously, and can no longer be used. Use the same session exclusively, or use different sessions",
	"AUTH_KEY_INVALID":                    "The Authorization Key is invalid",
	"AUTH_KEY_PERM_EMPTY":                 "The method is unavailable for temporary authorization key, not bound to permanent",
	"AUTH_KEY_UNREGISTERED":               "The key is not registered in the system",
	"AUTH_RESTART":                        "Restart the authorization process",
	"AUTH_TOKEN_ALREADY_ACCEPTED":         "The authorization token was already used",
	"AUTH_TOKEN_EXPIRED":                  "The provided authorization token has expired and the updated QR-code must be re-scanned",
	"AUTH_TOKEN_INVALID":                  "An invalid authorization token was provided",
	"AUTOARCHIVE_NOT_AVAILABLE":           "You cannot use this feature yet",
	"BANK_CARD_NUMBER_INVALID":            "Incorrect credit card number",
	"BASE_PORT_LOC_INVALID":               "Base port location invalid",
	"BANNED_RIGHTS_INVALID":               "You cannot use that set of permissions in this request, i.e. restricting view_messages as a default",
	"BOTS_TOO_MUCH":                       "There are too many bots in this chat/channel",
	"BOT_ONESIDE_NOT_AVAIL":               "One-sided pinned messages are not available in this chat",
	"BOT_CHANNELS_NA":                     "Bots can't edit admin privileges",
	"BOT_COMMAND_DESCRIPTION_INVALID":     "The command description was empty, too long or had invalid characters used",
	"BOT_COMMAND_INVALID":                 "The bot command was in invalid format",
	"BOT_DOMAIN_INVALID":                  "The domain used for the auth button does not match the one configured in @BotFather",
	"BOT_GAMES_DISABLED":                  "Bot games cannot be used in this type of chat",
	"BOT_GROUPS_BLOCKED":                  "This bot can't be added to groups",
	"BOT_INLINE_DISABLED":                 "This bot can't be used in inline mode",
	"BOT_INVALID":                         "This is not a valid bot",
	"BOT_METHOD_INVALID":                  "The API access for bot users is restricted. The method you tried to invoke cannot be executed as a bot",
	"BOT_MISSING":                         "This method can only be run by a bot",
	"BOT_PAYMENTS_DISABLED":               "This method can only be run by a bot",
	"BOT_POLLS_DISABLED":                  "You cannot create polls under a bot account",
	"BOT_RESPONSE_TIMEOUT":                "The bot did not answer to the callback query in time",
	"BOT_SCORE_NOT_MODIFIED":              "There are too many bots in this chat/channel",
	"BROADCAST_CALLS_DISABLED":            "",
	"BROADCAST_FORBIDDEN":                 "The request cannot be used in broadcast channels",
	"BROADCAST_ID_INVALID":                "The channel is invalid",
	"BROADCAST_PUBLIC_VOTERS_FORBIDDEN":   "You cannot broadcast polls where the voters are public",
	"BROADCAST_REQUIRED":                  "The request can only be used with a broadcast channel",
	"BUTTON_DATA_INVALID":                 "The provided button data is invalid",
	"BUTTON_TYPE_INVALID":                 "The type of one of the buttons you provided is invalid",
	"BUTTON_URL_INVALID":                  "Button URL invalid",
	"CALL_ALREADY_ACCEPTED":               "The call was already accepted",
	"CALL_ALREADY_DECLINED":               "The call was already declined",
	"CALL_OCCUPY_FAILED":                  "The call failed because the user is already making another call",
	"CALL_PEER_INVALID":                   "The provided call peer object is invalid",
	"CALL_PROTOCOL_FLAGS_INVALID":         "Call protocol flags invalid",
	"CDN_METHOD_INVALID":                  "This method cannot be invoked on a CDN server. Refer to https://core.telegram.org/cdn#schema for available methods",
	"CHANNELS_ADMIN_PUBLIC_TOO_MUCH":      "You're admin of too many public channels, make some channels private to change the username of this channel",
	"CHANNELS_ADMIN_LOCATED_TOO_MUCH":     "Returned if both the check_limit and the by_location flags are set and the user has reached the limit of public geogroups",
	"CHANNELS_TOO_MUCH":                   "You have joined too many channels/supergroups",
	"CHANNEL_ADD_INVALID":                 "",
	"CHANNEL_BANNED":                      "The channel is banned",
	"CHANNEL_INVALID":                     "Invalid channel object. Make sure to pass the right types, for instance making sure that the request is designed for channels or otherwise look for a different one more suited",
	"CHANNEL_PRIVATE":                     "The channel specified is private and you lack permission to access it. Another reason may be that you were banned from it",
	"CHANNEL_PUBLIC_GROUP_NA":             "channel/supergroup not available",
	"CHANNEL_TOO_LARGE":                   "Channel is too large to be deleted; this error is issued when trying to delete channels with more than 1000 members (subject to change)",
	"CHAT_ABOUT_NOT_MODIFIED":             "About text has not changed",
	"CHAT_ABOUT_TOO_LONG":                 "Chat about too long",
	"CHAT_ADMIN_INVITE_REQUIRED":          "You do not have the rights to do this",
	"CHAT_ADMIN_REQUIRED":                 "Chat admin privileges are required to do that in the specified chat (for example, to send a message in a channel which is not yours), or invalid permissions used for the channel or group",
	"CHAT_FORBIDDEN":                      "You cannot write in this chat",
	"CHAT_FORWARDS_RESTRICTED":            "You cannot forward messages from this chat",
	"CHAT_ID_EMPTY":                       "The provided chat ID is empty",
	"CHAT_ID_INVALID":                     "Invalid object ID for a chat. Make sure to pass the right types, for instance making sure that the request is designed for chats (not channels/megagroups) or otherwise look for a different one more suited\nAn example working with a megagroup and AddChatUserRequest, it will fail because megagroups are channels. Use InviteToChannelRequest instead",
	"CHAT_INVALID":                        "The chat is invalid for this request",
	"CHAT_LINK_EXISTS":                    "The chat is linked to a channel and cannot be used in that request",
	"CHAT_NOT_MODIFIED":                   "The chat or channel wasn't modified (title, invites, username, admins, etc. are the same)",
	"CHAT_RESTRICTED":                     "The chat is restricted and cannot be used in that request",
	"CHAT_SEND_GIFS_FORBIDDEN":            "You can't send gifs in this chat",
	"CHAT_SEND_INLINE_FORBIDDEN":          "You cannot send inline results in this chat",
	"CHAT_SEND_MEDIA_FORBIDDEN":           "You can't send media in this chat",
	"CHAT_SEND_STICKERS_FORBIDDEN":        "You can't send stickers in this chat",
	"CHAT_TITLE_EMPTY":                    "No chat title provided",
	"CHAT_TOO_BIG":                        "The chat is too big",
	"CHAT_WRITE_FORBIDDEN":                "You can't write in this chat",
	"CHP_CALL_FAIL":                       "The statistics cannot be retrieved at this time",
	"CODE_EMPTY":                          "The provided code is empty",
	"CODE_HASH_INVALID":                   "Code hash invalid",
	"CODE_INVALID":                        "Code invalid (i.e. from email)",
	"CONNECTION_API_ID_INVALID":           "The provided API id is invalid",
	"CONNECTION_DEVICE_MODEL_EMPTY":       "Device model empty",
	"CONNECTION_LANG_PACK_INVALID":        "The specified language pack is not valid. This is meant to be used by official applications only so far, leave it empty",
	"CONNECTION_LAYER_INVALID":            "The very first request must always be InvokeWithLayerRequest",
	"CONNECTION_NOT_INITED":               "Connection not initialized",
	"CONNECTION_SYSTEM_EMPTY":             "Connection system empty",
	"CONNECTION_SYSTEM_LANG_CODE_EMPTY":   "The system language string was empty during connection",
	"CONTACT_ID_INVALID":                  "The provided contact ID is invalid",
	"CONTACT_NAME_EMPTY":                  "The provided contact name cannot be empty",
	"CURRENCY_TOTAL_AMOUNT_INVALID":       "The total amount of the invoice is invalid",
	"DATA_INVALID":                        "Encrypted data invalid",
	"DATA_JSON_INVALID":                   "The provided JSON data is invalid",
	"DATE_EMPTY":                          "Date empty",
	"DC_ID_INVALID":                       "This occurs when an authorization is tried to be exported for the same data center one is currently connected to",
	"DH_G_A_INVALID":                      "g_a invalid",
	"DOCUMENT_INVALID":                    "The document file was invalid and can't be used in inline mode",
	"EMAIL_HASH_EXPIRED":                  "The email hash expired and cannot be used to verify it",
	"EMAIL_INVALID":                       "The given email is invalid",
	"EMOJI_INVALID":                       "The specified emoji cannot be used or was not a emoji",
	"EMOJI_NOT_MODIFIED":                  "The emoji was not modified",
	"EMOTICON_EMPTY":                      "The emoticon field cannot be empty",
	"EMOTICON_INVALID":                    "The specified emoticon cannot be used or was not a emoticon",
	"EMOTICON_STICKERPACK_MISSING":        "The emoticon sticker pack you are trying to get is missing",
	"ENCRYPTED_MESSAGE_INVALID":           "Encrypted message invalid",
	"ENCRYPTION_ALREADY_ACCEPTED":         "Secret chat already accepted",
	"ENCRYPTION_ALREADY_DECLINED":         "The secret chat was already declined",
	"ENCRYPTION_DECLINED":                 "The secret chat was declined",
	"ENCRYPTION_ID_INVALID":               "The provided secret chat ID is invalid",
	"ENCRYPTION_OCCUPY_FAILED":            "TDLib developer claimed it is not an error while accepting secret chats and 500 is used instead of 420",
	"ENTITIES_TOO_LONG":                   "It is no longer possible to send such long data inside entity tags (for example inline text URLs)",
	"ENTITY_MENTION_USER_INVALID":         "You can't use this entity",
	"ERROR_TEXT_EMPTY":                    "The provided error message is empty",
	"EXPIRE_DATE_INVALID":                 "The provided expire date is invalid",
	"EXPIRE_FORBIDDEN":                    "The provided expire date is forbidden",
	"EXPORT_CARD_INVALID":                 "Provided card is invalid",
	"EXTERNAL_URL_INVALID":                "External URL invalid",
	"FIELD_NAME_EMPTY":                    "The field with the name FIELD_NAME is missing",
	"FIELD_NAME_INVALID":                  "The field with the name FIELD_NAME is invalid",
	"FILEREF_UPGRADE_NEEDED":              "The file reference needs to be refreshed before being used again",
	"FILE_CONTENT_TYPE_INVALID":           "The provided file content type is invalid",
	"FILE_ID_INVALID":                     "The provided file id is invalid. Make sure all parameters are present, have the correct type and are not empty (ID, access hash, file reference, thumb size ...)",
	"FILE_PARTS_INVALID":                  "The number of file parts is invalid",
	"FILE_PART_EMPTY":                     "The provided file part is empty",
	"FILE_PART_INVALID":                   "The file part number is invalid",
	"FILE_PART_LENGTH_INVALID":            "The length of a file part is invalid",
	"FILE_PART_SIZE_CHANGED":              "The file part size (chunk size) cannot change during upload",
	"FILE_PART_SIZE_INVALID":              "The provided file part size is invalid",
	"FILE_REFERENCE_EMPTY":                "The file reference must exist to access the media and it cannot be empty",
	"FILE_REFERENCE_EXPIRED":              "The file reference has expired and is no longer valid or it belongs to self-destructing media and cannot be resent",
	"FILE_REFERENCE_INVALID":              "The file reference is invalid or you can't do that operation on such message",
	"FILE_TITLE_EMPTY":                    "The provided file title is empty",
	"FIRSTNAME_INVALID":                   "The first name is invalid",
	"FILTER_NOT_SUPPORTED":                "The provided filter is not supported",
	"FOLDER_ID_EMPTY":                     "The folder you tried to delete was already empty",
	"FOLDER_ID_INVALID":                   "The folder you tried to use was not valid",
	"FRESH_CHANGE_ADMINS_FORBIDDEN":       "Recently logged-in users cannot add or change admins",
	"FRESH_CHANGE_PHONE_FORBIDDEN":        "Recently logged-in users cannot use this request",
	"FRESH_RESET_AUTHORISATION_FORBIDDEN": "The current session is too new and cannot be used to reset other authorisations yet",
	"FROM_PEER_INVALID":                   "The given from_user peer cannot be used for the parameter",
	"GAME_BOT_INVALID":                    "You cannot send that game with the current bot",
	"GIF_CONTENT_TYPE_INVALID":            "The provided GIF content type is invalid",
	"GIF_ID_INVALID":                      "The provided GIF ID is invalid",
	"GRAPH_INVALID_RELOAD":                "Data can't be used for the channel statistics, graphs invalid",
	"GRAPH_OUTDATED_RELOAD":               "Data can't be used for the channel statistics, graphs outdated",
	"GROUPCALL_ADD_PARTICIPANTS_FAILED":   "Failed to add participants to the group call",
	"GROUPCALL_ALREADY_DISCARDED":         "The group call was already discarded",
	"GROUPCALL_FORBIDDEN":                 "The group call is forbidden",
	"GROUPCALL_INVALID":                   "The group call is invalid",
	"GROUPCALL_JOIN_MISSING":              "The group call join is missing",
	"GROUPCALL_SSRC_DUPLICATE_MUCH":       "Too many duplicate SSRCs",
	"GROUPCALL_NOT_MODIFIED":              "The group call was not modified",
	"GROUPED_MEDIA_INVALID":               "Invalid grouped media",
	"GROUP_CALL_INVALID":                  "Group call invalid",
	"HASH_INVALID":                        "The provided hash is invalid",
	"HISTORY_GET_FAILED":                  "Fetching of history failed",
	"IMAGE_PROCESS_FAILED":                "Failure while processing image",
	"IMPORT_FILE_INVALID":                 "The file is too large to be imported",
	"IMPORT_FORMAT_UNRECOGNIZED":          "Unknown import format",
	"IMPORT_ID_INVALID":                   "The provided import ID is invalid",
	"INLINE_BOT_REQUIRED":                 "The action must be performed through an inline bot callback",
	"INLINE_RESULT_EXPIRED":               "The inline query expired",
	"INPUT_CONSTRUCTOR_INVALID":           "The provided constructor is invalid",
	"INPUT_FETCH_ERROR":                   "An error occurred while deserializing TL parameters",
	"INPUT_FETCH_FAIL":                    "Failed deserializing TL payload",
	"INPUT_FILTER_INVALID":                "The search query filter is invalid",
	"INPUT_LAYER_INVALID":                 "The provided layer is invalid",
	"INPUT_METHOD_INVALID":                "The invoked method does not exist anymore or has never existed",
	"INPUT_REQUEST_TOO_LONG":              "The input request was too long. This may be a bug in the library as it can occur when serializing more bytes than it should (like appending the vector constructor code at the end of a message)",
	"INPUT_USER_DEACTIVATED":              "The specified user was deleted",
	"INVITE_FORBIDDEN_WITH_JOINAS":        "You cannot use this method to invite users to a channel if you are using the join_as parameter",
	"INVITE_HASH_EMPTY":                   "The invite hash is empty",
	"INVITE_HASH_EXPIRED":                 "The chat the user tried to join has expired and is not valid anymore",
	"INVITE_HASH_INVALID":                 "The invite hash is invalid",
	"LANG_CODE_INVALID":                   "The provided language code is invalid",
	"LANG_PACK_INVALID":                   "The provided language pack is invalid",
	"LASTNAME_INVALID":                    "The last name is invalid",
	"LIMIT_INVALID":                       "An invalid limit was provided. See https://core.telegram.org/api/files#downloading-files",
	"LINK_NOT_MODIFIED":                   "The channel is already linked to this group",
	"LOCATION_INVALID":                    "The location given for a file was invalid. See https://core.telegram.org/api/files#downloading-files",
	"MAX_ID_INVALID":                      "The provided max ID is invalid",
	"MAX_QTS_INVALID":                     "The provided QTS were invalid",
	"MD5_CHECKSUM_INVALID":                "The MD5 check-sums do not match",
	"MEDIA_CAPTION_TOO_LONG":              "The caption is too long",
	"MEDIA_EMPTY":                         "The provided media object is invalid or the current account may not be able to send it (such as games as users)",
	"MEDIA_GROUPED_INVALID":               "The provided media group is invalid",
	"MEDIA_INVALID":                       "Media invalid",
	"MEDIA_NEW_INVALID":                   "The new media to edit the message with is invalid (such as stickers or voice notes)",
	"MEDIA_PREV_INVALID":                  "The old media cannot be edited with anything else (such as stickers or voice notes)",
	"MEDIA_TTL_INVALID":                   "The provided media TTL is invalid",
	"MEGAGROUP_ID_INVALID":                "The group is invalid",
	"MEGAGROUP_PREHISTORY_HIDDEN":         "You can't set this discussion group because it's history is hidden",
	"MEGAGROUP_REQUIRED":                  "The request can only be used with a megagroup channel",
	"MEMBER_NO_LOCATION":                  "An internal failure occurred while fetching user info (couldn't find location)",
	"MEMBER_OCCUPY_PRIMARY_LOC_FAILED":    "Occupation of primary member location failed",
	"MESSAGE_AUTHOR_REQUIRED":             "Message author required",
	"MESSAGE_DELETE_FORBIDDEN":            "You can't delete one of the messages you tried to delete, most likely because it is a service message.",
	"MESSAGE_EDIT_TIME_EXPIRED":           "You can't edit this message anymore, too much time has passed since its creation.",
	"MESSAGE_EMPTY":                       "Empty or invalid UTF-8 message was sent",
	"MESSAGE_IDS_EMPTY":                   "No message ids were provided",
	"MESSAGE_ID_INVALID":                  "The specified message ID is invalid or you can't do that operation on such message",
	"MESSAGE_NOT_MODIFIED":                "Content of the message was not modified",
	"MESSAGE_POLL_CLOSED":                 "The poll was closed and can no longer be voted on",
	"MESSAGE_TOO_LONG":                    "Message was too long. Current maximum length is 4096 UTF-8 characters",
	"METHOD_INVALID":                      "The API method is invalid and cannot be used",
	"MSGID_DECREASE_RETRY":                "The request should be retried with a lower message ID",
	"MSG_ID_INVALID":                      "The message ID used in the peer was invalid",
	"MSG_WAIT_FAILED":                     "A waiting call returned an error",
	"MT_SEND_QUEUE_TOO_LONG":              "The message was not sent because the send queue is too long",
	"MULTI_MEDIA_TOO_LONG":                "Too many media files were included in the same album",
	"NEED_CHAT_INVALID":                   "The provided chat is invalid",
	"NEED_MEMBER_INVALID":                 "The provided member is invalid or does not exist (for example a thumb size)",
	"NEW_SALT_INVALID":                    "The new salt is invalid",
	"NEW_SETTINGS_INVALID":                "The new settings are invalid",
	"NEXT_OFFSET_INVALID":                 "The value for next_offset is invalid. Check that it has normal characters and is not too long",
	"OFFSET_INVALID":                      "The given offset was invalid, it must be divisible by 1KB. See https://core.telegram.org/api/files#downloading-files",
	"OFFSET_PEER_ID_INVALID":              "The provided offset peer is invalid",
	"OPTIONS_TOO_MUCH":                    "You defined too many options for the poll",
	"OPTION_INVALID":                      "The option specified is invalid and does not exist in the target poll",
	"PACK_SHORT_NAME_INVALID":             `Invalid sticker pack name. It must begin with a letter, can't contain consecutive underscores and must end in "_by_<bot username>".`,
	"PACK_SHORT_NAME_OCCUPIED":            "A stickerpack with this name already exists",
	"PARTICIPANTS_TOO_FEW":                "Not enough participants",
	"PARTICIPANT_CALL_FAILED":             "Failure while making call",
	"PARTICIPANT_JOIN_MISSING":            "The participant is not in the chat",
	"PARTICIPANT_ID_INVALID":              "The participant ID is invalid",
	"PARTICIPANT_VERSION_OUTDATED":        "The other participant does not use an up to date telegram client with support for calls",
	"PASSWORD_EMPTY":                      "The provided password is empty",
	"PASSWORD_HASH_INVALID":               "The password (and thus its hash value) you entered is invalid",
	"PASSWORD_MISSING":                    "The account must have 2-factor authentication enabled (a password) before this method can be used",
	"PASSWORD_RECOVERY_EXPIRED":           "The password recovery code has expired",
	"PASSWORD_RECOVERY_NA":                "The password recovery code is invalid",
	"PASSWORD_REQUIRED":                   "The account must have 2-factor authentication enabled (a password) before this method can be used",
	"PAYMENT_PROVIDER_INVALID":            "The payment provider was not recognised or its token was invalid",
	"PEER_FLOOD":                          "Too many requests",
	"PEER_ID_INVALID":                     "An invalid Peer was used. Make sure to pass the right peer type and that the value is valid (for instance, bots cannot start conversations)",
	"PEER_ID_NOT_SUPPORTED":               "The provided peer ID is not supported",
	"PERSISTENT_TIMESTAMP_EMPTY":          "Persistent timestamp empty",
	"PERSISTENT_TIMESTAMP_INVALID":        "Persistent timestamp invalid",
	"PERSISTENT_TIMESTAMP_OUTDATED":       "Persistent timestamp outdated",
	"PHONE_CODE_EMPTY":                    "The phone code is missing",
	"PHONE_CODE_EXPIRED":                  "The confirmation code has expired",
	"PHONE_CODE_HASH_EMPTY":               "The phone code hash is missing",
	"PHONE_CODE_INVALID":                  "The phone code entered was invalid",
	"PHONE_NOT_OCCUPIED":                  "The phone number is not yet occupied",
	"PHONE_NUMBER_APP_SIGNUP_FORBIDDEN":   "You can't sign up using this app",
	"PHONE_NUMBER_BANNED":                 "The used phone number has been banned from Telegram and cannot be used anymore. Maybe check https://www.telegram.org/faq_spam",
	"PHONE_NUMBER_FLOOD":                  "You asked for the code too many times.",
	"PHONE_NUMBER_INVALID":                "The phone number is invalid",
	"PHONE_NUMBER_OCCUPIED":               "The phone number is already in use",
	"PHONE_NUMBER_UNOCCUPIED":             "The phone number is not yet being used",
	"PHONE_PASSWORD_FLOOD":                "You have tried logging in too many times",
	"PHONE_PASSWORD_PROTECTED":            "This phone is password protected",
	"PHOTO_CONTENT_TYPE_INVALID":          "The content type of photo is invalid",
	"PHOTO_CONTENT_URL_EMPTY":             "The content from the URL used as a photo appears to be empty or has caused another HTTP error",
	"PHOTO_CROP_SIZE_SMALL":               "Photo is too small",
	"PHOTO_EXT_INVALID":                   "The extension of the photo is invalid",
	"PHOTO_FILE_MISSING":                  "The photo file is missing",
	"PHOTO_ID_INVALID":                    "Photo id is invalid",
	"PHOTO_INVALID":                       "Photo invalid",
	"PHOTO_INVALID_DIMENSIONS":            "The photo dimensions are invalid (hint: `pip install pillow` for `send_file` to resize images)",
	"PHOTO_SAVE_FILE_INVALID":             "The photo you tried to send cannot be saved by Telegram. A reason may be that it exceeds 10MB. Try resizing it locally",
	"PHOTO_THUMB_URL_EMPTY":               "The URL used as a thumbnail appears to be empty or has caused another HTTP error",
	"PIN_RESTRICTED":                      "You can't pin messages in private chats with other people",
	"PINNED_DIALOGS_TOO_MUCH":             "You have too many pinned dialogs",
	"POLL_ANSWER_INVALID":                 "The poll answer is invalid",
	"POLL_ANSWERS_INVALID":                "The poll did not have enough answers or had too many",
	"POLL_OPTION_DUPLICATE":               "A duplicate option was sent in the same poll",
	"POLL_OPTION_INVALID":                 "A poll option used invalid data (the data may be too long)",
	"POLL_QUESTION_INVALID":               "The poll question was either empty or too long",
	"POLL_UNSUPPORTED":                    "This layer does not support polls in the issued method",
	"POLL_VOTE_REQUIRED":                  "You must vote in a non-anonymous poll before requesting the results",
	"PRIVACY_KEY_INVALID":                 "The privacy key is invalid",
	"PRIVACY_TOO_LONG":                    "Cannot add that many entities in a single request",
	"PRIVACY_VALUE_INVALID":               "The privacy value is invalid",
	"PTS_CHANGE_EMPTY":                    "No PTS change",
	"PUBLIC_KEY_REQUIRED":                 "The public key is required",
	"QUERY_ID_EMPTY":                      "The query ID is empty",
	"QUERY_ID_INVALID":                    "The query ID is invalid",
	"QUERY_TOO_SHORT":                     "The query string is too short",
	"QUIZ_ANSWER_MISSING":                 "A quiz must specify at least one answer",
	"QUIZ_CORRECT_ANSWERS_EMPTY":          "A quiz must specify one correct answer",
	"QUIZ_CORRECT_ANSWERS_TOO_MUCH":       "There can only be one correct answer",
	"QUIZ_CORRECT_ANSWER_INVALID":         "The correct answer is not an existing answer",
	"QUIZ_MULTIPLE_INVALID":               "A poll cannot be both multiple choice and quiz",
	"RANDOM_ID_DUPLICATE":                 "You provided a random ID that was already used",
	"RANDOM_ID_INVALID":                   "A provided random ID is invalid",
	"RANDOM_LENGTH_INVALID":               "Random length invalid",
	"RANGES_INVALID":                      "Invalid range provided",
	"REACTION_EMPTY":                      "No reaction provided",
	"REACTION_INVALID":                    "Invalid reaction provided (only emoji are allowed) or you cannot use the reaction in the specified chat",
	"REFLECTOR_NOT_AVAILABLE":             "Invalid call reflector server",
	"REG_ID_GENERATE_FAILED":              "Failure while generating registration ID",
	"REPLY_MARKUP_GAME_EMPTY":             "The provided reply markup for the game is empty",
	"REPLY_MARKUP_INVALID":                "The provided reply markup is invalid",
	"REPLY_MARKUP_TOO_LONG":               "The data embedded in the reply markup buttons was too much",
	"RESET_REQUEST_MISSING":               "You must request a password reset first",
	"RESULTS_TOO_MUCH":                    "You sent too many results, see https://core.telegram.org/bots/api#answerinlinequery for the current limit",
	"RESULT_ID_DUPLICATE":                 "Duplicated IDs on the sent results. Make sure to use unique IDs",
	"RESULT_ID_INVALID":                   "The given result cannot be used to send the selection to the bot",
	"RESULT_TYPE_INVALID":                 "Result type invalid",
	"REVOTE_NOT_ALLOWED":                  "You cannot change your vote in a non-anonymous poll",
	"RIGHT_FORBIDDEN":                     "Either your admin rights do not allow you to do this or you passed the wrong rights combination (some rights only apply to channels and vice versa)",
	"RPC_CALL_FAIL":                       "Telegram is having internal issues, please try again later.",
	"RPC_MCGET_FAIL":                      "Telegram is having internal issues, please try again later.",
	"RSA_DECRYPT_FAILED":                  "Internal RSA decryption failed",
	"SCHEDULE_BOT_NOT_ALLOWED":            "Bots are not allowed to schedule messages",
	"SCHEDULE_DATE_INVALID":               "The date you tried to schedule is invalid",
	"SCHEDULE_DATE_TOO_LATE":              "The date you tried to schedule is too far in the future (last known limit of 1 year and a few hours)",
	"SCHEDULE_STATUS_PRIVATE":             "You cannot schedule a message until the person comes online if their privacy does not show this information",
	"SCHEDULE_TOO_MUCH":                   "You cannot schedule more messages in this chat (last known limit of 100 per chat)",
	"SCORE_INVALID":                       "The given game score is invalid",
	"SEARCH_QUERY_EMPTY":                  "The search query is empty",
	"SECONDS_INVALID":                     "Slow mode only supports certain values (e.g. 0, 10s, 30s, 1m, 5m, 15m and 1h)",
	"SEND_AS_PEER_INVALID":                "The sendAs Peer is Invalid or cannot be used as Bot",
	"SEND_CODE_UNAVAILABLE":               "The code you tried to send is not available",
	"SEND_MESSAGE_MEDIA_INVALID":          "The message media was invalid or not specified",
	"SEND_MESSAGE_TYPE_INVALID":           "The message type is invalid",
	"SENSITIVE_CHANGE_FORBIDDEN":          "Your sensitive content settings cannot be changed at this time",
	"SESSION_EXPIRED":                     "The authorization has expired",
	"SESSION_PASSWORD_NEEDED":             "Two-steps verification is enabled and a password is required",
	"SESSION_REVOKED":                     "The authorization has been invalidated, because of the user terminating all sessions",
	"SHA256_HASH_INVALID":                 "The provided SHA256 hash is invalid",
	"SHORTNAME_OCCUPY_FAILED":             "An error occurred when trying to register the short-name used for the sticker pack. Try a different name",
	"SHORT_NAME_INVALID":                  "The sticker set short name is invalid",
	"SHORT_NAME_OCCUPIED":                 "The sticker set short name is already occupied",
	"SLOWMODE_MULTI_MSGS_DISABLED":        "You cannot send multiple messages at once in slow mode",
	"SRP_ID_INVALID":                      "The SRP ID is invalid",
	"START_PARAM_EMPTY":                   "The start parameter is empty",
	"START_PARAM_INVALID":                 "Start parameter invalid",
	"STICKERPACK_STICKERS_TOO_MUCH":       "The sticker pack contains too many stickers",
	"STICKERSET_INVALID":                  "The provided sticker set is invalid",
	"STICKERSET_OWNER_ANONYMOUS":          "This sticker set can't be used as the group's official stickers because it was created by one of its anonymous admins",
	"STICKERS_EMPTY":                      "No sticker provided",
	"STICKERS_TOO_MUCH":                   "The pack contains too many stickers",
	"STICKER_DOCUMENT_INVALID":            "The sticker file was invalid (this file has failed Telegram internal checks, make sure to use the correct format and comply with https://core.telegram.org/animated_stickers)",
	"STICKER_EMOJI_INVALID":               "Sticker emoji invalid",
	"STICKER_FILE_INVALID":                "Sticker file invalid",
	"STICKER_GIF_DIMENSIONS":              "The sticker GIF dimensions are invalid",
	"STICKER_ID_INVALID":                  "The provided sticker ID is invalid",
	"STICKER_INVALID":                     "The provided sticker is invalid",
	"STICKER_PNG_DIMENSIONS":              "Sticker png dimensions invalid",
	"STICKER_PNG_NOPNG":                   "Stickers must be a png file but the used image was not a png",
	"STICKER_TGS_NODOC":                   "Sticker TGS not uploaded as document",
	"STICKER_TGS_NOTGS":                   "Stickers must be a tgs file but the used file was not a tgs",
	"STICKER_THUMB_PNG_NOPNG":             "Stickerset thumb must be a png file but the used file was not png",
	"STICKER_THUMB_TGS_NOTGS":             "Stickerset thumb must be a tgs file but the used file was not tgs",
	"STICKER_VIDEO_NOWEBM":                "Stickers must be a webm file but the used file was not webm",
	"STICKER_VIDEO_BIG":                   "The provided sticker video is too big",
	"STORAGE_CHECK_FAILED":                "Server storage check failed",
	"STORE_INVALID_SCALAR_TYPE":           "Invalid scalar type",
	"TAKEOUT_INVALID":                     "The takeout session has been invalidated by another data export session",
	"TAKEOUT_REQUIRED":                    "You must initialize a takeout request first",
	"TEMP_AUTH_KEY_EMPTY":                 "No temporary auth key provided",
	"TIMEOUT":                             "A timeout occurred while fetching data from the worker",
	"TITLE_INVALID":                       "The sticker set title is invalid",
	"THEME_INVALID":                       "Theme invalid",
	"THEME_MIME_INVALID":                  "You cannot create this theme, the mime-type is invalid",
	"TMP_PASSWORD_DISABLED":               "The temporary password is disabled",
	"TMP_PASSWORD_INVALID":                "Password auth needs to be regenerated",
	"TOKEN_INVALID":                       "The provided token is invalid",
	"TTL_DAYS_INVALID":                    "The provided TTL is invalid",
	"TTL_MEDIA_INVALID":                   "The provided media cannot be used with a TTL",
	"TTL_PERIOD_INVALID":                  "The provided TTL Period is invalid",
	"TYPES_EMPTY":                         "The types field is empty",
	"TYPE_CONSTRUCTOR_INVALID":            "The type constructor is invalid",
	"UNKNOWN_ERROR":                       "The server has returned an unknown error",
	"UNKNOWN_METHOD":                      "The method you tried to call cannot be called on non-CDN DCs",
	"UNTIL_DATE_INVALID":                  "That date cannot be specified in this request (try using None)",
	"UPDATE_APP_TO_LOGIN":                 "This layer no longer supports logging in, please update your app",
	"URL_INVALID":                         "The URL used was invalid (e.g. when answering a callback with a URL that's not t.me/yourbot or your game's URL)",
	"USAGE_LIMIT_INVALID":                 "The invite link usage limit is invalid",
	"USER_VOLUME_INVALID":                 "The provided volume is invalid",
	"USERNAME_INVALID":                    `Nobody is using this username, or the username is unacceptable. If the latter, it must match r"[a-zA-Z][\w\d]{3,30}[a-zA-Z\d]"`,
	"USERNAME_NOT_MODIFIED":               "The username is not different from the current username",
	"USERNAME_NOT_OCCUPIED":               "The username is not in use by anyone else yet",
	"USERNAME_OCCUPIED":                   "The username is already taken",
	"USERS_TOO_FEW":                       "Not enough users (to create a chat, for example)",
	"USERS_TOO_MUCH":                      "The maximum number of users has been exceeded (to create a chat, for example)",
	"USER_ADMIN_INVALID":                  "Either you're not an admin or you tried to ban an admin that you didn't promote",
	"USER_ALREADY_INVITED":                "The user has already been invited to the group call",
	"USER_ALREADY_PARTICIPANT":            "The authenticated user is already a participant of the chat",
	"USER_BANNED_IN_CHANNEL":              "You're banned from sending messages in supergroups/channels",
	"USER_BLOCKED":                        "User blocked",
	"USER_BOT":                            "Bots can only be admins in channels.",
	"USER_BOT_INVALID":                    "This method can only be called by a bot",
	"USER_BOT_REQUIRED":                   "This method can only be called by a bot",
	"USER_CHANNELS_TOO_MUCH":              "One of the users you tried to add is already in too many channels/supergroups",
	"USER_CREATOR":                        "You can't leave this channel, because you're its creator",
	"USER_DEACTIVATED":                    "The user has been deleted/deactivated",
	"USER_DEACTIVATED_BAN":                "The user has been deleted/deactivated",
	"USER_ID_INVALID":                     "Invalid object ID for a user. Make sure to pass the right types, for instance making sure that the request is designed for users or otherwise look for a different one more suited",
	"USER_INVALID":                        "The given user was invalid",
	"USER_IS_BLOCKED":                     "User is blocked",
	"USER_IS_BOT":                         "Bots can't send messages to other bots",
	"USER_KICKED":                         "This user was kicked from this supergroup/channel",
	"USER_NOT_MUTUAL_CONTACT":             "The provided user is not a mutual contact",
	"USER_NOT_PARTICIPANT":                "The target user is not a member of the specified megagroup or channel",
	"USER_PRIVACY_RESTRICTED":             "The user's privacy settings do not allow you to do this",
	"USER_RESTRICTED":                     "You're spamreported, you can't create channels or chats.",
	"USERPIC_UPLOAD_REQUIRED":             "You must have a profile picture before using this method",
	"VIDEO_CONTENT_TYPE_INVALID":          "The video content type is not supported with the given parameters (i.e. supports_streaming)",
	"VIDEO_FILE_INVALID":                  "The given video cannot be used",
	"VIDEO_TITLE_EMPTY":                   "The video title of inline result is empty",
	"WALLPAPER_FILE_INVALID":              "The given file cannot be used as a wallpaper",
	"WALLPAPER_INVALID":                   "The input wallpaper was not valid",
	"WALLPAPER_MIME_INVALID":              "The MimeType of the wallpaper is Invalid",
	"WC_CONVERT_URL_INVALID":              "WC convert URL invalid",
	"WEBDOCUMENT_MIME_INVALID":            "The MimeType of the URL is Invalid",
	"WEBDOCUMENT_URL_INVALID":             "The given URL cannot be used",
	"WEBPAGE_CURL_FAILED":                 "Failure while fetching the webpage with cURL",
	"WEBPAGE_MEDIA_EMPTY":                 "Webpage media empty",
	"WORKER_BUSY_TOO_LONG_RETRY":          "Telegram workers are too busy to respond immediately",
	"YOU_BLOCKED_USER":                    "You blocked this user",

	// errors with additional data
	"2FA_CONFIRM_WAIT_X":                    "You'll be able to reset your account in X seconds. If not, account will be deleted in 1 week for security reasons",
	"EMAIL_UNCONFIRMED_X":                   "Email unconfirmed, the length of the code must be %v",
	"FILE_MIGRATE_X":                        "The file to be accessed is currently stored in DC %v",
	"FILE_PART_X_MISSING":                   "Part %v of the file is missing from storage",
	"FLOOD_TEST_PHONE_WAIT_X":               "A wait of %v seconds is required in the test servers",
	"FLOOD_WAIT_X":                          "A wait of %v seconds is required",
	"INTERDC_X_CALL_ERROR":                  "An error occurred while communicating with DC %v",
	"INTERDC_X_CALL_RICH_ERROR":             "A rich error occurred while communicating with DC %v",
	"NETWORK_MIGRATE_X":                     "The source IP address is associated with DC %v",
	"PASSWORD_TOO_FRESH_X":                  "The password was added too recently and %v seconds must pass before using the method",
	"PHONE_MIGRATE_X":                       "The phone number a user is trying to use for authorization is associated with DC %v",
	"SESSION_TOO_FRESH_X":                   "The session logged in too recently and %v seconds must pass before calling the method",
	"SLOWMODE_WAIT_X":                       "A wait of %v seconds is required before sending another message in this chat",
	"STATS_MIGRATE_X":                       "The channel statistics must be fetched from DC %v",
	"TAKEOUT_INIT_DELAY_X":                  "A wait of %v seconds is required before being able to initiate the takeout",
	"USER_MIGRATE_X":                        "The user whose identity is being used to execute queries is associated with DC %v",
	"PREVIOUS_CHAT_IMPORT_ACTIVE_WAIT_XMIN": "Similar to a flood wait, must wait %v minutes",

	// next errors was found, but they're too strange and looks like misspelling
	// "FILE_REFERENCE_*":                  "The file reference expired, it must be refreshed",
	//! pony floodwait https://core.telegram.org/method/messages.forwardMessages
	// "P0NY_FLOODWAIT":                    " ", //! No any description provided
	// "INPUT_METHOD_INVALID_1192227_X":    "Invalid method",
	// "INPUT_METHOD_INVALID_1400137063_X": "Invalid method",
	// "INPUT_METHOD_INVALID_1604042050_X": "Invalid method",
}

type BadMsgError struct {
	*objects.BadMsgNotification
	Description string
}

func BadMsgErrorFromNative(in *objects.BadMsgNotification) *BadMsgError {
	return &BadMsgError{
		BadMsgNotification: in,
		Description:        badMsgErrorCodes[uint8(in.Code)],
	}
}

func (e *BadMsgError) Error() string {
	return fmt.Sprintf("%v (code %v)", e.Description, e.Code)
}

// https://core.telegram.org/mtproto/service_messages_about_messages#notice-of-ignored-error-message
var badMsgErrorCodes = map[uint8]string{
	16: "msg_id too low (most likely, client time is wrong; it would be worthwhile to synchronize it using msg_id notifications and re-send the original message with the “correct” msg_id or wrap it in a container with a new msg_id if the original message had waited too long on the client to be transmitted)",
	17: "msg_id too high (similar to the previous case, the client time has to be synchronized, and the message re-sent with the correct msg_id",
	18: "incorrect two lower order msg_id bits (the server expects client message msg_id to be divisible by 4)",
	19: "container msg_id is the same as msg_id of a previously received message (this must never happen)",
	20: "message too old, and it cannot be verified whether the server has received a message with this msg_id or not",
	32: "msg_seqno too low (the server has already received a message with a lower msg_id but with either a higher or an equal and odd seqno)",
	33: "msg_seqno too high (similarly, there is a message with a higher msg_id but with either a lower or an equal and odd seqno)",
	34: "an even msg_seqno expected (irrelevant message), but odd received",
	35: "odd msg_seqno expected (relevant message), but even received",
	48: "incorrect server salt (in this case, the bad_server_salt response is received with the correct salt, and the message is to be re-sent with it)",
	64: "invalid container",
}

type BadSystemMessageCode int32

const (
	ErrBadMsgUnknown             BadSystemMessageCode = 0
	ErrBadMsgIdTooLow            BadSystemMessageCode = 16
	ErrBadMsgIdTooHigh           BadSystemMessageCode = 17
	ErrBadMsgIncorrectMsgIdBits  BadSystemMessageCode = 18
	ErrBadMsgWrongContainerMsgId BadSystemMessageCode = 19 // this must never happen
	ErrBadMsgMessageTooOld       BadSystemMessageCode = 20
	ErrBadMsgSeqNoTooLow         BadSystemMessageCode = 32
	ErrBadMsgSeqNoTooHigh        BadSystemMessageCode = 33
	ErrBadMsgSeqNoExpectedEven   BadSystemMessageCode = 34
	ErrBadMsgSeqNoExpectedOdd    BadSystemMessageCode = 35
	ErrBadMsgServerSaltIncorrect BadSystemMessageCode = 48
	ErrBadMsgInvalidContainer    BadSystemMessageCode = 64
)

// internal errors for internal purposes

type errorSessionConfigsChanged struct{}

func (*errorSessionConfigsChanged) Error() string {
	return "session configuration was changed, need to repeat request"
}

func (*errorSessionConfigsChanged) CRC() uint32 {
	return 0x00000000
}
