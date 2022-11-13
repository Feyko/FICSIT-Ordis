//go:build generate
// +build generate

package smr

import _ "git.sr.ht/~emersion/gqlclient/cmd/gqlintrospect"
import _ "git.sr.ht/~emersion/gqlclient/cmd/gqlclientgen"

//go:generate sh -c "go run git.sr.ht/~emersion/gqlclient/cmd/gqlintrospect https://api.ficsit.app/v2/query > schema.graphqls"
//go:generate go run git.sr.ht/~emersion/gqlclient/cmd/gqlclientgen -s schema.graphqls -o generated.go -q queries.graphql
