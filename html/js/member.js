$('#editMember').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var memberId = button.data('member-id');
	var modal = $(this);
	if (0 != memberId) {
		$.ajax({
			url:'/guild/msg/member-info',
			contentType:'application/json',
			data:JSON.stringify({Id:memberId}),
			type:'post',
			dataType:'json',
			success:function(data) {
				modal.find('#editMemberLabel').text("修改成员");
				modal.find('#member-name').val(data.Name);
				modal.find('#member-mobile').val(data.Mobile);
				modal.find('#member-ability').val(data.Ability);
				modal.find('#member-role').val(data.Role);
			}
		});
	} else {
		modal.find('#editMemberLabel').text("新成员");
		modal.find('#member-mobile').val("");
		modal.find('#member-name').val("");
		modal.find('#member-ability').val("");
		modal.find('#member-role').val(0);
	}
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveMember('"+memberId+"')");
})
function saveMember(memberId)
{
	var data = {
		Id:memberId,
		Mobile:$('#editMember .modal-body #member-mobile').val(),
		Name:$('#editMember .modal-body #member-name').val(),
		Ability:$('#editMember .modal-body #member-ability').val(),
		Role:parseInt($('#editMember .modal-body #member-role').val()),
		GuildId:$('body .main .page-header span').text()
	};
	$.ajax({
		url:'/guild/msg/member-save',
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
		url:'/guild/msg/member-delete',
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

function sendCode()
{
	$.ajax({
		url:'/guild/msg/get-mobile-code',
		contentType:'application/json',
		headers:{SessionId:localStorage.sessionId},
		data:JSON.stringify({Mobile:$('#login .modal-body #mobile').val()}),
		type:'post',
		dataType:'json',
		success:function(data) {
			// setTimeout()
			var secs = 60;
			var button = $('#login .modal-body #mobile').next().find('button');
			button.attr('disabled', 'disabled');
			button.text(secs+'秒后重发');
			var t = setInterval(function(){
				button.text((--secs)+'秒后重发');
				if (0 == secs) {
					clearInterval(t);
					button.removeAttr('disabled');
					button.text('发送验证码');
				}
			}, 1000);
		}
	});
}
function userLogin()
{
	$.ajax({
		url:'/guild/msg/login',
		contentType:'application/json',
		headers:{SessionId:localStorage.sessionId},
		data:JSON.stringify({VerifyCode:parseInt($('#login .modal-body #verify-code').val())}),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result) {
				$('#login').modal('hide');
				loginButton2out();
				window.location.reload();
			}
		}
	});
}
function userLogout()
{
	$.ajax({
		url:'/guild/msg/logout',
		contentType:'application/json',
		headers:{SessionId:localStorage.sessionId},
		data:JSON.stringify({Mobile:"11"}),
		type:'post',
		dataType:'json',
		success:function(data) {
			if (data.Result)
				window.location.reload();
				// loginButton2in();
		}
	})
}
function loginButton2out()
{
	var loginButton = $('body .sidebar button.btn-warning');
	loginButton.removeAttr('data-toggle').removeAttr('data-target');
	loginButton.attr('onclick', 'userLogout()');
	loginButton.text('登出');
}
function loginButton2in()
{
	var loginButton = $('body .sidebar button.btn-warning');
	loginButton.removeAttr('onclick');
	loginButton.attr('data-toggle', 'modal').attr('data-target', '#login');
	loginButton.text('登入/注册');
}
function roleDoc()
{

}
function roleChgDoc(role) //role able!
{
	//nil role
	if (role) {
		//RoleSysAdmin
		if (1 & role) {

		} else {

		}
		//RoleMaster
		if (2 & role) {

		} else {

		}
		//RoleAdmin
		if (4 & role) {

		} else {

		}
	} else {

	}
}
