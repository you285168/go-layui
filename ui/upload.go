package ui

type UIUpload struct {
	*ElemBase
	Text     string
	OnUpload UploadFile
}

type UploadFile func(username string, param map[string]string, filename string, data []byte)

var HtmlUpload = `<div class="layui-upload layui-form-item">
<button type="button" class="layui-btn" id="{{.Id}}"><i class="layui-icon"></i>{{.Text}}</button>
<script>
	layui.use(['upload'], function () {
	     var $ = layui.jquery,
	         upload = layui.upload;

		upload.render({
			elem: '#{{.Id}}'
			,url: '/api/upload?event_id={{.Id}}&url_router={{.Rout}}'
			,accept: 'file' //普通文件
			//,exts: 'zip|rar|7z' //只允许上传压缩文件
			,choose: function(obj){
				obj.preview(function(index, file, result){
					$("#{{.Id}}").val(file.name);
				});
			}
			,progress: function(n, elem){
				var percent = n + '%'; //获取进度百分比
				$("#{{.Id}}").text("正在上传..."+percent);
			}
			,done: function(res){
				handleRsp(res);
			}
			,error: function(index, upload){
				$("#{{.Id}}").text("上传失败");
			}
		});
	});
</script>
</div>`

func newUpload(e *ElemBase, text string, onupload UploadFile) *UIUpload {
	b := &UIUpload{e, text, onupload}
	b.self = b
	return b
}

func NewUpload(text string, onupload UploadFile) *UIUpload {
	return newUpload(newElem("upload", HtmlUpload), text, onupload)
}

func (u *UIUpload) Clone() HtmlElem {
	nu := newUpload(cloneElem(u.Id, "upload", HtmlUpload), u.Text, u.OnUpload)
	nu.ElemBase.clone(u.ElemBase)
	return nu
}
