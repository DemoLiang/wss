//index.js
//获取应用实例
const app = getApp()

Page({
  data: {
    motto: 'Hello World',
    userInfo: {},
    hasUserInfo: false,
    code: '',
    canIUse: wx.canIUse('button.open-type.getUserInfo')
  },
  //事件处理函数
  bindViewTap: function() {
    wx.connectSocket({
      url: 'ws://127.0.0.1:7777/ws',
      method:"GET",
      success:res => {
        console.log("aaaaaaaaaaaa")
      }
    })
    wx.onSocketOpen(function (res) {
      console.log('WebSocket连接已打开！')
      wx.login({
        success: function (res) {
          if (res.code) {
            //发起网络请求
            app.globalData.code = res.code
            wx.sendSocketMessage({
              data:  res.code
            })
          } else {
            console.log('获取用户登录态失败！' + res.errMsg)
          }
        }
      });
      // wx.sendSocketMessage({
      //   data: 'Hello,World:' + Math.round(Math.random() * 0xFFFFFF).toString(),
      // })
      wx.onSocketMessage(function (res) {
        console.log('收到服务器内容：' + res.data)
      })
    })
    wx.navigateTo({
      url: '../logs/logs'
    })
  },
  onLoad: function () {
    if (app.globalData.userInfo) {
      this.setData({
        userInfo: app.globalData.userInfo,
        hasUserInfo: true
      })
    } else if (this.data.canIUse){
      // 由于 getUserInfo 是网络请求，可能会在 Page.onLoad 之后才返回
      // 所以此处加入 callback 以防止这种情况
      app.userInfoReadyCallback = res => {
        this.setData({
          userInfo: res.userInfo,
          hasUserInfo: true
        })
      }
    } else {
      // 在没有 open-type=getUserInfo 版本的兼容处理
      wx.getUserInfo({
        success: res => {
          app.globalData.userInfo = res.userInfo
          this.setData({
            userInfo: res.userInfo,
            hasUserInfo: true
          })
        }
      })
    }
  },
  getUserInfo: function(e) {
    console.log(e)
    app.globalData.userInfo = e.detail.userInfo
    this.setData({
      userInfo: e.detail.userInfo,
      hasUserInfo: true
    })
    // wxRequest
  },
  // wxRequest:function(e){
  //   wx.request({
  //   url: "127.0.0.1:7777/ws",
  //   data: "aaaaaaaa",
  //   method: "POST",
  //   dataType: "json",
  //   success:res =>{
  //       console.log('aaaaaaaaaa')
  //     }
  //   })
  // }
  bindFormSubmit: function (e) {
    console.log(e.detail.value.textarea)
    wx.sendSocketMessage({
      data: e.detail.value.textarea + app.globalData.code,
      // code:dataContent.code
      success: function (res) { },
      fail: function (res) { },
      complete: function (res) { },
    })
  },
})
