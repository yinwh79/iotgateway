package device

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	simplejson "github.com/bitly/go-simplejson"
	"strconv"
	//"strings"
	//"sync"
)

var MASTER_SLAVE = map[int]string{
	0: "主机",
	1: "副机",
}
var RUN_STATUS = map[int]string{
	0: "不变",
	1: "自动",
	2: "制冷",
	3: "除湿",
	4: "制热",
	5: "送风",
}
var ON_OFF = map[int]string{
	0: "不变",
	1: "停止",
	2: "运行",
}
var ARIFLOW_STATUS = map[int]string{
	0: "不变",
	1: "自动",
	2: "安静",
	3: "低",
	4: "中",
	5: "高",
	6: "中低",
	7: "中高",
}
var MALFUNCTION = map[int]string{
	0: "无故障",
	1: "故障",
}
var VERTICAL_HORIZONTAL = map[int]string{
	0: "不变",
	1: "摆动",
	2: "位置1",
	3: "位置2",
	4: "位置3",
	5: "位置4",
	6: "位置5",
}
var REMOTE_SET = map[int]string{
	0: "不抑制",
	1: "抑制",
}
var FILTER_STATUS = map[int]string{
	0: "无标志",
	1: "过滤网标志",
}
var NORMAL_ENEREGYSAVE = map[int]string{
	0: "无变化",
	1: "正常运行",
	2: "节能运行",
}
var NORMAL_FANGDONG = map[int]string{
	1: "正常运行",
	2: "防冻液运行",
}
var NORMAL_SPACIAL = map[int]string{
	0: "特殊状态",
	1: "正常状态",
}
var CHUSHUANG = map[int]string{
	0: "无除霜状态",
	1: "除霜状态",
}
var YOUHUISHOU = map[int]string{
	0: "无油回收状态",
	1: "油回收状态",
}
var BENGGUZHANG = map[int]string{
	0: "无泵故障状态",
	1: "泵故障状态",
}
var WAIBUGUANRE = map[int]string{
	0: "无变化",
	1: "释放",
	2: "关热",
}

var SET_FANGDONG = map[int]string{
	0: "无变化",
	1: "释放",
	2: "防冻液运行",
}

//out_m
var LOW_NOISE = map[int]string{
	0: "性能优先无效",
	1: "性能优先有效",
}
var LOW_NOISE_LEVEL = map[int]string{
	0: "释放",
	1: "第1级",
	2: "第2级",
	3: "第3级",
}
var EDRLJSYX = map[int]string{
	1: "释放",
	2: "100%",
	3: "90%",
	4: "80%",
	5: "70%",
	6: "60%",
	7: "50%",
	8: "40%",
}

//all_m
var ALL_ON_OFF = map[int]string{
	0: "不变",
	1: "所有室内机均停止",
	2: "有些室内机正在运行",
}
var ALL_MALFUNCTION = map[int]string{
	0: "所有室内机无故障",
	1: "有些室内机处于故障状态",
}
var JINGJITINGZHI = map[int]string{
	0: "不变",
	1: "释放请求",
	2: "紧急停止请求",
}

type FUJITSU struct {
	//继承于ModebusRtu
	ModbusRtu
	/**************按不同设备自定义*************************/
	sub_addr string
	mtype    string
	/**************按不同设备自定义*************************/
}

func init() {
	RegDevice["FUJITSU"] = &FUJITSU{}
}

func (d *FUJITSU) NewDev(id string, ele map[string]string) (DeviceRWer, error) {
	ndev := new(FUJITSU)
	ndev.Device = d.Device.NewDev(id, ele)
	/***********************初始化设备的特有的参数*****************************/
	ndev.sub_addr = ele["sub_addr"]
	ndev.mtype = ele["mtype"]
	ndev.BaudRate = 19200
	ndev.DataBits = 8
	ndev.StopBits = 1
	ndev.Parity = "E"
	//	saint, _ := strconv.Atoi(ele["Starting_address"])
	//ndev.Starting_address = 2
	//	qint, _ := strconv.Atoi(ele["Quantity"])
	//ndev.Quantity = 22
	/***********************初始化设备的特有的参数*****************************/
	return ndev, nil
}

func (d *FUJITSU) GetElement() (dict, error) {
	conn := dict{
		"devaddr": d.devaddr,
		/***********************设备的特有的参数*****************************/
		"mtype":    d.mtype,
		"sub_addr": d.sub_addr,
		"commif":   d.commif,
		/***********************设备的特有的参数*****************************/
	}
	data := dict{
		"_devid": d.devid,
		"_type":  d.devtype,
		"_conn":  conn,
	}
	return data, nil
}

/***********************设备的参数说明帮助***********************************/
func (d *FUJITSU) HelpDoc() interface{} {
	conn := dict{
		"devaddr": "设备地址",
		/***********FUJITSU设备的参数*****************************/
		"commif":   "通信接口,比如(rs485-1)",
		"sub_addr": "子地址",
		/***********FUJITSU设备的参数*****************************/
	}
	r_parameter := dict{
		"_devid": "被读取设备对象的id",
		/***********读取设备的参数*****************************/
		/***********读取设备的参数*****************************/
	}

	w_parameter := dict{
		"_devid": "被操作设备对象的id",
		/***********操作设备的参数*****************************/
		/***********操作设备的参数*****************************/
	}
	data := dict{
		"_devid": "添加设备对象的id",
		"_type":  "FUJITSU", //设备类型
		"_conn":  conn,
	}
	dev_update := dict{
		"request": dict{
			"cmd":  "manager/dev/update.do",
			"data": data,
		},
	}
	readdev := dict{
		"request": dict{
			"cmd":  "do/getvar",
			"data": r_parameter,
		},
	}
	writedev := dict{
		"request": dict{
			"cmd":  "do/setvar",
			"data": w_parameter,
		},
	}
	helpdoc := dict{
		"1.添加设备": dev_update,
		"2.读取设备": readdev,
		"3.操作设备": writedev,
	}
	return helpdoc
}

/***********************设备的参数说明帮助***********************************/

/***************************************添加设备参数检验**********************************************/
func (d *FUJITSU) CheckKey(ele dict) (bool, error) {
	_, sb_ok := ele["sub_addr"].(json.Number)
	if !sb_ok {
		return false, errors.New(fmt.Sprintf("FUJITSU device must have int type element 室内外机编号:sub_addr"))
	}
	_, mt_ok := ele["mtype"].(string)
	if !mt_ok {
		return false, errors.New(fmt.Sprintf("FUJITSU device must have  element 室内还是室外机{in_m|out_m}):mtype"))
	}
	return true, nil
}

/***************************************添加设备参数检验**********************************************/

func (d *FUJITSU) inside_status(ret dict, iarray []int) {
	ret["VRF地址"] = fmt.Sprintf("%d-%d", iarray[1], iarray[0])
	ret["主副机信息"] = MASTER_SLAVE[iarray[5]]
	ret["运行模式状态"] = RUN_STATUS[iarray[7]]
	ret["运行开关状态"] = ON_OFF[iarray[9]]
	ret["设置温度状态"] = (iarray[10]*0x100 + iarray[11]) / 4
	ret["气流状态"] = ARIFLOW_STATUS[iarray[13]]
	ret["室内温度状态"] = (iarray[14]*0x100 + iarray[15]) / 4
	ret["故障监控"] = MALFUNCTION[iarray[17]]
	ret["垂直空气方向位置状态"] = VERTICAL_HORIZONTAL[iarray[19]]
	ret["水平空气方向位置状态"] = VERTICAL_HORIZONTAL[iarray[21]]
	ret["遥控器运行禁止设置状态"] = map[string]string{
		"全部运行设置":    REMOTE_SET[iarray[23]&0x01],
		"定时器设置":     REMOTE_SET[iarray[23]&0x02/0x02],
		"室温设置":      REMOTE_SET[iarray[23]&0x04/0x04],
		"运行模式设置":    REMOTE_SET[iarray[23]&0x08/0x08],
		"启动/停止运行设置": REMOTE_SET[iarray[23]&0x10/0x10],
		"启动行设置":     REMOTE_SET[iarray[23]&0x20/0x20],
		"过滤网重置运行":   REMOTE_SET[iarray[23]&0x40/0x40],
	}
	ret["过滤网标志状态"] = FILTER_STATUS[iarray[25]]
	ret["经济模式运行状态"] = NORMAL_ENEREGYSAVE[iarray[27]]
	ret["防冻液运行状态"] = NORMAL_FANGDONG[iarray[29]]
	ret["温度上下限设置状态(制冷/干燥)"] = fmt.Sprintf("上限=%0.1f,下限=%0.1f", float64(iarray[31])/2, float64(iarray[30])/2)
	ret["温度上下限设置状态(加热)"] = fmt.Sprintf("上限=%0.1f,下限=%0.1f", float64(iarray[33])/2, float64(iarray[32])/2)
	ret["温度上下限设置状态(自动)"] = fmt.Sprintf("上限=%0.1f,下限=%0.1f", float64(iarray[35])/2, float64(iarray[34])/2)
	ret["室内机状态"] = map[string]string{
		"正常状态": NORMAL_SPACIAL[iarray[37]&0x01],
		"除霜":   CHUSHUANG[iarray[37]&0x02/0x02],
		"油回收":  YOUHUISHOU[iarray[37]&0x04/0x04],
		"泵故障":  BENGGUZHANG[iarray[37]&0x08/0x08],
	}
	ret["外部关热状态"] = WAIBUGUANRE[iarray[39]]
}
func (d *FUJITSU) outside_status(ret dict, iarray []int) {
	ret["室外机低噪音运行状态监控"] = map[string]string{
		"性能优先": LOW_NOISE[iarray[5]&0x01],
		"级别":   LOW_NOISE_LEVEL[iarray[5]>>1],
	}
	ret["室外机额定容量节省运行监控"] = EDRLJSYX[iarray[7]]
	ret["主副机信息"] = MASTER_SLAVE[iarray[3]]
	ret["VRF地址"] = fmt.Sprintf("%d-%d", iarray[1], iarray[0])

}

func (d *FUJITSU) all_in_status(ret dict, iarray []int) {
	ret["所有室内机故障监控"] = ALL_MALFUNCTION[iarray[1]]
	ret["所有室内机开/关状态"] = ALL_ON_OFF[iarray[3]]
}

func (d *FUJITSU) encode(m dict) (json.Number, error) {
	name, _ := m["_varname"]
	var results json.Number
	switch name {
	case "运行模式设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					for k, v := range RUN_STATUS {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 0
							log.Debugln("set 运行模式状态 = ", results)
						}
					}
				}
			}
		}
	case "运行开关设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					for k, v := range ON_OFF {
						if v == sval {
							results = json.Number(strconv.Itoa(k))
							d.Starting_address += 1
							log.Debugln("set 运行开关设置 = ", results)
						}
					}
				}
			}
		}
	case "设置温度设定值":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 2
						log.Debugln("set 设置温度设定值 = ", results)
					}
				}
			}
		}
	case "气流设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					for k, v := range ARIFLOW_STATUS {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 3
							log.Debugln("set 气流设置 = ", results)
						}
					}
				}
			}
		}
	case "垂直空气方向位置状态":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "in_m" {
				if sval, ok := val.(string); ok {
					for k, v := range VERTICAL_HORIZONTAL {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 4
							log.Debugln("set 垂直空气方向位置状态 = ", results)
						}
					}
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	case "水平空气方向位置状态":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "in_m" {
				if sval, ok := val.(string); ok {
					for k, v := range VERTICAL_HORIZONTAL {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 5
							log.Debugln("set 水平空气方向位置状态 = ", results)
						}
					}
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	case "遥控器运行禁止设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if sval == "允许" {
						results = json.Number("255")
					} else {
						results = json.Number("0")
					}
					d.Starting_address += 6
					if d.mtype == "all_m" {
						d.Starting_address -= 2
					}
					log.Debugln("set 遥控器运行禁止设置 = ", results)
				}
			}
		}
	case "过滤网标志重置":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "in_m" {
				if sval, ok := val.(string); ok {
					if sval == "重置" {
						results = json.Number("1")
					} else {
						results = json.Number("0")
					}
					d.Starting_address += 7
					log.Debugln("set 过滤网标志重置 = ", results)
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	case "经济运行模式设置":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "in_m" {
				if sval, ok := val.(string); ok {
					for k, v := range NORMAL_ENEREGYSAVE {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 8
							log.Debugln("set 经济运行模式设置 = ", results)
						}
					}
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	case "防冻液运行设置":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "in_m" {
				if sval, ok := val.(string); ok {
					for k, v := range SET_FANGDONG {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 9
							log.Debugln("set 防冻液运行设置 = ", results)
						}
					}
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	case "制冷/干燥温度上限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 10
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 制冷/干燥温度上限设置 = ", results)
					}
				}
			}
		}
	case "制冷/干燥温度下限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 11
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 制冷/干燥温度下限设置 = ", results)
					}
				}
			}
		}
	case "加热温度上限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 12
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 加热温度上限设置 = ", results)
					}
				}
			}
		}
	case "加热温度下限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 13
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 加热温度下限设置 = ", results)
					}
				}
			}
		}
	case "自动温度上限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 14
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 自动温度上限设置 = ", results)
					}
				}
			}
		}
	case "自动温度下限设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					if isvalm, err := strconv.Atoi(sval); err == nil {
						si := isvalm*8 + 1
						results = json.Number(strconv.Itoa(si))
						d.Starting_address += 15
						if d.mtype == "all_m" {
							d.Starting_address -= 5
						}
						log.Debugln("set 自动温度下限设置 = ", results)
					}
				}
			}
		}
	case "外部关热设置":
		{
			if val, ok := m["_varvalue"]; ok {
				if sval, ok := val.(string); ok {
					for k, v := range WAIBUGUANRE {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 16
							if d.mtype == "all_m" {
								d.Starting_address -= 5
							}
							log.Debugln("set 外部关热设置 = ", results)
						}
					}
				}
			}
		}
	case "紧急停止":
		{
			if val, ok := m["_varvalue"]; ok && d.mtype == "all_m" {
				if sval, ok := val.(string); ok {
					for k, v := range JINGJITINGZHI {
						if v == sval {
							sk := strconv.Itoa(k)
							results = json.Number(sk)
							d.Starting_address += 12
							log.Debugln("set 紧急停止说明 = ", results)
						}
					}
				}
			} else {
				return results, errors.New("设置参数错误")
			}
		}
	default:
		{
			return json.Number("0"), errors.New("错误的_varname")
		}
	}

	return results, nil
}

/***************************************读写接口实现**************************************************/
func (d *FUJITSU) RWDevValue(rw string, m dict) (ret dict, err error) {
	ret = make(dict)
	defer func() {
		if drive_err := recover(); drive_err != nil {
			log.Errorf("drive programer  error : (%s)", drive_err)
			errstr := fmt.Sprintf("drive programer  error : %s", drive_err)
			err = errors.New(errstr)
		}
	}()
	d.mutex.Lock()
	defer d.mutex.Unlock()
	//log.SetLevel(log.DebugLevel)
	ret["_devid"] = d.devid
	if rw == "r" {
		if d.mtype == "in_m" {
			d.Quantity = 25
			d.Function_code = 4
			var addr uint16
			if dno, err := strconv.Atoi(d.sub_addr); err == nil {
				addr = uint16(dno)
				if 1 > addr || addr > 128 {
					return nil, errors.New("室内机地址参数错误")
				}
				d.Starting_address = 60*(addr-1) + 50
				log.Debugln("start_address=", d.Starting_address)
				bmdict, berr := d.ModbusRtu.RWDevValue("r", nil)
				if berr != nil {
					bmdict, berr = d.ModbusRtu.RWDevValue("r", nil)
				}
				if berr == nil {
					btdl := bmdict["Modbus-value"]
					bdl, _ := btdl.([]int)
					log.Debugf("室内机-%d receive data = %d", addr, bdl)
					d.inside_status(ret, bdl)

				} else {
					ret["error"] = err.Error()
					log.Debugln(ret)
					return ret, nil
				}
			}
		} else if d.mtype == "out_m" {
			d.Quantity = 4
			d.Function_code = 4
			var addr uint16
			if dno, err := strconv.Atoi(d.sub_addr); err == nil {
				addr = uint16(dno)
				if 1 > addr || addr > 100 {
					return nil, errors.New("室外机地址参数错误")
				}
				d.Starting_address = 15*(addr-1) + 7740
				log.Debugln("start_address=", d.Starting_address)
				bmdict, berr := d.ModbusRtu.RWDevValue("r", nil)
				if berr == nil {
					btdl := bmdict["Modbus-value"]
					bdl, _ := btdl.([]int)
					log.Debugf("室外机-%d receive data = %d", addr, bdl)
					d.outside_status(ret, bdl)

				} else {
					ret["error"] = err.Error()
					log.Debugln(ret)
					return ret, nil
				}
			}
		} else if d.mtype == "all_m" {
			d.Quantity = 2
			d.Function_code = 4
			var addr uint16
			d.Starting_address = 7730
			log.Debugln("start_address=", d.Starting_address)
			bmdict, berr := d.ModbusRtu.RWDevValue("r", nil)
			if berr == nil {
				log.Debugln(bmdict)
				btdl := bmdict["Modbus-value"]
				bdl, _ := btdl.([]int)
				log.Debugf("ALL室内机-%d receive data = %d", addr, bdl)
				d.all_in_status(ret, bdl)
			} else {
				ret["error"] = err.Error()
				log.Debugln(ret)
				return ret, nil
			}
		}

	} else {
		if d.mtype == "in_m" {
			if _, ok := m["_varname"]; ok {
				d.Function_code = 6
				var addr uint16
				if dno, err := strconv.Atoi(d.sub_addr); err == nil {
					addr = uint16(dno)
					if 1 > addr || addr > 128 {
						return nil, errors.New("室内机地址参数错误")
					}
					//if dno, ok := m["_varvalue"]; ok {
					//addr = getnm(dno)
					d.Starting_address = 60*(addr-1) + 1
					wval, werr := d.encode(m)
					if werr != nil {
						ret["error"] = werr.Error()
						return ret, nil
					}
					log.Debugln("wval", wval)
					log.Debugln("start_address=", d.Starting_address)
					bmdict, berr := d.ModbusRtu.RWDevValue("w", dict{"value": wval})
					if berr == nil {
						log.Infof("设置-%s-%d receive data = %v", m["_varname"], addr, bmdict)
					} else {
						ret["error"] = berr.Error()
						log.Debugln(ret)
						return ret, nil
					}
				} else {
					return nil, errors.New("地址参数错误")
				}
			} else {
				return nil, errors.New("缺少_varname")
			}
		} else if d.mtype == "all_m" {
			d.Function_code = 6
			d.Starting_address = 7680
			wval, werr := d.encode(m)
			if werr != nil {
				ret["error"] = werr.Error()
				log.Debugf("设置all室内机-(%s)-%v", m["_varname"], werr)
				return ret, nil
			}
			log.Debugln("wval", wval)
			log.Debugln("start_address=", d.Starting_address)
			bmdict, berr := d.ModbusRtu.RWDevValue("w", dict{"value": wval})
			if berr == nil {
				log.Infof("设置all室内机-%s receive data = %v", m["_varname"], bmdict)
			} else {
				ret["error"] = berr.Error()
				log.Debugln(ret)
				return ret, nil
			}
		} else {
			return nil, errors.New("地址参数错误")
		}

	}
	jsret, _ := json.Marshal(ret)
	inforet, _ := simplejson.NewJson(jsret)
	pinforet, _ := inforet.EncodePretty()
	log.Info(string(pinforet))
	if err != nil {
		log.Debugln(err)
	}
	return ret, err
}

func getnm(inf interface{}) uint16 {
	var ruint uint16 = 0
	if nm, ok := inf.(json.Number); ok {
		nmi64, _ := nm.Int64()
		ruint = uint16(nmi64)
	} else {
		if nm, ok := inf.(string); ok {
			snm, _ := strconv.Atoi(nm)
			ruint = uint16(snm)
		}
	}
	return ruint
}
