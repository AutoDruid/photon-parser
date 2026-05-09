package session_test

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/AutoDruid/photon-parser"
	"github.com/AutoDruid/photon-parser/internal/assembler"
	"github.com/AutoDruid/photon-parser/internal/context"
	v16 "github.com/AutoDruid/photon-parser/internal/parameters/v16"
	v18 "github.com/AutoDruid/photon-parser/internal/parameters/v18"
	"github.com/AutoDruid/photon-parser/internal/reader"
	"github.com/AutoDruid/photon-parser/internal/session"
	"github.com/AutoDruid/photon-parser/internal/types"
)

func TestParseSessionv16(t *testing.T) {

	payload := "00:00:00:07:4e:71:2d:5f:40:41:ad:15:06:00:01:00:00:00:01:5e:00:00:00:0e:f3:04:01:00:29:00:69:00:04:6c:cb:01:73:00:0d:4a:4f:53:45:47:41:4d:45:52:33:31:31:33:02:62:73:03:62:1a:05:78:00:00:00:05:01:00:03:02:05:06:78:00:00:00:05:00:03:0a:00:05:07:78:00:00:00:10:12:29:0d:f6:67:16:00:43:a8:fa:58:55:10:2c:4c:cf:0a:73:00:00:0b:69:00:00:00:00:0c:69:00:00:00:00:0d:73:00:00:10:78:00:00:00:08:c2:46:fb:d4:dd:d8:bd:e8:11:78:00:00:00:08:c2:46:fb:d4:dd:d8:bd:e8:12:69:01:96:b9:09:13:66:40:12:99:68:14:66:40:9e:66:66:16:66:44:c1:c0:00:17:66:44:c1:c0:00:19:66:41:78:02:de:1a:69:01:96:bd:52:1b:66:43:20:00:00:1c:66:43:20:00:00:1e:66:3f:ff:95:84:1f:69:01:96:bc:f3:20:66:44:bb:e0:00:21:66:44:bb:e0:00:23:66:41:b6:3d:71:24:69:01:96:bd:53:25:66:47:43:19:f2:26:6c:08:d8:a9:fd:9a:b5:35:c1:27:62:00:28:79:00:0a:6b:00:00:00:00:00:00:00:00:0e:1f:0b:26:15:97:0b:ab:00:00:00:00:2b:79:00:0e:6b:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:0f:5e:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:ff:33:73:00:00:35:62:00:36:62:00:37:79:00:07:6b:00:0c:ff:ff:00:33:00:af:00:f7:ff:ff:01:43:38:62:04:39:62:0c:3f:69:00:00:00:00:fc:6b:00:1d:06:00:01:00:00:00:01:a1:00:00:00:0f:f3:04:01:00:0c:00:69:00:04:6c:cb:01:79:00:0b:6b:09:10:04:98:02:21:0f:d9:04:97:03:d2:02:e4:03:9a:03:b7:03:b8:03:cc:02:79:00:0b:66:42:c8:00:00:43:8a:7e:5d:42:c8:00:00:43:4b:13:26:43:8a:7e:5d:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:03:79:00:0b:66:42:c8:00:00:43:3c:fa:21:42:c8:00:00:43:44:7d:8a:43:3c:fa:21:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:42:c8:00:00:04:79:00:0b:6c:08:de:8e:a4:47:26:4d:8a:08:de:8e:a4:ac:1e:41:0d:08:de:86:b2:41:9d:ca:71:08:de:8f:4b:16:7f:6b:81:08:de:8f:4b:17:35:2d:e1:08:de:8f:4b:2e:e5:f4:8f:08:de:8f:4b:2e:e6:53:c7:08:de:8f:4b:2e:e6:53:c7:08:de:8f:4b:2e:e6:53:c7:08:de:8f:4b:2e:e6:53:c7:08:de:8f:4b:2e:e5:e7:ff:05:78:00:00:00:0b:01:01:01:01:01:01:01:01:01:01:01:06:78:00:00:00:00:07:78:00:00:00:11:5b:82:6d:49:82:24:49:82:01:db:b6:6d:db:b6:6d:db:02:08:79:00:19:69:80:00:00:00:80:00:00:00:00:1b:77:40:00:1b:77:40:00:1b:77:40:0d:d6:c7:08:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:75:30:00:00:3a:98:00:00:ea:60:00:00:ea:60:00:00:ea:60:00:00:ea:60:00:00:ea:60:00:00:ea:60:00:00:ea:60:09:78:00:00:00:02:64:04:0a:79:00:04:6b:00:00:00:00:00:00:02:8c:fc:6b:00:0b:06:00:01:00:00:00:00:45:00:00:00:10:f3:04:01:00:0a:00:69:00:04:6c:cb:01:69:01:96:bd:75:02:6b:0b:ab:03:6b:2b:be:04:66:44:bb:e0:00:05:66:41:b6:3d:71:06:6c:08:de:8f:4b:2f:8d:40:bf:08:6f:01:0b:6f:01:fc:6b:00:d1:06:00:01:00:00:00:00:18:00:00:00:11:f3:04:01:00:02:00:62:64:fc:6b:01:43:06:00:01:00:00:00:00:5e:00:00:00:12:f3:84:f5:26:3e:06:6f:00:dd:1c:b9:23:98:33:a4:47:4c:5c:e4:ae:99:1c:ce:9f:21:02:6d:98:1c:a4:7e:1f:65:0a:8c:3b:e5:a5:b2:e3:0b:4b:67:aa:e2:3c:13:8f:5c:e0:d4:b6:94:d0:d2:9c:1f:d4:3c:0b:e4:13:1b:78:ad:de:b5:2e:9f:00:5d:93:54:ba:92:f8:35:98:2a:97:cc:60:06:00:01:00:00:00:00:5d:00:00:00:13:f3:04:01:00:0b:00:69:00:04:6c:f1:01:6b:18:74:02:62:01:04:6c:00:00:00:00:ae:4b:61:10:05:73:00:0f:47:72:7a:65:67:6f:72:7a:4e:69:65:4d:6f:67:65:06:62:02:07:69:1c:a1:fc:7d:08:79:00:03:6b:0b:13:0b:1d:0b:20:09:78:00:00:00:01:12:0a:62:00:fc:6b:00:1e:06:00:01:00:00:00:00:53:00:00:00:14:f3:04:01:00:0b:00:69:00:04:6c:f3:01:6b:0f:fc:02:62:01:04:69:48:d9:3c:01:05:73:00:0d:53:61:62:65:72:6c:65:61:74:68:65:72:31:06:62:02:07:69:0f:a8:fd:5e:08:79:00:01:6b:0f:10:09:78:00:00:00:01:ac:0a:62:00:fc:6b:00:1e"
	cleared, _ := hex.DecodeString(strings.ReplaceAll(payload, ":", ""))

	ctx := context.NewContext(
		reader.NewReader(cleared),
		nil,
		nil,
		context.Decoders[v16.Parameter]{
			ParameterParser:              &v16.Parameter{},
			ReliableHeaderParameterCount: &v16.ReliableHeaderParameterCountV16{},
		},
	)

	var sess types.Session
	err := session.Parse(ctx, &sess)

	if err != nil {
		t.Fatalf("error parsing session: %v", err)
	}

	if sess.PeerID != 0 {
		t.Fatalf("error peerid should be 0: %v", sess.PeerID)
	}

	if sess.CRCEnabled != 0 {
		t.Fatalf("error CRCEnabled should be 0: %v", sess.CRCEnabled)
	}

	if sess.CommandCount != 7 {
		t.Fatalf("error CommandCount should be 7: %v", sess.CommandCount)
	}

	if sess.Timestamp != 1316040031 {
		t.Fatalf("error Timestamp should be 1316040031: %v", sess.Timestamp)
	}

	if sess.Challenge != 1078045973 {
		t.Fatalf("error Timestamp should be 1078045973: %v", sess.Challenge)
	}

	if len(sess.Commands) != 7 {
		t.Fatalf("error len commands should be 7: %v", len(sess.Commands))
	}
}

func TestParseSessionv18(t *testing.T) {

	payload := []byte{0x0, 0x0, 0x0, 0x2, 0x16, 0x28, 0x80, 0x6b, 0x48, 0x74, 0x8c, 0x42, 0x1, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2e, 0x3, 0xff, 0x1, 0x0, 0x0, 0x0, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0xd7, 0xa1, 0x4, 0xb0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x13, 0x88, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x2}

	ctx := context.NewContext(
		reader.NewReader(payload),
		assembler.NewAssembler(),
		nil,
		context.Decoders[v18.Parameter]{
			ParameterParser:              &v18.Parameter{},
			ReliableHeaderParameterCount: &v18.ReliableHeaderParameterCountV18{},
		},
	)

	var sess types.Session
	err := session.Parse(ctx, &sess)

	if err != nil {
		t.Fatalf("error parsing session: %v", err)
	}

	if sess.PeerID != 0 {
		t.Fatalf("error peerid should be 0: %v", sess.PeerID)
	}

	if sess.CRCEnabled != 0 {
		t.Fatalf("error CRCEnabled should be 0: %v", sess.CRCEnabled)
	}

	if sess.CommandCount != 2 {
		t.Fatalf("error CommandCount should be 2: %v", sess.CommandCount)
	}

	if sess.Timestamp != 371753067 {
		t.Fatalf("error Timestamp should be 371753067: %v", sess.Timestamp)
	}

	if sess.Challenge != 1215597634 {
		t.Fatalf("error Timestamp should be 1215597634: %v", sess.Challenge)
	}

	if len(sess.Commands) != 2 {
		t.Fatalf("error len commands should be 2: %v", len(sess.Commands))
	}
}

func TestMalFormedHeader(t *testing.T) {

	payload := []byte{0x0, 0x0}

	ctx := context.NewContext(
		reader.NewReader(payload),
		nil,
		nil,
		context.Decoders[v18.Parameter]{},
	)

	var sess types.Session
	err := session.Parse(ctx, &sess)

	if err == nil {
		t.Fatalf("error parsing session should be error")
	}
}

func TestSyncHookOnSession(t *testing.T) {

	payload := []byte{0x0, 0x0, 0x0, 0x2, 0x16, 0x28, 0x80, 0x6b, 0x48, 0x74, 0x8c, 0x42, 0x1, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2e, 0x3, 0xff, 0x1, 0x0, 0x0, 0x0, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0xd7, 0xa1, 0x4, 0xb0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x13, 0x88, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x2}

	parser := photon.NewV18()

	var sess types.Session
	parser.OnSessionSync(func(s types.Session) {
		sess = s
	})

	_, err := parser.ParsePacket(payload)

	if err != nil {
		t.Fatalf("error parsing session: %v", err)
	}

	if sess.PeerID != 0 {
		t.Fatalf("error peerid should be 0: %v", sess.PeerID)
	}

	if sess.CRCEnabled != 0 {
		t.Fatalf("error CRCEnabled should be 0: %v", sess.CRCEnabled)
	}

	if sess.CommandCount != 2 {
		t.Fatalf("error CommandCount should be 2: %v", sess.CommandCount)
	}

	if sess.Timestamp != 371753067 {
		t.Fatalf("error Timestamp should be 371753067: %v", sess.Timestamp)
	}

	if sess.Challenge != 1215597634 {
		t.Fatalf("error Timestamp should be 1215597634: %v", sess.Challenge)
	}

	if len(sess.Commands) != 2 {
		t.Fatalf("error len commands should be 2: %v", len(sess.Commands))
	}
}

func TestASyncHookOnSession(t *testing.T) {

	payload := []byte{0x0, 0x0, 0x0, 0x2, 0x16, 0x28, 0x80, 0x6b, 0x48, 0x74, 0x8c, 0x42, 0x1, 0xff, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x2e, 0x3, 0xff, 0x1, 0x0, 0x0, 0x0, 0x0, 0x2c, 0x0, 0x0, 0x0, 0x0, 0xd7, 0xa1, 0x4, 0xb0, 0x0, 0x0, 0x80, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x13, 0x88, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x2}

	parser := photon.NewV18()

	ch := parser.OnSessionAsync(types.HookOptions{Size: 1})

	_, err := parser.ParsePacket(payload)

	if err != nil {
		t.Fatalf("error parsing session: %v", err)
	}

	select {
	case got := <-ch:
		if got.PeerID != 0 {
			t.Fatalf("error peerid should be 0: %v", got.PeerID)
		}

		if got.CRCEnabled != 0 {
			t.Fatalf("error CRCEnabled should be 0: %v", got.CRCEnabled)
		}

		if got.CommandCount != 2 {
			t.Fatalf("error CommandCount should be 2: %v", got.CommandCount)
		}

		if got.Timestamp != 371753067 {
			t.Fatalf("error Timestamp should be 371753067: %v", got.Timestamp)
		}

		if got.Challenge != 1215597634 {
			t.Fatalf("error Timestamp should be 1215597634: %v", got.Challenge)
		}

		if len(got.Commands) != 2 {
			t.Fatalf("error len commands should be 2: %v", len(got.Commands))
		}
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for async session")
	}
}
