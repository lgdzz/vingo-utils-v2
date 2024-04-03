#!/bin/bash

# 定义单选选项数组
options=("提交代码" "发布版本" "删除版本")

# 显示单选选项
echo "请选择执行方法:"
select opt in "${options[@]}"; do
  case $opt in
    "提交代码")
      action="push"
      break
      ;;
    "发布版本")
      action="addTag"
      read -p "请输入要发布的版本号，按回车确认: " version
      break
      ;;
    "删除版本")
      action="delTag"
      read -p "请输入要删除的版本号，按回车确认: " version
      break
      ;;
    *) echo "Invalid option";;
  esac
done

if [ "$action" = "push" ]
then
  git add .
  git commit -m "-"
  git push
fi

if [ "$action" = "addTag" ]
then
  git tag -a $version -m '-'
  git push origin $version
fi

if [ "$action" = "delTag" ]
then
  git tag -d $version
  git push --delete origin $version
fi
