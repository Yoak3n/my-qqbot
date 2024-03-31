package bilibili

import "testing"

func TestRefreshCookie(t *testing.T) {
	CK = "DedeUserID=33922058;DedeUserID__ckMd5=6c4288b6371877bc;Expires=;SESSDATA=027d4d1d%2C1726675069%2C090cb%2A32CjAh9LOxepDjF6W0Q050KnZddju6j5-m1I2YMohr9ui8DAp3QkuCDD1jiDZXUV8gDZoSVkdIekhXdGxYZmdKakRNWFhBY0J2TTRqek5UODlaS2dDT3JNbG5FZU13ZGtWLTZhOVFMMWJHaDBudGlPZXNsb1hJS3c1ajAwSEs3a0xrckVTRE4wbGtRIIEC;bili_jct=66e231b3f5894ad17d9fc388937119f3;"
	refreshToken = "7d2c4eee8535c08d1dc7a1fbce5c4232"
	err := RefreshCookie()
	if err != nil {
		t.Error(err)
	}
}

func TestGetCsrf(t *testing.T) {
	CK = "DedeUserID=33922058;DedeUserID__ckMd5=6c4288b6371877bc;Expires=;SESSDATA=027d4d1d%2C1726675069%2C090cb%2A32CjAh9LOxepDjF6W0Q050KnZddju6j5-m1I2YMohr9ui8DAp3QkuCDD1jiDZXUV8gDZoSVkdIekhXdGxYZmdKakRNWFhBY0J2TTRqek5UODlaS2dDT3JNbG5FZU13ZGtWLTZhOVFMMWJHaDBudGlPZXNsb1hJS3c1ajAwSEs3a0xrckVTRE4wbGtRIIEC;bili_jct=66e231b3f5894ad17d9fc388937119f3;"
	refreshToken = "7d2c4eee8535c08d1dc7a1fbce5c4232"
	csrf := getCsrf()
	if csrf == "" {
		t.Error("获取csrf失败")
	}
	t.Log(csrf)
}
