$('#editGuild').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget);
	var guildId = button.data('guild-id');
	//ajax here
	var modal = $(this);
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveGuild('"+guildId+"')");//"saveGuild("+guildId+")")
})
function saveGuild(guildId)
{
	var modal = $(this);
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
			$('#editGuild').modal('hide');
		}
	});
}
