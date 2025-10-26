package client

type Stack struct {
	ID              int              `json:"Id"`
	Name            string           `json:"Name"`
	Type            int              `json:"Type"`
	EndpointID      int              `json:"EndpointId"`
	SwarmID         string           `json:"SwarmId,omitempty"`
	EntryPoint      string           `json:"EntryPoint"`
	Env             []StackEnv       `json:"Env"`
	Status          int              `json:"Status"`
	CreationDate    int64            `json:"CreationDate"`
	CreatedBy       string           `json:"CreatedBy"`
	UpdateDate      int64            `json:"UpdateDate"`
	UpdatedBy       string           `json:"UpdatedBy"`
	ProjectPath     string           `json:"ProjectPath"`
	AutoUpdate      *StackAutoUpdate `json:"AutoUpdate"`
	GitConfig       *StackGitConfig  `json:"GitConfig"`
	FromAppTemplate bool             `json:"FromAppTemplate"`
	Namespace       string           `json:"Namespace,omitempty"`
	IsComposeFormat bool             `json:"IsComposeFormat"`
	ResourceControl *ResourceControl `json:"ResourceControl"`
}

type StackEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type StackAutoUpdate struct {
	Interval string `json:"Interval"`
	Webhook  string `json:"Webhook"`
}

type StackGitConfig struct {
	URL            string             `json:"URL"`
	ReferenceName  string             `json:"ReferenceName"`
	ConfigFilePath string             `json:"ConfigFilePath"`
	Authentication *GitAuthentication `json:"Authentication"`
	TLSSkipVerify  bool               `json:"TLSSkipVerify"`
}

type GitAuthentication struct {
	Username string `json:"Username"`
	Password string `json:"Password"`
}

type ResourceControl struct {
	ID                 int          `json:"Id"`
	ResourceID         string       `json:"ResourceId"`
	SubResourceIDs     []string     `json:"SubResourceIds"`
	Type               int          `json:"Type"`
	UserAccesses       []UserAccess `json:"UserAccesses"`
	TeamAccesses       []TeamAccess `json:"TeamAccesses"`
	Public             bool         `json:"Public"`
	AdministratorsOnly bool         `json:"AdministratorsOnly"`
	System             bool         `json:"System"`
}

type UserAccess struct {
	UserID      int `json:"UserId"`
	AccessLevel int `json:"AccessLevel"`
}

type TeamAccess struct {
	TeamID      int `json:"TeamId"`
	AccessLevel int `json:"AccessLevel"`
}

//type CreateStackRequest struct {
//	Name             string     `json:"Name"`
//	StackFileContent string     `json:"StackFileContent,omitempty"`
//	Env              []StackEnv `json:"Env,omitempty"`
//	RepositoryURL               string             `json:"RepositoryURL,omitempty"`
//	RepositoryReferenceName     string             `json:"RepositoryReferenceName,omitempty"`
//	ComposeFilePathInRepository string             `json:"ComposeFilePathInRepository,omitempty"`
//	RepositoryAuthentication    *GitAuthentication `json:"RepositoryAuthentication,omitempty"`
//	TLSSkipVerify   bool `json:"TLSSkipVerify,omitempty"`
//	FromAppTemplate bool `json:"FromAppTemplate,omitempty"`
//}

type CreateStackRequest struct {
	Name                        string             `json:"Name"`
	SwarmID                     string             `json:"SwarmID,omitempty"`
	StackFileContent            string             `json:"StackFileContent,omitempty"`
	Env                         []StackEnv         `json:"Env,omitempty"`
	RepositoryURL               string             `json:"RepositoryURL,omitempty"`
	RepositoryReferenceName     string             `json:"RepositoryReferenceName,omitempty"`
	ComposeFilePathInRepository string             `json:"ComposeFilePathInRepository,omitempty"`
	RepositoryAuthentication    *GitAuthentication `json:"RepositoryAuthentication,omitempty"`
	TLSSkipVerify               bool               `json:"TLSSkipVerify,omitempty"`
	FromAppTemplate             bool               `json:"FromAppTemplate,omitempty"`
}

type UpdateStackRequest struct {
	StackFileContent string     `json:"StackFileContent,omitempty"`
	Env              []StackEnv `json:"Env,omitempty"`
	Prune            bool       `json:"Prune,omitempty"`
	PullImage        bool       `json:"PullImage,omitempty"`
}
