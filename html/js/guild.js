$('#editGuild').on('show.bs.modal', function (event) {
	var button = $(event.relatedTarget)
	var guildId = button.data('guild-id')
	//ajax here
	var modal = $(this)
	modal.find('.modal-footer .btn-primary').attr('onclick', "saveGuild("+guildId+")")//"saveGuild("+guildId+")")
})
function saveGuild(guildId)
{
	var modal = $(this)
	$.ajax({
		data:{
			name:modal.find('.modal-body #guild-name').val()
			introduce:modal.find('.modal-body #guild-introduce').val()
		}
		type:'post',
		dataType:'json',
		url:'/guild/save/'+guildId,
		success:function(data) {

		}
	});
}
