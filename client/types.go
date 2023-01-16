package client

import "net/url"

// RequestConfig //
type RequestConfig struct {
	Endpoint   string
	Params     url.Values
	Pagination RequestPagination
}

func NewRequestConfig(endpoint string) RequestConfig {
	return RequestConfig{
		Endpoint:   endpoint,
		Params:     url.Values{},
		Pagination: RequestPagination{},
	}
}

type RequestPagination struct {
	Limit  int
	Offset int
}

// Pagination //
type Pagination struct {
	SortBy  string  `json:"sortBy"`
	Limit   int     `json:"limit"`
	SortDir SortDir `json:"sortDir"`
	Offset  int     `json:"offset"`
}

type PaginatedResponse struct {
	Pg Pagination `json:"pg"`
}

type SortDir string

const (
	SortDirAsc  SortDir = "asc"
	SortDirDesc SortDir = "desc"
)

// ApiToken //
type ApiToken struct {
	Token   string   `json:"token"`
	Label   string   `json:"label"`
	Scope   []string `json:"scope"`
	Created string   `json:"created"`
}

type ApiTokenListResponse struct {
	ApiTokens []ApiToken `json:"apiTokens"`
}

// Connection //
type Connection struct {
	Id           int                `json:"id"`
	Name         string             `json:"name"`
	AccountName  string             `json:"accountName"`
	AccountLabel string             `json:"accountLabel"`
	PackageName  string             `json:"packageName"`
	Expire       string             `json:"expire"`
	Metadata     ConnectionMetadata `json:"metadata,omitempty"`
	TeamId       int                `json:"teamId"`
	Upgradeable  bool               `json:"upgradeable"`
	Scoped       bool               `json:"scoped"`
	Scopes       []ConnectionScope  `json:"scopes,omitempty"`
	AccountType  string             `json:"accountType"`
	Editable     bool               `json:"editable"`
	Uid          string             `json:"uid"`
}

type ConnectionMetadata struct {
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type ConnectionScope struct {
	Id      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Account string `json:"account,omitempty"`
}

type ConnectionResponse struct {
	Connection Connection `json:"connection"`
}

type ConnectionListResponse struct {
	Connections []Connection `json:"connections"`
	Pg          Pagination   `json:"pg"`
}

// DataStore //
type DataStore struct {
	Id              int    `json:"id"`
	Name            string `json:"name"`
	Records         any    `json:"records"`
	Size            string `json:"size"`
	MaxSize         string `json:"maxSize"`
	DatastructureId int    `json:"datastructureId"`
	TeamId          int    `json:"teamId"`
}

type DataStoreResponse struct {
	DataStore DataStore `json:"dataStore"`
}

type DataStoreListResponse struct {
	DataStores []DataStore `json:"dataStores"`
	Pg         Pagination  `json:"pg"`
}

// Organization //
type Organization struct {
	Id          int                 `json:"id"`
	Name        string              `json:"name"`
	CountryId   int                 `json:"countryId"`
	TimezoneId  int                 `json:"timezoneId"`
	License     OrganizationLicence `json:"license"`
	Zone        string              `json:"zone"`
	ServiceName string              `json:"serviceName"`
	IsPaused    bool                `json:"isPaused"`
	ExternalId  string              `json:"externalId"`
	Teams       []OrganizationTeam  `json:"teams"` // used to load make_connection
}

type OrganizationTeam struct {
	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type OrganizationLicence struct {
	Apps       []string `json:"apps"`
	Users      int      `json:"users"`
	Dslimit    int64    `json:"dslimit"`
	Fslimit    int64    `json:"fslimit"`
	Iolimit    int64    `json:"iolimit"`
	Dsslimit   int64    `json:"dsslimit"`
	Fulltext   bool     `json:"fulltext"`
	Interval   int      `json:"interval"`
	Transfer   int64    `json:"transfer"`
	Operations int64    `json:"operations"`
}

type OrganizationResponse struct {
	Organization Organization `json:"organization"`
}

type OrganizationListResponse struct {
	Organizations []Organization `json:"organizations"`
	Pagination    Pagination     `json:"pg"`
}

// OrganizationVariable //
type OrganizationVariable struct {
	Name           string `json:"name"`
	TypeId         int    `json:"typeId"`
	Value          any    `json:"value"`
	IsSystem       bool   `json:"isSystem"`
	OrganizationId int
}

type OrganizationVariableListResponse struct {
	OrganizationVariables []OrganizationVariable `json:"organizationVariables"`
}

// Team //
type Team struct {
	Id             int    `json:"id"`
	Name           string `json:"name"`
	OrganizationId int    `json:"organizationId"`
}

type TeamResponse struct {
	Team Team `json:"team"`
}

type TeamListResponse struct {
	Teams      []Team     `json:"teams"`
	Pagination Pagination `json:"pg"`
}

// TeamVariable //
type TeamVariable struct {
	Name     string `json:"name"`
	TypeId   int    `json:"typeId"`
	Value    any    `json:"value"`
	IsSystem bool   `json:"isSystem"`
	TeamId   int
}

type TeamVariableListResponse struct {
	TeamVariables []TeamVariable `json:"teamVariables"`
}

// User //
type User struct {
	Id             int          `json:"id"`
	Name           string       `json:"name"`
	Email          string       `json:"email"`
	Language       string       `json:"language"`
	TimezoneId     int          `json:"timezoneId"`
	LocaleId       int          `json:"localeId"`
	CountryId      int          `json:"countryId"`
	Features       UserFeatures `json:"features"`
	Avatar         string       `json:"avatar"`
	LastLogin      string       `json:"lastLogin"`
	OrganizationId int
	TeamId         int
}

type UserFeatures struct {
	AllowApps       bool `json:"allow_apps"`
	AllowAppsJs     bool `json:"allow_apps_js"`
	PrivateModules  bool `json:"private_modules"`
	AllowAppsCommit bool `json:"allow_apps_commit"`
	LocalAccess     bool `json:"local_access"`
}

type UserListResponse struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pg"`
}

// UserOrganizationRole //
type UserOrganizationRole struct {
	UserId         int    `json:"userId"`
	UsersRoleId    int    `json:"usersRoleId"`
	OrganizationId int    `json:"organizationId"`
	Invitation     string `json:"invitation"`
}

type UserOrganizationRoleListResponse struct {
	Users      []UserOrganizationRole `json:"userOrganizationRoles"`
	Pagination Pagination             `json:"pg"`
}

// UserRole //
type UserRole struct {
	Id          int      `json:"id"`
	Name        string   `json:"name"`
	Subsidiary  bool     `json:"subsidiary"`
	Category    string   `json:"category"`
	Permissions []string `json:"permissions"`
}

type UserRoleListResponse struct {
	UserRoles  []UserRole `json:"usersRoles"`
	Pagination Pagination `json:"pg"`
}

// UserTeamRole //
type UserTeamRole struct {
	UserId      int  `json:"userId"`
	UsersRoleId int  `json:"usersRoleId"`
	TeamId      int  `json:"teamId"`
	Changeable  bool `json:"changeable"`
}

type UserTeamRoleListResponse struct {
	UserTeamRoles []UserTeamRole `json:"userTeamRoles"`
	Pagination    Pagination     `json:"pg"`
}
