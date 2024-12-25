package domain

import (
	"testing"
	"time"

	"github.com/ixre/go2o/core/domain/interface/sys"
	"github.com/ixre/go2o/core/inject"
)

func TestGetOptions(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	arr := ia.Options().GetChildOptions(0, "BIZ")
	t.Logf("options = %#v \n", arr)
}

func TestGetAllCities(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	cities := ia.Location().GetAllCities()
	if len(cities) == 0 {
		t.Error("No cities found")
	}
	for _, city := range cities {
		t.Log(city)
	}
}

// 测试根据城市获取站点
func TestGetStationByCityCode(t *testing.T) {
	city := 110100
	ia := inject.GetSystemRepo().GetSystemAggregateRoot() // TODO: write test code here
	d := ia.Location().GetDistrict(city)
	t.Logf("district = %s \n", d.Name)
	station := ia.Stations().FindStationByCity(city)
	t.Logf("station = %#v \n", station.GetValue())
}

// 测试获取应用分发
func TestGetAppDistribution(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot()
	dist := ia.Application().GetAppDistributionByName("Go2o")
	if dist == nil {
		dist = &sys.SysAppDistribution{
			Id:             1,
			AppName:        "go2o",
			AppIcon:        "",
			AppDesc:        "测试应用",
			UpdateMode:     1,
			DistributeUrl:  "https://github.com/ixre/go2o",
			UrlScheme:      "go2o://net.fze.go2o/open",
			DistributeName: "测试应用",
			StableVersion:  "1.0.0",
			StableDownUrl:  "https://github.com/ixre/go2o/releases/download/v1.0.0/go2o-v1.0.0.apk",
			BetaVersion:    "1.0.0",
			BetaDownUrl:    "https://github.com/ixre/go2o/releases/download/v1.0.0/go2o-v1.0.0.apk",
		}
		ia.Application().SaveAppDistribution(dist)
	}
	err := ia.Application().SaveAppVersion(&sys.SysAppVersion{
		Id:              0,
		DistributionId:  dist.Id,
		Version:         "1.0.9",
		VersionCode:     0,
		TerminalOs:      "android",
		TerminalChannel: "stable",
		StartTime:       int(time.Now().Unix()),
		UpdateMode:      1,
		UpdateContent:   "测试更新",
		PackageUrl:      "https://github.com/ixre/go2o/releases/download/v1.0.1/go2o-v1.0.1.apk",
		IsForce:         1,
		IsNotified:      0,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	t.Logf("dist = %#v \n", dist)

	ver := ia.Application().GetLatestVersion(dist.Id, "android", "stable")
	if ver == nil {
		t.Error("No such version")
		t.FailNow()
	}
	t.Logf("latest version = %#v \n", ver)

	dist = ia.Application().GetAppDistributionByName("go2o")
	t.Logf("distribution = %#v \n", dist)
}

// 测试提交应用日志
func TestSubmitAppLog(t *testing.T) {
	ia := inject.GetSystemRepo().GetSystemAggregateRoot()
	err := ia.Application().AddLog(&sys.SysLog{
		UserId:          0,
		Username:        "jarrysix",
		LogLevel:        3,
		Message:         "测试错误信息",
		Arguments:       "{id:1}",
		TerminalModel:   "vscode",
		TerminalName:    "unittest",
		TerminalEntry:   "test",
		TerminalVersion: "1.0.0",
		ExtraInfo:       "",
		CreateTime:      0,
	})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
}
