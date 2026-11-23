package main

// Playlist mutations use a cloud protocol or a raw image stream, so their
// request models and transport metadata cannot be inferred from createRequest.
func applyPlaylistModel(m *apiModel) {
	switch m.Identifier {
	case "playlist_pic":
		m.Fields = map[string]fieldType{"data": tString, "total_ver": tInt}
	case "playlist_pic_upload":
		m.Name = "PlaylistUpload"
		m.Spec.BaseURL = "http://imgphp.kugou.com"
		m.Spec.UpstreamURL = "/imageupload/stream.php"
		m.Fields = map[string]fieldType{
			"data": tBytes, "file": tString, "type": tString,
			"extendName": tString, "md5": tString,
		}
	case "playlist_update":
		m.Fields = map[string]fieldType{
			"total_ver": tInt, "listid": tInt, "type": tInt, "sort": tInt,
			"name": tString, "tags": tString, "intro": tString, "pic": tString,
		}
	default:
		return
	}
	m.Method = "POST"
	m.Spec.Method = "POST"
	m.Spec.ClearDefaultParams = true
	m.Spec.NotSignature = true
	m.Spec.UseParams = false
	m.Spec.UseData = true
	if m.Identifier != "playlist_pic_upload" {
		m.Spec.Headers["x-router"] = "cloudlist.service.kugou.com"
	}
}
