Dropzone.autoDiscover = false;
// Get the template HTML and remove it from the doumenthe template HTML and remove it from the doument
var previewNode = document.querySelector("#template");
previewNode.id = "";
var previewTemplate = previewNode.parentNode.innerHTML;
previewNode.parentNode.removeChild(previewNode);

var myDropzone = new Dropzone(".container", { // Make the whole body a dropzone
    url: "/upload", // Set the url
    thumbnailWidth: 80,
    thumbnailHeight: 80,
    parallelUploads: 50,
    previewTemplate: previewTemplate,
    autoQueue: false, // Make sure the files aren't queued until manually added
    previewsContainer: "#previews", // Define the container to display the previews
    clickable: ".fileinput-button" // Define the element that should be used as click trigger to select files.
});

myDropzone.on("addedfile", function (file) {
    // Hookup the start button
    file.previewElement.querySelector(".start").onclick = function (event) {
        event.preventDefault()
        myDropzone.enqueueFile(file);
    };
});

// Update the total progress bar
myDropzone.on("totaluploadprogress", function (progress) {
    document.querySelector("#total-progress .progress-bar").style.width = progress + "%";
});

myDropzone.on("sending", function (file) {
    // Show the total progress bar when upload starts
    document.querySelector("#total-progress").style.opacity = "1";
    // And disable the start button
    file.previewElement.querySelector(".start").setAttribute("disabled", "disabled");
});

// Hide the total progress bar when nothing's uploading anymore
myDropzone.on("queuecomplete", function (progress) {
    document.querySelector("#total-progress").style.opacity = "0";
    reload_download_list()
});


// Setup the buttons for all transfers
// The "add files" button doesn't need to be setup because the config
// `clickable` has already been specified.
document.querySelector("#actions .start").onclick = function () {
    myDropzone.enqueueFiles(myDropzone.getFilesWithStatus(Dropzone.ADDED));
};
document.querySelector("#actions .cancel").onclick = function () {
    myDropzone.removeAllFiles(true);
    SecondDropzone.removeAllFiles(true)
};


var SecondDropzone = new Dropzone(
    "#second_dropzone", {
    url: "/upload",
    maxFiles: 50,
    parallelUploads: 50,
    dictDefaultMessage: "或是 將檔案與資料夾拖曳至此直接上傳...",
    init: function () {
        this.hiddenFileInput.setAttribute("webkitdirectory", true);
    }
}
)
SecondDropzone.on("addedfile", function (file) {
    console.log(file.type)
})

SecondDropzone.on("dragenter", function () {

    document.getElementById("second_dropzone").style.backgroundColor = "LightGray"
})

SecondDropzone.on("dragleave", function () {
    document.getElementById("second_dropzone").style.backgroundColor = "white"
})

SecondDropzone.on("dragleave", function () {
    document.getElementById("second_dropzone").style.backgroundColor = "white"
})

SecondDropzone.on("dragend", function () {
    document.getElementById("second_dropzone").style.backgroundColor = "white"
})

SecondDropzone.on("success", function () {
    document.getElementById("second_dropzone").style.backgroundColor = "white"
    reload_download_list()
})


// start
function formatBytes(a, b = 2, k = 1024) { with (Math) { let d = floor(log(a) / log(k)); return 0 == a ? "0 Bytes" : parseFloat((a / pow(k, d)).toFixed(max(0, b))) + " " + ["Bytes", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"][d] } }


function getFileIcon(ext) {
    if (ext.length == 0) {
        return '<i class="far fa-file"></i>'
    }
    // <i class="fas fa-file-code"></i>
    switch (ext[0].toLowerCase()) {
        case 'txt':
            return '<span class="iconify" data-icon="grommet-icons:document-txt"></span>'
        case 'pdf':
            return '<i class="far fa-file-pdf"></i>'
        //word 
        case 'docx':
        case 'doc':
            return '<span class="iconify" data-icon="vscode-icons:file-type-word2"></span>'
        //excel
        case 'xls':
        case 'xlsx':
        case 'xlsm':
            return '<span class="iconify" data-icon="vscode-icons:file-type-excel2"></span>'
        //csv
        case 'csv':
            return '<i class="fas fa-file-csv"></i>'
        //tab
        case 'tab':
            return '<span class="iconify" data-icon="carbon:cross-tab"></span>'
        //powerpoint
        case 'ppt':
        case 'pptx':
            return '<span class="iconify" data-icon="vscode-icons:file-type-powerpoint2"></span>'
        //image
        case 'jpeg':
        case 'jpg':
        case 'png':
        case 'svg':
        case 'webp':
        case 'bmp':
        case 'gif':
        case 'img':
            return '<i class="fas fa-images"></i>'
        //video
        case 'avi':
        case 'flv':
        case 'wmv':
        case 'mov':
        case 'mp4':
        case 'mpeg':
        case 'heic':
        case 'heif':
        case 'hevc':
        case 'mts':
            return '<i class="fas fa-video"></i>'
        //music
        case 'wav':
        case 'flac':
        case 'ape':
        case 'alac':
        case 'mp3':
        case 'aac':
        case 'wma':
        case 'rmi':
        case 'wv':
            return '<i class="fas fa-music"></i>'
        //git
        case 'git':
        case 'gitignore':
            return '<span class="iconify" data-icon="fa-brands:git-alt"></span>'
        //python 
        case 'py':
            return '<span class="iconify" data-icon="logos:python"></span>'
        case 'ipynb':
            return '<span class="iconify" data-icon="logos:jupyter"></span>'
        //golang
        case 'go':
            return '<span class="iconify" data-icon="grommet-icons:golang"></span>'
        //matlab
        case 'm':
        case 'mat':
            return '<span class="iconify" data-icon="vscode-icons:file-type-matlab"></span>'
        // c-like
        case 'c':
            return '<span class="iconify" data-icon="logos:c"></span>'
        case 'cpp':
            return '<span class="iconify" data-icon="logos:c-plusplus"></span>'

        // Other languages
        case 'pl':
            return '<span class="iconify" data-icon="logos:perl"></span>'
        case 'rb':
            return '<span class="iconify" data-icon="logos:ruby"></span>'
        case 'php':
            return '<span class="iconify" data-icon="vscode-icons:file-type-php3"></span>'
        case 'jsp':
        case 'jspx':
            return '<span class="iconify" data-icon="logos:java"></span>'
        //js
        case 'js':
            return '<span class="iconify" data-icon="logos:javascript"></span>'
        case 'json':
            return '<span class="iconify" data-icon="vscode-icons:file-type-light-json"></span>'
        //css 
        case 'css':
            return '<span class="iconify" data-icon="logos:css-3"></span>'
        //html:
        case 'html':
        case 'htm':
            return '<span class="iconify" data-icon="logos:html-5"></span>'
        //xml
        case 'xml':
            return '<span class="iconify" data-icon="vscode-icons:file-type-xml"></span>'
        //markdown:
        case 'md':
            return '<span class="iconify" data-icon="vscode-icons:file-type-markdown"></span>'
        //zip
        case 'zip':
        case 'tar':
        case 'gz':
        case 'rar':
        case '7z':
            return '<i class="fas fa-file-archive"></i>'
        //db
        case 'db':
        case 'sql':
        case 'dat':
            return '<span class="iconify" data-icon="vscode-icons:file-type-sql"></span>'
        //bin
        case 'exe':
        case 'bat':
        case 'dll':
        case 'ini':
            return '<span class="iconify" data-icon="vscode-icons:file-type-binary"></span>'
        case 'ttf':
            return '<span class="iconify" data-icon="ant-design:font-size-outlined"></span>'
        case 'cer':
            return '<span class="iconify" data-icon="vscode-icons:file-type-cert"></span>'
        case 'apk':
            return '<span class="iconify" data-icon="flat-color-icons:android-os"></span>'
        default:
            return '<i class="far fa-file"></i>'
    }
}

function delFile(filename) {
    $.ajax({
        type: "POST",
        url: "/delete/" + filename,
        async: false,
        cache: false,
        timeout: 30000,
        fail: function () {
            return true;
        },
        done: function () {
            reload_download_list()
        }
    });
}

function ChangePath(dir) {
    $.ajax({
        type: "POST",
        url: "/ChangePath",
        data: { "subfolder": dir },
        done: function (msg) {
            reload_download_list()
        }
    });
}

function display_file(ele) {
    var tableData = ""
    for (var i = 0; i < ele.length; i++) {
        var create_time = new Date(ele[i]["Time"])
        var filename = ele[i]["Name"]
        var ext = (/[.]/.exec(filename)) ? /[^.]+$/.exec(filename) : [];
        tableData += "<tr class='right_click'>\n"
        tableData += "<td class='fname' style='font-size: 150%; text-align: left;'>" + getFileIcon(ext) + '&ensp;<strong class="strong_title">' + filename + " </strong></td>\n"
        tableData += "<td>" + formatBytes(ele[i]["Size"]) + "</td>\n"
        tableData += "<td> " + create_time.toISOString().substring(0, 19).replace('T', ' '); + " </td>\n"
        new_url = 'downloads/' + filename
        tableData += '<td> <a class="btn btn-success" style="center" \
                                        href="'+ new_url + '" download="' + filename + '" > \
                                        <i class="fas fa-download"></i> \
                                        <span>下載</span> \
                                    </a>'
        tableData += '<a href="/" class="btn btn-danger del_button" style="center" \
                        onclick="delFile('+ "'" + filename + "'" + ')"  > \
                        <i class="fas fa-trash"></i> \
                        <span>刪除</span> \
                    </a>\
                </td>\n'
        tableData += "</tr>\n"
    }
    return tableData
}

function display_folder(ele) {
    var tableData = ""
    for (var i = 0; i < ele.length; i++) {
        var create_time = new Date(ele[i]["Time"])
        var filename = ele[i]["Name"]
        tableData += "<tr class='bg-secondary text-white right_click'>\n"
        tableData += "<td class='fname' style='font-size: 150%; text-align: left;' >" + '<i class="far fa-folder"></i>' + '&ensp;<strong class="strong_title">' + filename + " </strong></td>\n"
        tableData += "<td> 資料夾 </td>\n"
        tableData += "<td> " + create_time.toISOString().substring(0, 19).replace('T', ' '); + " </td>\n"
        new_url = 'downloads/' + filename
        tableData += '<td> <a class="btn btn-success" style="center" \
                                        href="'+ new_url + '" download="' + filename + '" > \
                                        <i class="fas fa-cloud-download-alt"></i> \
                                        <span>下載全部</span> \
                                    </a>'
        tableData += '<a  href="/" class="btn btn-danger del_button" style="center" \
                        onclick="delFile('+ "'" + filename + "'" + ')"  > \
                        <i class="fas fa-trash"></i> \
                        <span>刪除</span> \
                    </a>'
        tableData += '<a  href="/" class="btn btn-warning" style="center" \
                onclick="ChangePath('+ "'" + filename + "'" + ')"  > \
                <i class="fas fa-external-link-alt"></i> \
                <span>進入</span> \
            </a></td>\n'
        tableData += "</tr>\n"
    }
    return tableData
}

current_path = ""
list_folders = []
permission = ""
function reload_download_list() {
    $.ajax({
        type: "GET",
        url: "/ls",
        success: function (all_files) {
            permission=all_files["permission"]
            current_path = all_files["current_path"].replace("\\", "/")
            list_folders = all_files["folder"].map(f => f.Name)
            var tableData = display_folder(all_files["folder"]);
            tableData += display_file(all_files["file"]);
            if (all_files["file"].length == 0 && all_files["folder"].length == 0) {
                tableData += '<td colspan="4" class="align-middle">這個資料夾現在是空的喔</td>'
            }
            $("#tbody1").html(tableData);
            if(all_files["permission"]=="visitor"){
                $(".del_button").hide()
                $("#trash_button").hide()
                $("#setting_button").hide()
            }
            $('.breadcrumb-item').remove();
            if (current_path.length == 0) {
                $("#back_button").hide();
                $("#current_path_bread").hide();
            } else {
                var path_list = current_path.split("/")
                // console.log(path_list)
                for (var i = 0; i < path_list.length; i++) {
                    if (i == path_list.length - 1) {
                        $("#current_path_bread").append('<li class="breadcrumb-item" aria-current="page">' + path_list[i] + '</li>')
                    } else {
                        var action = "Go_abs_Path('" + path_list.slice(0, i + 1).join('/') + "')"
                        $("#current_path_bread").append('<li class="breadcrumb-item"><a href="#" onclick="' + action + '">' + path_list[i] + '</a></li>')
                    }
                }
                $("#back_button").show()
                $("#current_path_bread").show();
            }
        }
    });
}

function Go_abs_Path(dir) {
    $.ajax({
        type: "POST",
        url: "/Go_abs_Path",
        data: { "pathname": dir },
        success: function (msg) {
            reload_download_list()
        }
    });
}

function Go_Back() {
    $.ajax({
        type: "GET",
        url: "/Go_Back",
        success: function (msg) {
            reload_download_list()
        }
    });
}

function GetFreeSpace(){
    $.ajax({
        method:"GET",
        url:"/freespace",
        success:function(data){
            $("#disk_space").css("width",data["percent"])
            var label_text = `儲存空間使用了 ${data['free']} 中的 ${data['use']}， 約為(${data['percent']})`
            $("#disk_space_text").text(label_text);
            var num = parseFloat(data["percent"].slice(0,data["percent"].length-1))
            if(num<60.0){
                $("#disk_space").css("background-color", "green")
            }else if(num>=60 && num<=80){
                $("#disk_space").css("background-color", "yellow")
            }else{
                $("#disk_space").css("background-color", "red")
            }
        }
    });
}

window.onload = function () {
    reload_download_list()
    GetFreeSpace()
}


function CreateFolder() {
    var folder = document.getElementById("input_folder").value
    document.getElementById("input_folder").value = ""
    // var regexss = /[^\w_-]+/
    // if (folder.match(regexss) != null) {
    //     alert("添加的路徑不得含有除了底線與減號外的字元")
    //     return false
    // }
    $.ajax({
        type: "GET",
        url: "/create/" + folder,
        success: function (msg) {
            reload_download_list()
        }
    });
}

// the selected right click zone
SELECT_ZONE = undefined

//  add right click menu to table
$("table").on("DOMSubtreeModified", function () {
    // console.log("permission is",permission)
    if(permission=="" || permission=="visitor"){
        return
    }
    $(".right_click").on('contextmenu', function (e) {
        $('.right_click').css('box-shadow', 'none');
        var top = e.pageY + 10;
        var left = e.pageX + 10;
        $(this).css('box-shadow', 'inset 1px 1px 0px 0px red, inset -1px -1px 1px 1px red');
        $("#menu").css({
            display: "block",
            top: top,
            left: left
        });
        SELECT_ZONE = $(this);
        return false; //blocks default Webbrowser right click menu
    });

    $("body").on("click", function () {
        if ($("#menu").css('display') == 'block') {
            $(" #menu ").hide();
        }
        $('.right_click').css('box-shadow', 'none');
    });

    // $("#menu button").on("click", function () {
    //     $(this).parent().hide();
    // });
})

$(".rename").on("click", function (e) {
    if (SELECT_ZONE != undefined) {
        var new_fname = prompt("輸入新的檔案名稱(不包含副檔名): ");
        if (new_fname != null) {
            var old_name = SELECT_ZONE.children(".fname").children(".strong_title").html().trim();

            // var regexss = /[^\w_-]+/
            // if (new_fname.match(regexss) != null) {
            //     alert("添加的路徑不得含有除了底線與減號外的字元")
            //     return false
            // }
            // else {
                
            // }
            $.ajax({
                type: "POST",
                url: "/rename",
                data: {
                    "oldname": old_name,
                    "newname": new_fname,
                },
                success: function (msg) {
                    reload_download_list();
                },
                fail: function (msg) {
                    alert(msg.responseText);
                }
            });
        }
    }
    SELECT_ZONE = undefined;
})


$(".moveto").on("mouseenter", function () {

    if (SELECT_ZONE != undefined) {
        var fname = SELECT_ZONE.children(".fname").children(".strong_title").html().trim();
        // alert(fname)
        submenu = $(this).parent().children("ul")
        submenu.text("")
        var html_text = '<li><button class="dropdown-item bg-light" onclick="MoveToFolder(' + "'$-parent-$'" + ')">上一層</button></li>'
        submenu.append(html_text)
        for (var i = 0; i < list_folders.length; i++) {
            console.log(fname, list_folders[i])
            if (fname != list_folders[i]) {
                var html_text = '<li><button class="dropdown-item" onclick="MoveToFolder(' + "'" + list_folders[i] + "'" + ')" >' + list_folders[i] + '</button></li>';
                submenu.append(html_text);
            }

        }
    }
})

function MoveToFolder(foldername) {
    if (SELECT_ZONE != undefined) {
        var fname = SELECT_ZONE.children(".fname").children(".strong_title").html().trim();
        $.ajax({
            type: "POST",
            async: false,
            url: "/movetofolder",
            data: {
                "filename": fname,
                "foldername": foldername,
            },
            success: function () {
                reload_download_list()
            },
            error: function (msg) {
                alert(msg.responseText)
            }
        });
    }
}


