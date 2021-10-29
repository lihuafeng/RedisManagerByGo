/*
 * @Description: 
 * @Author: gphper
 * @Date: 2021-10-29 19:38:18
 */
var sec1 = `<li class="list-group-item" style="padding-left:{pleft}px !important;">
    <div class="float-left" style="margin-top:5px;width: 200px">
          <i class="bi bi-key">&nbsp;</i>{name}
    </div>
    <button type="button" tag="{tag}" class="add btn btn-primary btn-sm">查看</button>
    <button type="button" tag="{tag}" class="add btn btn-danger btn-sm">删除</button>
</li>
`
    var sec2 = `<li data-toggle="collapse" href="#{id}" class="list-group-item father-node" style="padding-left:{pleft}px !important;">
        <div class="float-left" style="margin-top:5px;width: 200px">
            <i class="bi bi-folder">&nbsp;</i>{name}
        </div>
        <button type="button" tag="{tag}" class="add btn btn-primary btn-sm">添加</button>
        <button type="button" tag="{tag}" class="add btn btn-danger btn-sm">删除</button>
</li>
<div id="{id}" class="collapse">
    {item}
</div>
`
function temp(template, json) {
    var pattern = /\{(\w*[:]*[=]*\w+)\}(?!})/g;
    return template.replace(pattern, function (match, key, value) {
        return json[key];
    });
}

function loadHandleFunc(){
    $(".father-node").click(function(){
        if($(this).children("div").children('i').attr('class') == "bi bi-folder2-open"){
            $(this).children("div").children('i').attr('class',"bi bi-folder")
        }else{
            $(this).children("div").children('i').attr('class',"bi bi-folder2-open")
        }
    })

    $(".add").click(function (event){
        alert("hello");
        event.stopPropagation();
    })
}

function htmlParse(node,layer){
    layer++;
    var html = ``
    for(var i = 0;i < node.length;i++){
        if(node[i].children != null){
            console.log(node[i].title)
            var item = htmlParse(node[i].children,layer)
            html += temp(sec2,{name: node[i].title,id:node[i].all.replaceAll(":","1"),item:item,pleft:layer*20,tag:node[i].all})
        }else{
            console.log(node[i].title)
            html += temp(sec1,{name:node[i].title,pleft:layer*20,tag:node[i].all});
        }
    }
    return html
}


function showTree(id,data){
    html = htmlParse(data,0)
    $("#"+id).html(html)
    loadHandleFunc()
}