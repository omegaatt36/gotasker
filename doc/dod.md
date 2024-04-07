Definition of Done
===

需要以 Golang 完成以下四個 RESTful API：

1. `GET /tasks`：獲取任務清單
1. `POST /tasks/`：創建任務
1. `PUT /tasks/{id}`：修改任務
1. `DELETE /tasks/{id}`：刪除任務

任務詳細資訊：

- name
  type: string
  description：task name
- status
  type: integer
  enum：[0,1]
  description：0 represents an incomplete task, while 1 represents a completed task

假設：

1. 需求中沒有提到 name 不可重複，則假設可以重複。
1. 需求中沒有提到 name 與 status 可以是 nullable，則假設只會是 enum 0,1。
1. 需求中沒有提到獲取任務清單的排序，假設使用創建時間升冪作為排序。
1. 需求中沒有提到 tasks id 的型別是 string 或是 numeric，假設使用 unit 作為非負整數 identification。
