/**
 * @author ivy <guoxuivy@gmail.com>
 * @copyright Copyright &copy; 2013-2017 Ivy Software LLC
 * @license https://github.com/guoxuivy/ivy/
 * @package framework
 * @link https://github.com/guoxuivy/ivy 
 * @since 1.0 
 *
 *对话框插件 Idialog 基于jquery 封装
 *
 *
 *
 *
 *
 需要样式
.idialog{ overflow:hidden; position:fixed; left:50%; top:200px;  display:none; z-index:999;}
.idialog_title{ height:30px; background:#07aaff; padding-left:20px; line-height:30px;font-family:"微软雅黑"; font-size:14px; color:#ffffff;}
.idialog_title span{ float:right; display:inline; margin-right:5px; height:30px; width:30px; text-align:center; cursor:pointer;}
.idialog_body{  overflow:hidden;font-family:"微软雅黑"; font-size:12px; color:#666666; background:#FFF; border-left:1px solid #dddddd; border-right:1px solid #dddddd; border-bottom:1px solid #dddddd;}
.idialog_content{padding:20px 30px 5px 5px;}
.idialog_active{ height:60px; background:#fafafa; width:100%;border-top:1px solid #dddddd; text-align:right;}
.idialog_active a{ display:inline-block; height:28px; width:80px; border:1px solid #cccccc; border-radius:5px; margin:0 5px;font-family:"微软雅黑"; font-size:12px; color:#666666; text-align:center; line-height:28px; margin-top:16px;}


使用示例
var d = Idialog({
	top:100,
	width:500,
	content:$('#test'),
	init:function(body){
		//console.log(body);
	},
	ok:function(obj){
		//console.log(obj);
		return false;
	},
	cancel:true
});

 */
(function($,win,dom,undef){

	var template  = '<div class="idialog" style="display: block;">';
		template += '	<div class="idialog_title"><font></font><span></span></div>';
		template += '	<div class="idialog_body">';
		template += '		<div class="idialog_content"></div>';
		template += '		<div class="idialog_active"><a class="idialog_cancel" href="javascript:">取消</a><a href="javascript:" class="idialog_ok">确定</a></div>';
		template += '	</div>';
		template += '</div>';


	var Idialog=function(settings){
		this.settings=$.extend({},Idialog.defaults,settings);
		this._self=$(template);
		//执行
		this.run=function(){
			this.bind();
			this.show();
			//渲染后内容初始化
			this.init();
			return this;
		}
	};

	/**
	 * 弹出框 默认配置 可扩展
	 * @type {Object}
	 */
	Idialog.defaults={
		top:200,
		width:474,
		title:'通知',
		content:'系统错误',
		ok:true,
		cancel:true,
		init:function(){},
		close:function(){}
	};

	Idialog.prototype={
		show:function(){
			var obj=this;
			var _self=this._self;

			_self.width(this.settings.width);
			_self.css('margin-left','-'+this.settings.width/2+'px');
			_self.css('top',this.settings.top+'px');

			_self.find('.idialog_title font').html(this.settings.title);

			if(this.settings.content instanceof $){
				_self.find('.idialog_content').html(this.settings.content.html());
			}else{
				_self.find('.idialog_content').html(this.settings.content);
			}

			if(this.settings.ok===false){
				_self.find('.idialog_ok').remove();
			}
			if(this.settings.cancel===false){
				_self.find('.idialog_cancel').remove();
			}
			
			_self.appendTo("body");
		},

		//对话框本身事件绑定
		bind:function(){
			var obj=this;
			var _self=this._self;
			//关闭
			_self.find('.idialog_title span').click(function(){
				obj.close();
				//_self.remove();
			});
			//取消  待扩展
			_self.find('.idialog_cancel').click(function(){
				obj.close();
				//_self.remove();
			});

			//拖拽绑定
			_self.find(".idialog_title").unbind("mousedown").mousedown(function(e){
				var marginLeft = parseInt( _self.css('marginLeft') ); //margin 左偏移，奇葩css导致
				var offset = $(this).offset();//DIV在页面的位置
				var x = e.pageX - offset.left+marginLeft;//获得鼠标指针离DIV元素左边界的距离
				var y = e.pageY - offset.top+$(document).scrollTop();//获得鼠标指针离DIV元素上边界的距离
				$(document).bind("mousemove",function(ev)//绑定鼠标的移动事件，因为光标在DIV元素外面也要有效果，所以要用doucment的事件，而不用DIV元素的事件
				{
					$(".idialog").stop();//加上这个之后
					var _x = ev.pageX - x;//获得X轴方向移动的值
					var _y = ev.pageY - y;//获得Y轴方向移动的值
					$(".idialog").animate({left:_x+"px",top:_y+"px"},0);
				});
			});
			//拖拽接触	
			$(document).mouseup(function(){
				$(this).unbind("mousemove");
			});

			//确定按钮回调
			_self.find('.idialog_ok').click(function(){
				var res = true;
				if(typeof(obj.settings.ok) == 'function'){
					var r = obj.settings.ok(obj);
					if(r===false) res=false;
				}
				if(res){
					obj.close();
					//_self.remove();	
				}
				
			});
		},
		//内容初始化回调
		init:function(){
			var obj=this;
			return obj.settings.init(obj._self.find('.idialog_content'));
		},
		//关闭弹窗回调
		close:function(){
			var obj=this;
			obj.settings.close(obj._self.find('.idialog_content'))
			obj._self.remove();
		},
	};

	/**
	 * 弹出框  气泡形式
	 * @type {Object}
	 */
	Idialog.tips=function(msg,time){
		if(time==undefined) time=2;
		var tips=$('<div class="idialog_tips">'+msg+'</div>');
		tips.appendTo("body");
		tips.css('margin-left','-'+tips.width()/2+'px');
		setTimeout(function() {
			tips.remove();
		},time*1000);
	};

	win['Idialog']=function(settings){
		var dialog=new Idialog(settings);
		return dialog.run();
	};
	win['Idialog']['tips']=Idialog.tips;


})(jQuery,window,document);