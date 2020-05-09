const express = require('express');
const request = require('request');
const bodyParser = require('body-parser');
var app = express();


var invoke = require('./invoke.js'); // this are .js file which used to communicated with
var query = require('./query.js');
var registerUser = require('./registerUser.js');
app.use(express.static('public')); // make css file static access from easly
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({
  extended: true
}));


var urlbodyParser = bodyParser.urlencoded({
  extended: true
}); // midleware parse the data
app.set('views', __dirname + '/views'); // tell the directory of webpage .ejs extention b/c run JavaScript within html page
app.set('view engine', 'ejs');


app.get('/createuser', function(req, res) {
  callback = function(queryResult) { // this check user account is still active or not
    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')

    if (queryResult != '') {
      arr = queryResult.split(',');
    } else {
      arr = ''
    }
    var data = {
      message: "Fill Form To create new User",
      bossList: arr
    }
    res.render('userCreationForm', {
      data: data
    });
  }

  query.query('users', 'QueryAllBossId', '', callback, 'admin')
}); // this api call is used to create user in system it open user creation form

app.post('/create', function(req, res) {
  console.log(req.body);
  callback = function(queryResult) { // this check user account is still active or not
    if (queryResult == 200) {
      callback = function(queryResult) { // this check user account is still active or not
        if (queryResult == 200) {
          var data = {
            message: "User Successfully Created Login Please"
          }
          res.render('login', {
            data: data
          });
        } else {
          callback = function(queryResult) { // this check user account is still active or not
            queryResult = queryResult.replace(/\"/g, '')
            queryResult = queryResult.replace('[', '')
            queryResult = queryResult.replace(']', '')
            queryResult = queryResult.replace('{', '')
            queryResult = queryResult.replace('}', '')
            arr = queryResult.split(',');

            var data = {
              message: "Error occur Please try again",
              bossList: arr
            }
            res.render('userCreationForm', {
              data: data
            });
          }
          query.query('users', 'QueryAllBossId', '', callback, 'admin')
        }
      } // if key created Successfully then add user details into blockchain
      invoke.createUser('users', 'CreateUser', req.body.Name, req.body.userid, req.body.password, req.body.role, req.body.boss, callback, 'admin')
    } else {
      res.send('User Already exit or some other error occur')
    }

  }
  registerUser.registerUser(req.body.userid, callback) // first we create public and private keys of used
}); // this will received request and data to create user


app.get('/login', function(req, res) {
  var data = {
    message: ""
  }
  res.render('login', {
    data: data // load login page
  });

});



app.post('/login', function(req, res) {
  callback = function(queryResult) { // this check user account is still active or not
    if (queryResult = 'true') {
      callback = function(queryResult) { // this check user account is still active or not
        var obj = JSON.parse(queryResult);
        console.log(obj);
        if (obj.role == 'Boss') { // here we check the role of user and redirect them to their respective page
          var data = {
            message: "",
            obj: queryResult
          }
          res.render('bossMainPage', {
            data: data
          });
        } else {
          var data = {
            message: "",
            obj: queryResult
          }
          res.render('employeeMainPage', {
            data: data
          });
        }
      }
      query.getUserDetails('users', 'QueryUser', req.body.userid, callback)
    } else {
      var data = {
        message: "Worong id password"
      }
      res.render('login', {
        data: data
      });
    }
  }
  console.log(req.body);

  query.login('users', 'LogIn', req.body.password, callback, req.body.userid) // here we query to blockchain and run SmartContract to check user id and password

});




app.post('/createTransaction', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  var data = {
    tranCreated: obj.userid,
    validated: obj.whoboss,
    obj: req.body.obj
  }
  res.render('createTransaction', {
    data: data // here we load forrm fro create tran in system
  });

});


app.post('/transfer', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not
    if (queryResult == 200) {
      var data = {
        message: "Successfully Transfer",
        obj: req.body.obj
      }
      res.render('employeeMainPage', {
        data: data
      });
    } else {
      var data = {
        message: "Failed Transfer",
        obj: req.body.obj
      }
      res.render('employeeMainPage', {
        data: data
      });
    }
  } // here we put details of tran into blockchain and mark as unvalidated tran
  invoke.transfer("tran", "CreateTransaton", req.body.to, req.body.from, req.body.date, req.body.amount, obj.userid, obj.whoboss, callback, obj.userid)
});



app.post('/ViewAllTransactonCreatedByMe', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not

    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')
    if (queryResult != '') {
      arr = queryResult.split(',');
    } else {
      arr = '';
    }

    var data = {
      arr: arr,
      obj: req.body.obj

    }
    res.render('ViewAllTransactonCreatedByMe', {
      data: data
    });
  }
  query.getAllTransacionCreatedByMe('tran', 'QueryAllTransationCreatedByMe', callback, obj.userid)
}); // here we get all the tran create by login user

app.post('/ViewAllTransaction', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not
    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')
    arr = queryResult.split(',');
    let transactionid = [];
    let permission = [];
    if (queryResult != '') {
      for (var i = 0; i < arr.length; i++) {
        values = arr[i].split(":");
        transactionid.push(values[0])
        permission.push(values[1])
      }
    }
    var data = {
      transactionid: transactionid,
      permission: permission,
      obj: req.body.obj
    }
    res.render('ViewAllTransaction', {
      data: data
    });
  }
  query.ViewAllTransaction('tran', 'QueryAllTransationInCompany', obj.whoboss, callback, obj.userid)
}); // here we get all the tran under one boss or company


app.post('/requestForViewDetails', function(req, res) {

  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    if (queryResult == 200) {
      var data = {
        message: "Request Has Been submitted",
        obj: values[0]
      }
      res.render('employeeMainPage', {
        data: data
      });
    } else {
      var data = {
        message: "Request Submission Problem Try Again",
        obj: values[0]
      }
      res.render('employeeMainPage', {
        data: data
      });
    }
  }
  invoke.createViwRequest("tran", "CreateViwRequest", values[1], obj.whoboss, obj.userid, callback, obj.userid)
}); // here we get details of transation if we have permission

app.post('/viewTransationDetailsAfterRequest', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      result: queryResult,
      obj: values[0]
    }
    res.render('viewDetails', {
      data: data
    });
  }
  query.queryTransation("tran", "QueryTransation", values[1], callback, obj.userid)
}); // here we sendrequest to view transation in systm


app.get('/logout', function(req, res) { // student login page
  var data = {
    message: ""
  }
  res.render('login', {
    data: data
  });
  //  res.render('user');
});

//////////////////////////////////////////////////////////////////////////////////////////Boss api Call /////////////////////////////////////////////////
app.post('/queryAllUnvalidatedTransationId', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not
    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')
    if (queryResult != '') {
      transactionid = queryResult.split(',');
    } else {
      transactionid = ''
    }
    var data = {
      transactionid: transactionid,
      obj: req.body.obj
    }
    res.render('queryAllUnvalidatedTransationId', {
      data: data
    });
  }
  query.queryAllUnvalidatedTransationId('tran', 'QueryAllUnvalidatedTransationId', callback, obj.userid)
}); // here we get all UnValidated tran// here we get all the UnValidated tran id and send tem to html page



app.post('/ApprovedRequest', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      message: 'Request ' + values[1] + ' Successfully Approved',
      obj: values[0]
    }
    res.render('bossMainPage', {
      data: data
    });
  }
  invoke.requestProcess("tran", "RequestProcess", values[1], 'allow', callback, obj.userid)
}); // here we sendrequest to view transation in systm  // here we set permission of tran View



app.post('/RejecedRequest', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      message: 'Request ' + values[1] + ' Successfully Rejected',
      obj: values[0]
    }
    res.render('bossMainPage', {
      data: data
    });
  }
  invoke.requestProcess("tran", "requestProcess", values[1], 'deny', callback, obj.userid)
}); // here we sendrequest to view transation in systm// here we reject request permission of View tran



app.post('/RejectTransation', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      message: 'Transaction Rejected',
      obj: values[0]
    }
    res.render('bossMainPage', {
      data: data
    });
  }
  invoke.validateTransation("tran", "ValidateTransation", values[1], 'discard', callback, obj.userid)
}); // here set Transaction Reject    // here we reject Transaction


app.post('/queryAllValidatedTransationIdByMe', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not
    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')
    console.log("queryAllValidatedTransationIdByMe" ,queryResult);
    if (queryResult != '') {
      transactionid = queryResult.split(',');
    } else {
      transactionid = ''
    }
    var data = {
      transactionid: transactionid,
      obj: req.body.obj
    }
    res.render('queryAllValidatedTransationIdByMe', {
      data: data
    });
  }
  query.queryAllUnvalidatedTransationId('tran', 'QueryAllValidatedTransationIdByMe', callback, obj.userid)
}); // here we get all the Validated  by Me // here we get all tran Validated by login user



app.post('/ViewTransationDetails', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      result: queryResult,
      obj: values[0]
    }
    res.render('ViewTransationDetails', {
      data: data
    });
  }
  query.queryTransation("tran", "QueryTransation", values[1], callback, obj.userid)
}); // here we sendrequest to view transation in systm    // here we get details of transation




app.post('/queryAllRequestToViewTransion', function(req, res) {
  var obj = JSON.parse(req.body.obj);
  callback = function(queryResult) { // this check user account is still active or not
    queryResult = queryResult.replace(/\"/g, '')
    queryResult = queryResult.replace('[', '')
    queryResult = queryResult.replace(']', '')
    queryResult = queryResult.replace('{', '')
    queryResult = queryResult.replace('}', '')
    arr = queryResult.split(',');
    let requestid = [];
    let transactionid = [];
    console.log("arr", arr);
    if (queryResult != '') {
      for (var i = 0; i < arr.length; i++) {
        values = arr[i].split(":");
        requestid.push(values[0])
        transactionid.push(values[1])
      }
    } else {
      console.log("here");
      requestid = '';
      transactionid = '';
    }

    var data = {
      requestid: requestid,
      transactionid: transactionid,
      obj: req.body.obj
    }
    res.render('queryAllRequestToViewTransion', {
      data: data
    });
  }
  query.queryAllUnvalidatedTransationId('tran', 'QueryAllRequestToViewTransion', callback, obj.userid)
}); // here we get all Request of employee to view the details of transation    // here we get request id which is issued by different user to  View particular Transaction



app.post('/ApprovedTransation', function(req, res) {
  values = req.body.data.split('~') // we conncatenat object of user data and transactionid b/c we don't send two hidden type data from one form
  var obj = JSON.parse(values[0]);
  callback = function(queryResult) { // this check user account is still active or not
    var data = {
      message: 'Transaction Successfully Approved',
      obj: values[0]
    }
    res.render('bossMainPage', {
      data: data
    });
  }
  invoke.validateTransation("tran", "ValidateTransation", values[1], 'true', callback, obj.userid)
}); // here we sendrequest to view transation in systm  // here we Approved Transaction Create by different users


////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

app.listen(3000);
