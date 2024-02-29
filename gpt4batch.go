/*
Copyright 2023 The gpt4batch Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package gpt4batch

import (
	"context"
	"io"
)

const (
	// MyFiles is the files' status. [my_files] replace ace_upload
	MyFiles = "my_files"
	// Multimodal is the multimodal status. [multimodal]
	Multimodal = "multimodal"
)

// Logger represents an abstracted structured logging implementation. It
// provides methods to trigger log messages at various alert levels and a
// WithField method to set keys for a structured log message.
type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Error(...interface{})
	Warn(...interface{})

	WithField(string, interface{}) Logger

	// Writer Logger can be transformed into an io.Writer.
	// That writer is the end of an io.Pipe and it is your responsibility to close it.
	Writer() *io.PipeWriter
}

// Source is the source of the data.
type Source struct {
	ID string `json:"id"`
	// URL is the URL of the source.
	URL string `json:"url"`
	// UploadURL is the upload URL of the source.
	UploadURL string `json:"upload_url"`
	// Name is the name of the source.
	Name string `json:"name"`
	// Pid is the pid of the source.
	Pid string `json:"pid"`
	// Prefix is the prefix of the source.
	Prefix string `json:"prefix"`
	// AccessToken is the access token of the source.
	AccessToken string `json:"access_token"`
	// Dir  is the dir of the source.
	Dir string `json:"dir"`
}

// UploadRequest is the request for uploading a file.
type UploadRequest struct {
	*Source `json:",inline"`
	// UploadPath is the path of the file to upload.need to specify the path to upload files or images.
	UploadPath string `json:"upload_path"`
	// ConversationId is the conversation id.if you need to associate previous conversations, please fill in the conversation id.
	ConversationId string `json:"conversation_id,omitempty"`
	// UploadType is the upload type. [file, image]
	UploadType string `json:"upload_type,omitempty"`
}

// ChatResponse is the response for chatting.
type ChatResponse struct {
	Created        int64         `json:"created"`
	MessageID      string        `json:"message_id"`
	ConversationID string        `json:"conversation_id"`
	EndTurn        bool          `json:"end_turn"`
	Contents       []interface{} `json:"contents"`
	Downloads      []string      `json:"downloads,omitempty"`
	// SpecDownloads is the spec downloads for chat service. [origin, local]
	SpecDownloads SpecDownloads `json:"spec_downloads,omitempty"`
}

// SpecDownloads is the spec downloads for chat service.
type SpecDownloads []*SpecDownload

// SpecDownload is the spec downloads for chat service.
type SpecDownload struct {
	// Origin is the origin url of the file.
	Origin string `json:"origin"`
	// Local is the local path of the file.
	Local string `json:"local"`
}

// ChatRequest is the request for chatting.
type ChatRequest struct {
	*Source                    `json:",inline"`
	GizmoId                    string      `json:"gizmo_id"`                    // GizmoId gizmo_id is the gizmo id.
	Message                    string      `json:"message"`                     // Messages: use in web (example: user , bot)
	ParentMessageID            string      `json:"parent_message_id,omitempty"` // ParentMessageID: use in web (example: 1234567890)
	ConversationID             string      `json:"conversation_id,omitempty"`   // ConversationID: use in web (example: 1234567890)
	Stream                     bool        `json:"stream,omitempty"`            // Stream: use in web (example: true , false)
	Model                      string      `json:"model"`                       // Model: use in web (example: gpt3 ,gpt4 )
	Attachments                Attachments `json:"attachments,omitempty"`
	Parts                      Parts       `json:"parts,omitempty"`
	HistoryAndTrainingDisabled bool        `json:"history_and_training_disabled,omitempty"`
}

// Openai is the openai chat request.
func (c ChatRequest) Openai() *OpenaiChatRequest {
	return &OpenaiChatRequest{
		GizmoId:                    c.GizmoId,
		Message:                    c.Message,
		ParentMessageID:            c.ParentMessageID,
		ConversationID:             c.ConversationID,
		Stream:                     c.Stream,
		Model:                      c.Model,
		Attachments:                c.Attachments,
		Parts:                      c.Parts,
		HistoryAndTrainingDisabled: c.HistoryAndTrainingDisabled,
	}
}

// OpenaiChatRequest is the openai request.
type OpenaiChatRequest struct {
	GizmoId                    string      `json:"gizmo_id"`                    // GizmoId gizmo_id is the gizmo id.
	Message                    string      `json:"message"`                     // Messages: use in web (example: user , bot)
	ParentMessageID            string      `json:"parent_message_id,omitempty"` // ParentMessageID: use in web (example: 1234567890)
	ConversationID             string      `json:"conversation_id,omitempty"`   // ConversationID: use in web (example: 1234567890)
	Stream                     bool        `json:"stream,omitempty"`            // Stream: use in web (example: true , false)
	Model                      string      `json:"model"`                       // Model: use in web (example: gpt3 ,gpt4 )
	Attachments                Attachments `json:"attachments,omitempty"`
	Parts                      Parts       `json:"parts,omitempty"`
	HistoryAndTrainingDisabled bool        `json:"history_and_training_disabled,omitempty"`
}

// ChatMessages is the messages for chat service.
type ChatMessages []ChatMessage

// ChatMessage is the message for chat service.
type ChatMessage struct {
	Role        string      `json:"role"`                  // Role: use in web (example: user , bot)
	Content     string      `json:"content"`               // Content: use in web (example: hello , how are you)
	Attachments Attachments `json:"attachments,omitempty"` // Attachments: use in web (example: files , video)
	Parts       Parts       `json:"parts,omitempty"`       // Parts: gpt-4 multimodal (image)
}

// Parts is the parts for chat service.
type Parts []*Part

// Part is the part for chat service.
type Part struct {
	Name         string `json:"name,omitempty"`
	AssetPointer string `json:"asset_pointer"`
	SizeBytes    int    `json:"size_bytes"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	MimeType     string `json:"mimeType,omitempty"`
}

// Attachments is the attachments for chat service.
type Attachments []*Attachment

// Attachment is the attachment for chat service.
type Attachment struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Size          int64  `json:"size"`
	FileTokenSize int    `json:"fileTokenSize,omitempty"`
	MimeType      string `json:"mimeType"`
	Width         int    `json:"width,omitempty"`
	Height        int    `json:"height,omitempty"`
}

// UploadResponse is the response for uploading a file.
type UploadResponse struct {
	ConversationId string      `json:"conversation_id"`
	Attachment     *Attachment `json:"attachment,omitempty"`
	Part           *Part       `json:"part"`
}

// DownloadRequest is the request for downloading a file.
type DownloadRequest struct {
	*Source `json:",inline"`
	// LocalDir is the local dir of the file. [current dir]
	LocalDir string `json:"local_dir"`
	// LocalFileName is the local file name of the file. [current filename]
	LocalFileName string `json:"local_file_name"`
}

// Client is the interface for the client.
type Client interface {
	// Upload uploads a file to the server.
	Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
	// Chat sends a message to the server and returns the response.
	Chat(ctx context.Context, req *ChatRequest) (*ChatResponse, error)
	// Download downloads a file from the server.
	Download(ctx context.Context, req *DownloadRequest) error
	// Close closes the client.
	Close(ctx context.Context) error
}

// Service is the interface for the service.
type Service interface {
	// Open opens the service.
	Open(ctx context.Context) error
	// Done is done.
	Done() <-chan struct{}
	// WithLogger sets the logger for the service.
	WithLogger(log Logger)
	// Close closes the service.
	Close(ctx context.Context) error
}

// Ins is the ins for the service.
type Ins []*In

// In is the in for the service.
type In struct {
	ID      string        `json:"id"`
	Asks    Asks          `json:"asks"`
	Answers []interface{} `json:"answers"`
	IErr    *IErr         `json:"iErr,omitempty"`
	Extra   interface{}   `json:"extra,omitempty"`
}

// Asks is the asks for the service.
type Asks []*Ask

// Ask is the ask for the service.
type Ask struct {
	ID      string   `json:"id"`
	Content string   `json:"content"`
	Images  []string `json:"images,omitempty"`
	Files   []string `json:"files,omitempty"`
}

// IErr is the error for the service.
type IErr struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
