//go:build !js

package mysql

import vitessmysql "github.com/dolthub/vitess/go/mysql"

const (
	DefaultCachingSha2PasswordHashIterations = vitessmysql.DefaultCachingSha2PasswordHashIterations
	DefaultConnBufferSize                    = vitessmysql.DefaultConnBufferSize

	MysqlNativePassword    = vitessmysql.MysqlNativePassword
	CachingSha2Password    = vitessmysql.CachingSha2Password
	CharacterSetBinary     = vitessmysql.CharacterSetBinary
	CharacterSetUtf8       = vitessmysql.CharacterSetUtf8
	ServerInTransaction    = vitessmysql.ServerInTransaction
	ServerStatusAutocommit = vitessmysql.ServerStatusAutocommit
	ServerCursorExists     = vitessmysql.ServerCursorExists

	CapabilityClientFoundRows = vitessmysql.CapabilityClientFoundRows

	ERAccessDeniedError                             = vitessmysql.ERAccessDeniedError
	ERBadDb                                         = vitessmysql.ERBadDb
	ERBadFieldError                                 = vitessmysql.ERBadFieldError
	ERBadNullError                                  = vitessmysql.ERBadNullError
	ERBadTable                                      = vitessmysql.ERBadTable
	ERBase64DecodeError                             = vitessmysql.ERBase64DecodeError
	ERBlobKeyWithoutLength                          = vitessmysql.ERBlobKeyWithoutLength
	ERCantDropFieldOrKey                            = vitessmysql.ERCantDropFieldOrKey
	ERDbCreateExists                                = vitessmysql.ERDbCreateExists
	ERDbDropExists                                  = vitessmysql.ERDbDropExists
	ERDupEntry                                      = vitessmysql.ERDupEntry
	ERFieldSpecifiedTwice                           = vitessmysql.ERFieldSpecifiedTwice
	ERInvalidCharacterString                        = vitessmysql.ERInvalidCharacterString
	ERKeyColumnDoesNotExist                         = vitessmysql.ERKeyColumnDoesNotExist
	ERLockDeadlock                                  = vitessmysql.ERLockDeadlock
	ERMixOfGroupFuncAndFields                       = vitessmysql.ERMixOfGroupFuncAndFields
	ERMultiplePriKey                                = vitessmysql.ERMultiplePriKey
	ERNoFormatDescriptionEventBeforeBinlogStatement = vitessmysql.ERNoFormatDescriptionEventBeforeBinlogStatement
	ERNoSuchTable                                   = vitessmysql.ERNoSuchTable
	ERNotSupportedYet                               = vitessmysql.ERNotSupportedYet
	EROnlyFDAndRBREventsAllowedInBinlogStatement    = vitessmysql.EROnlyFDAndRBREventsAllowedInBinlogStatement
	EROperandColumns                                = vitessmysql.EROperandColumns
	ERRowIsReferenced2                              = vitessmysql.ERRowIsReferenced2
	ERSubqueryNo1Row                                = vitessmysql.ERSubqueryNo1Row
	ERTooLongKey                                    = vitessmysql.ERTooLongKey
	ERTruncatedWrongValue                           = vitessmysql.ERTruncatedWrongValue
	ERTruncatedWrongValueForField                   = vitessmysql.ERTruncatedWrongValueForField
	ERUnknownError                                  = vitessmysql.ERUnknownError
	ERWrongAutoKey                                  = vitessmysql.ERWrongAutoKey
	ERWrongDbName                                   = vitessmysql.ERWrongDbName
	ERWrongSubKey                                   = vitessmysql.ERWrongSubKey
	ErNoReferencedRow2                              = vitessmysql.ErNoReferencedRow2

	SSAccessDeniedError = vitessmysql.SSAccessDeniedError
	SSClientError       = vitessmysql.SSClientError
	SSLockDeadlock      = vitessmysql.SSLockDeadlock

	AuthRejected     = vitessmysql.AuthRejected
	AuthAccepted     = vitessmysql.AuthAccepted
	AuthNeedMoreData = vitessmysql.AuthNeedMoreData
)

type (
	AuthMethod            = vitessmysql.AuthMethod
	AuthMethodDescription = vitessmysql.AuthMethodDescription
	AuthServer            = vitessmysql.AuthServer
	BinlogEvent           = vitessmysql.BinlogEvent
	BinlogReplicaHandler  = vitessmysql.BinlogReplicaHandler
	BoundQuery            = vitessmysql.BoundQuery
	CacheState            = vitessmysql.CacheState
	CachingStorage        = vitessmysql.CachingStorage
	Conn                  = vitessmysql.Conn
	ExtendedHandler       = vitessmysql.ExtendedHandler
	GTIDSet               = vitessmysql.GTIDSet
	Getter                = vitessmysql.Getter
	Handler               = vitessmysql.Handler
	HashStorage           = vitessmysql.HashStorage
	Listener              = vitessmysql.Listener
	ListenerConfig        = vitessmysql.ListenerConfig
	ParsedQuery           = vitessmysql.ParsedQuery
	PlainTextStorage      = vitessmysql.PlainTextStorage
	PrepareData           = vitessmysql.PrepareData
	ResultSpoolFn         = vitessmysql.ResultSpoolFn
	SQLError              = vitessmysql.SQLError
	UserValidator         = vitessmysql.UserValidator
)

var (
	DeserializeCachingSha2PasswordAuthString = vitessmysql.DeserializeCachingSha2PasswordAuthString
	NewListenerWithConfig                    = vitessmysql.NewListenerWithConfig
	NewMariadbBinlogEvent                    = vitessmysql.NewMariadbBinlogEvent
	NewMysqlClearAuthMethod                  = vitessmysql.NewMysqlClearAuthMethod
	NewMysqlNativeAuthMethod                 = vitessmysql.NewMysqlNativeAuthMethod
	NewSalt                                  = vitessmysql.NewSalt
	NewSha2CachingAuthMethod                 = vitessmysql.NewSha2CachingAuthMethod
	NewSQLError                              = vitessmysql.NewSQLError
	ParseMysql56GTIDSet                      = vitessmysql.ParseMysql56GTIDSet
	SerializeCachingSha2PasswordAuthString   = vitessmysql.SerializeCachingSha2PasswordAuthString
)
