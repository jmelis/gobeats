function getDir(path) {
    path = path || ".";

    $.getJSON("/dir/"+encodeURIComponent(path), function(data) {
        var dirs = data["Dirs"];
        var dirs_items = [];

        var files = data["Files"];
        var files_items = [];

        if (dirs != null) {
            $.each(dirs, function(i){
                var dir = dirs[i];
                var uri = encodeURIComponent(dir);
                dirs_items.push("<li><a href='#' data-dir='"+uri+"'>" + dir + "</li>");

            });
            $("#dirs").html($("<ul/>", {html: dirs_items.join("")}));
        }

        if (files != null) {
            $.each(files, function(i){
                var file = files[i];
                var uri = "/stream/" + encodeURIComponent(file);
                files_items.push("<li><a href='" + uri + "'>" + file + "</a></li>");
            });
            $("#files").html($("<ul/>", {html: files_items.join("")}))
        }
    });
}

function play(url) {
    var audio = $("#audio")[0];
    audio.src = url;
    audio.load();
    audio.play();
}

$(document).ready(function(){
    getDir();

    $("#files").on("click","li a", function(e){
        e.preventDefault();
        var href = $(this).attr('href');

        play(href);

    });

    $("#dirs").on("click","li a", function(e){
        e.preventDefault();
        var dir = $(this).attr('data-dir');
        getDir(dir);
    });
});
