$('#editGuild').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var guildId = button.data('guild-id');
	var modal = $(this);
	//ajax here
	if (0 != guildId) {
		$.ajax({
			url:'/msg/guild-info',
			contentType:'application/json',
			data:JSON.stringify({Id:guildId}),
			type:'post',
			dataType:'json',
			success:function(data) {
				modal.find('#editGuildLabel').text(data.Name);
				modal.find('#guild-name').val(data.Name);
				modal.find('#guild-introduce').text(data.Introduce);
			}
		});
	} else {
		modal.find('#editGuildLabel').text("新公会");
		modal.find('#guild-name').val("");
		modal.find('#guild-introduce').text("");
	}
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveGuild('"+guildId+"')");//"saveGuild("+guildId+")")
})
function saveGuild(guildId)
{
	var data = {
		Id:guildId,
		Name:$('#editGuild .modal-body #guild-name').val(),
		Introduce:$('#editGuild .modal-body #guild-introduce').val()
	};
	$.ajax({
		url:'/msg/guild-save',
		contentType: 'application/json',
		data:JSON.stringify(data),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				$('#editGuild').modal('hide');
		}
	});
}
function deleteGuild(guildId)
{
	$.ajax({
		url:'/msg/guild-delete',
		contentType:'application/json',
		data:JSON.stringify({Id:guildId}),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
		}
	});
}
function showOnRight(guildId)
{
	$.ajax({
		url:'/guild/detail',
		data:{Id:guildId},
		type:'get',
		dataType:'html',
		success:function(data) {
			$('body .main').html(data);
			$('body .sidebar ul li.active').removeClass('active');
			$('body .sidebar ul li a[onclick*="'+guildId+'"]').parent().addClass('active');
			$('body .main ul.nav-tabs li.active a').trigger('click');
		}
	});
}
function showTable(type)
{
	if (1 == type) {
		url = '/member';
	} else if (2 == type) {
		url = '/task';
	} else {
		return;
	}
	guildId = $('body .main .page-header span').text();
	$.ajax({
		url:url,
		data:{Id:guildId},
		type:'get',
		dataType:'html',
		success:function(data) {
			$('body .main .table-responsive').html(data);
			$('body .main ul.nav-tabs li.active').removeClass('active');
			$('body .main ul.nav-tabs li a[onclick*="'+type+'"]').parent().addClass('active');
		}
	})
}
$('#editTask').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var taskId = button.data('task-id');
	var modal = $(this);
	//ajax here
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
$('#editMember').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var memberId = button.data('member-id');
	var modal = $(this);
	//ajax here
	if (0 != memberId) {
		$.ajax({
			url:'/msg/member-info',
			contentType:'application/json',
			data:JSON.stringify({Id:memberId}),
			type:'post',
			dataType:'json',
			success:function(data) {
				modal.find('#editMemberLabel').text("修改成员");
				modal.find('#member-mobile').val(data.Mobile);
			}
		});
	} else {
		modal.find('#editMemberLabel').text("新成员");
		modal.find('#member-mobile').val("");
	}
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveMember('"+memberId+"')");
})
function saveMember(memberId)
{
	var data = {
		Id:memberId,
		Mobile:$('#editMember .modal-body #member-mobile').val(),
		GuildId:$('body .main .page-header span').text()
	};
	$.ajax({
		url:'/msg/member-save',
		contentType: 'application/json',
		data:JSON.stringify(data),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				$('#editMember').modal('hide');
		}
	});
}
function deleteMember(memberId)
{
	$.ajax({
		url:'/msg/member-delete',
		contentType:'application/json',
		data:JSON.stringify({Id:memberId}),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
		}
	});
}
$(function() {
	// $('body .sidebar ul li:first-child').trigger('click');
	$('body .sidebar ul li:first-child').attr('class', 'active');
	$('body .sidebar ul li.active a').trigger('click');
})
