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
		}
	});
}
function showTable(type)
{
	if (1 == type) {
		url = '/member/list';
	} else {
		url = '/task/list';
	}
	guildId = $('body .main .page-header span').text();
	$.ajax({
		url:url,
		data:{Id:guildId},
		type:'get',
		dataType:'html',
		success:function(data) {
			$('body .main .table-responsive').html(data);
		}
	})
}

$(function() {
	// $('body .sidebar ul li:first-child').trigger('click');
	$('body .sidebar ul li:first-child').attr('class', 'active');
	$('body .sidebar ul li.active a').trigger('click');
})
