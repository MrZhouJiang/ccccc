/**
 * Created with JetBrains PhpStorm.
 * User: kk
 * Date: 13-8-28
 * Time: 下午4:44
 */
function U() {
    var url = arguments[0] || [];
    var param = arguments[1] || {};
    var url_arr = url.split('/');

    if (!$.isArray(url_arr) || url_arr.length < 2 || url_arr.length > 3) {
        return '';
    }

    if (url_arr.length == 2)
        url_arr.unshift(_GROUP_);

    var pre_arr = ['g', 'm', 'a'];

    var arr = [];
    for (d in pre_arr)
        arr.push(pre_arr[d] + '=' + url_arr[d]);

    for (d in param)
        arr.push(d + '=' + param[d]);

    return _APP_ + '?' + arr.join('&');
}


function LoadBaseData() {
    GetFenWei()
    GetUnit()
    GetGoodsType()
}

function GetBaseHost(){
    return "http://192.168.202.5:8889"
   // return "http://127.0.0.1:8889"
}

function GetQueryString(name) {
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
    var r = decodeURI(window.location.search.substr(1)).match(reg);
    if (r != null) return unescape(r[2]);
    return null;
}

function GetFenWei() {
    str = ""
    //开始请求数据
    $.ajax({
        url: GetBaseHost() +"/v2/get_all_fen_wei_list",
        dataType: "json",
        data: {
            "shafa_id": "CESAE",
            "types": "裁工"
        },
        type: "get",
        async: false,
        success: function (data) {
            $.each(data.data, function (index, element) {
                //绘制导出表
                str += " <option  value='" + element.fwmc + "' >" + element.fwmc + "</option>\n"

            })

        },
    })
    return str
}

function GetUnit() {
    str = ""
    //开始请求数据
    $.ajax({
        url: GetBaseHost() +"/v2/get_uint",
        dataType: "json",
        data: {},
        type: "get",
        async: false,
        success: function (data) {
            $.each(data.data, function (index, element) {
                //绘制导出表
                str += " <option  value='" + element.id + "' >" + element.name + "</option>\n"

            })

        },
    })
    return str
}

function GetUnitNoId() {
    str = ""
    //开始请求数据
    $.ajax({
        url: GetBaseHost() +"/v2/get_uint",
        dataType: "json",
        data: {},
        type: "get",
        async: false,
        success: function (data) {
            $.each(data.data, function (index, element) {
                //绘制导出表
                str += " <option  value='" + element.name + "' >" + element.name + "</option>\n"

            })

        },
    })
    return str
}


function GetGoodsType() {
    str = ""
    //开始请求数据
    $.ajax({
        url: GetBaseHost() +"/v2/get_goods_type",
        dataType: "json",
        data: {},
        type: "get",
        async: false,
        success: function (data) {
            $.each(data.data, function (index, element) {
                //绘制导出表
                str += " <option  value='" + element.goods_type_name + "' >" + element.goods_type_name + "</option>\n"

            })

        },
    })

    return str
}


