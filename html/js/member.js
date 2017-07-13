$('#editMember').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var memberId = button.data('member-id');
	var modal = $(this);
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
