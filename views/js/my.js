Dropzone.autoDiscover = false;
// Get the template HTML and remove it from the doumenthe template HTML and remove it from the doument
var previewNode = document.querySelector("#template");
previewNode.id = "";
var previewTemplate = previewNode.parentNode.innerHTML;
previewNode.parentNode.removeChild(previewNode);

var myDropzone = new Dropzone(".container", { // Make the whole body a dropzone
    url: "/", // Set the url
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
    url: "/",
    maxFiles: 50,
    parallelUploads: 50,
    dictDefaultMessage: "或是 將檔案與資料夾拖曳至此直接上傳...",
    init: function () {
        this.hiddenFileInput.setAttribute("webkitdirectory", true);
    }
}
)

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

function delFile(filename) {
    $.ajax({
        type: "POST",
        url: "/delete/" + filename,
        async: false,
        cache: false,
        timeout: 30000,
        fail: function(){
            return true;
        },
        done: function () {
            reload_download_list()
        }
    });
}

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
            return '<span class="iconify" data-icon="vscode-icons:file-type-excel2"></span>'
        //csv
        case 'csv':
            return '<i class="fas fa-file-csv"></i>'
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
            return '<i class="fas fa-file-music"></i>'
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
            return '<i class="fas fa-database"></i>'
        //bin
        case 'exe':
        case 'bat':
        case 'dll':
        case 'dat':
        case 'ini':
            return '<span class="iconify" data-icon="vscode-icons:file-type-binary"></span>'
        default:
            return '<i class="far fa-file"></i>'
    }
}

function ChangePath(dir){
    $.ajax({
        type: "POST",
        url: "/ChangePath",
        data:{"subfolder":dir},
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
        console.log(ext)
        tableData += "<tr>\n"
        tableData += "<td style='font-size: 150%; text-align: left;'><strong>" + getFileIcon(ext) + '&ensp;' + filename + " </strong></td>\n"
        tableData += "<td>" + formatBytes(ele[i]["Size"]) + "</td>\n"
        tableData += "<td> " + create_time.toISOString().substring(0, 19).replace('T', ' '); + " </td>\n"
        new_url = 'downloads/' + filename
        tableData += '<td> <a class="btn btn-success" style="center" \
                                        href="'+ new_url + '" download="' + filename + '" > \
                                        <i class="fas fa-download"></i> \
                                        <span>下載</span> \
                                    </a>'
        tableData += '<a  href="/" class="btn btn-danger" style="center" \
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
        var ext = (/[.]/.exec(filename)) ? /[^.]+$/.exec(filename) : [];
        tableData += "<tr class='bg-secondary text-white'>\n"
        tableData += "<td  style='font-size: 150%; text-align: left;' ><strong>" + '<i class="far fa-folder"></i>' + '&ensp;' + filename + " </strong></td>\n"
        tableData += "<td> 資料夾 </td>\n"
        tableData += "<td> " + create_time.toISOString().substring(0, 19).replace('T', ' '); + " </td>\n"
        new_url = 'downloads/' + filename
        tableData += '<td> <a class="btn btn-success" style="center" \
                                        href="'+ new_url + '" download="' + filename + '" > \
                                        <i class="fas fa-cloud-download-alt"></i> \
                                        <span>下載全部</span> \
                                    </a>'
        tableData +='<a  href="/" class="btn btn-warning" style="center" \
        onclick="ChangePath('+ "'" + filename + "'" + ')"  > \
        <i class="fas fa-external-link-alt"></i> \
        <span>進入</span> \
    </a>'
        tableData += '<a  href="/" class="btn btn-danger" style="center" \
                        onclick="delFile('+ "'" + filename + "'" + ')"  > \
                        <i class="fas fa-trash"></i> \
                        <span>刪除</span> \
                    </a>\
                </td>\n'
        tableData += "</tr>\n"
    }
    return tableData
}

current_path = ""

function reload_download_list() {
    $.ajax({
        type: "GET",
        url: "/ls",
        success: function (all_files) {
            current_path = all_files["current_path"].replace("\\","/")

            var tableData = display_file(all_files["file"]);
            tableData += display_folder(all_files["folder"]);
            if (all_files["file"].length==0 && all_files["folder"].length==0){
                tableData+='<td colspan="4" class="align-middle">這個資料夾現在是空的喔</td>'
            }
            $("#tbody1").html(tableData);
            if(current_path.length==0){
                $("#back_button").hide()
                $("#current_path").html("/");
            }else{
                $("#back_button").show()
                $("#current_path").html(current_path);
            }
        }
    });
}

function Go_Back(){
    $.ajax({
        type: "GET",
        url: "/Go_Back",
        success: function (msg) {
            console.log("good")
            reload_download_list()
        }
    });
}

window.onload = reload_download_list()


function CreateFolder() {
    var folder = document.getElementById("input_folder").value
    var regexss = /[^\w_-]+/
    if(folder.match(regexss)!=null){
        alert("添加的路徑不得含有除了底線與減號外的字元")
        return false
    }
    console.log("The new folder is ", folder)
    $.ajax({
        type: "GET",
        url: "/create/" + folder,
        success: function (msg) {
            reload_download_list()
            console.log(msg)
        }
    });
}

function MoveTo() {
    oldLocation = document.getElementById("")
    newLocation = document.getElementById("")
    $.ajax({
        type: "POST",
        url: "/move",
        data: { "oldLocation": oldLocation, "newLocation": newLocation },
        dataType: "json",
        success: function (msg) {
            reload_download_list()
        }
    });
}
