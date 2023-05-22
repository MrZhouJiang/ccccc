(function($,window){
    var Page = function(ele,opt){
        this.$ele = ele;
        this.defaults ={
            curPage: 1,
            totalPage: 1,
            totalCount: 0,
            morePage: opt.morePage,
            perPageCount: opt.perPageCount
        }
        this.options = $.extend({},this.defaults,opt);
    }
    Page.prototype = {
        init: function () {
            //数据初始化
            this.dataInit();
            //显示当前页数、总页数
            this.pageInit();
            //分页处理
            this.pageFun();
            //销毁之前事件
            this.offEventFun();
            //事件处理
            this.eventFun();
            return this.$ele;
        },
        pageInit: function () {
            $(this.options.curPageEl).html("当前第" + this.options.curPage + "页");
            $(this.options.totalEl).html("共" + this.options.totalPage + "页");
        },
        pageFun: function () {
            var $list = this.$ele.children();
            $list.hide();
            var start = (this.options.curPage - 1) * this.options.perPageCount;
            if (this.options.curPage == this.options.totalPage) {
                var end = $list.length;
                for (var i = start; i < end; i++) {
                    $($list[i]).show();
                }
            } else {
                for (var i = start; i < start + this.options.perPageCount; i++) {
                    $($list[i]).show();
                }
            }
            this.pageInit();
        }
        ,
        dataInit: function () {
            var $list = this.$ele.children();
            this.options.curPage = 1;
            this.options.totalCount = $list.length;
            this.options.totalPage = Math.ceil($list.length / this.options.perPageCount);
        },
        eventFun:function(){
            //下一页
            var self = this;
            $(this.options.next).on("click", function () {
                if (self.options.curPage + 1 > self.options.totalPage) {
                    alert("已经是最后一页");
                    return;
                }
                self.options.curPage++;
                self.pageFun();
            });
            //上一页
            $(this.options.prev).on("click", function () {
                if (self.options.curPage - 1 < 1) {
                    alert("已经是第一页");
                    return;
                }
                self.options.curPage--;
                self.pageFun();
            });
            //下n页
            $(this.options.nextMore).on("click", function () {
                if (self.options.curPage + self.options.morePage > self.options.totalPage){
                    self.options.curPage = self.options.totalPage;
                    alert("已经是最后一页")
                }else{
                    self.options.curPage += self.options.morePage;
                }
                self.pageFun();
            });
            //上n页
            $(this.options.prevMore).on("click", function () {
                if (self.options.curPage - self.options.morePage < 1){
                    self.options.curPage = 1;
                    alert("已经是第一页")
                }else{
                    self.options.curPage -= self.options.morePage;
                }
                self.pageFun();
            });
        },
        offEventFun:function(){
            $(this.options.next).off("click");
            $(this.options.prev).off("click");
            $(this.options.nextMore).off("click");
            $(this.options.prevMore).off("click");
        }
    }
    $.fn.page = function(options){
        var page = new Page(this,options);
        return page.init();
    }
})(jQuery,window)


function firstPage(){
    hide();
    currPageNum = 1;
    showCurrPage(currPageNum);
    showTotalPage();
    for(i = 1; i < pageCount + 1; i++){
        blockTable.rows[i].style.display = "";
    }

    firstText();
    preText();
    nextLink();
    lastLink();
}

function prePage(){
    hide();
    currPageNum--;
    showCurrPage(currPageNum);
    showTotalPage();
    var firstR = firstRow(currPageNum);
    var lastR = lastRow(firstR);
    for(i = firstR; i < lastR; i++){
        blockTable.rows[i].style.display = "";
    }

    if(1 == currPageNum){
        firstText();
        preText();
        nextLink();
        lastLink();
    }else if(pageNum == currPageNum){
        preLink();
        firstLink();
        nextText();
        lastText();
    }else{
        firstLink();
        preLink();
        nextLink();
        lastLink();
    }

}

function nextPage(){
    hide();
    currPageNum++;
    showCurrPage(currPageNum);
    showTotalPage();
    var firstR = firstRow(currPageNum);
    var lastR = lastRow(firstR);
    for(i = firstR; i < lastR; i ++){
        blockTable.rows[i].style.display = "";
    }

    if(1 == currPageNum){
        firstText();
        preText();
        nextLink();
        lastLink();
    }else if(pageNum == currPageNum){
        preLink();
        firstLink();
        nextText();
        lastText();
    }else{
        firstLink();
        preLink();
        nextLink();
        lastLink();
    }
}

function lastPage(){
    hide();
    currPageNum = pageNum;
    showCurrPage(currPageNum);
    showTotalPage();
    var firstR = firstRow(currPageNum);
    for(i = firstR; i < numCount + 1; i++){
        blockTable.rows[i].style.display = "";
    }

    firstLink();
    preLink();
    nextText();
    lastText();
}

// 计算将要显示的页面的首行和尾行
function firstRow(currPageNum){
    return pageCount*(currPageNum - 1) + 1;
}

function lastRow(firstRow){
    var lastRow = firstRow + pageCount;
    if(lastRow > numCount + 1){
        lastRow = numCount + 1;
    }
    return lastRow;
}

function showCurrPage(cpn){
    currPageSpan.innerHTML = cpn;
}

function showTotalPage(){
    pageNumSpan.innerHTML = pageNum;
}

//隐藏所有行
function hide(){
    for(var i = 1; i < numCount + 1; i ++){
        blockTable.rows[i].style.display = "none";
    }
}

//控制首页等功能的显示与不显示
function firstLink(){firstSpan.innerHTML = "<a href='javascript:firstPage();'>First</a>";}
function firstText(){firstSpan.innerHTML = "First";}

function preLink(){preSpan.innerHTML = "<a href='javascript:prePage();'>Pre</a>";}
function preText(){preSpan.innerHTML = "Pre";}

function nextLink(){nextSpan.innerHTML = "<a href='javascript:nextPage();'>Next</a>";}
function nextText(){nextSpan.innerHTML = "Next";}

function lastLink(){lastSpan.innerHTML = "<a href='javascript:lastPage();'>Last</a>";}
function lastText(){lastSpan.innerHTML = "Last";}