$('#editGuild').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var guildId = button.data('guild-id');
	var modal = $(this);
	//ajax here
	if (0 != guildId) {
		$.ajax({
			url:'/guild/msg/guild-info',
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
		url:'/guild/msg/guild-save',
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
		url:'/guild/msg/guild-delete',
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
		headers:{SessionId:localStorage.sessionId},
		data:{Id:guildId},
		type:'get',
		dataType:'html',
		success:function(data) {
			$('body .main').html(data);
			$('body .sidebar ul li.active').removeClass('active');
			$('body .sidebar ul li a[onclick*="'+guildId+'"]').parent().addClass('active');
			$('body .main ul.nav-tabs li.active a').trigger('click');

			var role = parseInt($('body .main #role').text());
			if (role & 2) {
				$('#editMember #member-role').parent().removeClass('hidden').addClass('show');
			} else {
				$('#editMember #member-role').parent().removeClass('show').addClass('hidden');
			}
		}
	});
}
function showTable(type)
{
	if (1 == type) {
		url = '/guild/member';
	} else if (2 == type) {
		url = '/guild/task';
	} else {
		return;
	}
	guildId = $('body .main .page-header span').text();
	$.ajax({
		url:url,
		headers:{SessionId:localStorage.sessionId},
		data:{Id:guildId},
		type:'get',
		dataType:'html',
		success:function(data) {
			$('body .main .table-responsive').html(data);
			$('body .main ul.nav-tabs li.active').removeClass('active');
			$('body .main ul.nav-tabs li a[onclick*="'+type+'"]').parent().addClass('active');
		}
	});
}

function applySess()
{
	if (localStorage.sessionId) {
		$.ajax({
			url:'/guild/msg/get-account',
			contentType:'application/json',
			headers:{SessionId:localStorage.sessionId},
			data:JSON.stringify({
				SessionId:parseInt(localStorage.sessionId),
				GuildId:$('body .main .page-header span').text() //取不到
			}),
			type:'post',
			dataType:'json',
			success:function(data) {
				if ('' != data.AccountId) {
					loginButton2out();
				}
				if (data.Role & 1) {
					$('body .sidebar ul li:last-child').removeClass('hidden').addClass('show');
				} else {
					$('body .sidebar ul li:last-child').removeClass('show').addClass('hidden');
				}
			}
		});
	} else {
		$.ajax({
			url:'/guild/msg/apply-session',
			contentType:'application/json',
			data:JSON.stringify({Device:'test'}),
			type:'post',
			dataType:'json',
			success:function(data) {
				localStorage.sessionId = data.SessionId;
			}
		});
	}
}

$(function() {
	$('body .sidebar ul li:first-child').attr('class', 'active');
	$('body .sidebar ul li.active a').trigger('click');
	applySess()
})
