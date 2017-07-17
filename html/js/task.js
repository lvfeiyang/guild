$('#editTask').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var taskId = button.data('task-id');
	var modal = $(this);
	if (0 != taskId) {
		$.ajax({
			url:'/msg/task-info',
			contentType:'application/json',
			data:JSON.stringify({Id:taskId}),
			type:'post',
			dataType:'json',
			success:function(data) {
				modal.find('#editTaskLabel').text("修改任务");
				modal.find('#task-price').val(data.Price);
				modal.find('#task-description').text(data.Desc);
			}
		});
	} else {
		modal.find('#editTaskLabel').text("新任务");
		modal.find('#task-description').text("");
		modal.find('#task-price').val(0);
	}
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveTask('"+taskId+"')");
})
function saveTask(taskId)
{
	var data = {
		Id:taskId,
		Price:parseInt($('#editTask .modal-body #task-price').val()),//Number
		Desc:$('#editTask .modal-body #task-description').val(),
		GuildId:$('body .main .page-header span').text()
	};
	$.ajax({
		url:'/msg/task-save',
		contentType: 'application/json',
		data:JSON.stringify(data),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				$('#editTask').modal('hide');
		}
	});
}
function deleteTask(taskId)
{
	$.ajax({
		url:'/msg/task-delete',
		contentType:'application/json',
		data:JSON.stringify({Id:taskId}),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
		}
	});
}