package routes

var allAddTemplate = `
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>Document</title>
    <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/css/bootstrap.min.css" integrity="sha384-/Y6pD6FV/Vv2HJnA6t+vslU6fwYXjCFtcEpHbNJ0lyAFsXTsjBbfaDjzALeQsN6M" crossorigin="anonymous">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/4.7.0/css/font-awesome.min.css">
    <style>
        .full-form{
            margin-top: 20px;
        }
        .panel-heading p{
            border-bottom: 1px solid #ccc;
            text-align: center;
            margin-top: 5px;
            margin-bottom: 5px;
        }
        label.width{
            margin-left: 5px;
        }
        select.width{
            margin-left: 5px;
            width: 80%;
            height: 35px !important;
            margin-bottom: 10px;
            font-size: 13px !important;
        }
        label.count{
            margin-left: 5px;
        }
        select.count{
            margin-left: 5px;
            width: 80%;
            height: 35px !important;
            margin-bottom: 10px;
            font-size: 13px !important;
        }
        button{
            cursor: pointer;
        }
    </style>
</head>
<body>
<div class="container">
    <div class="row">
        <div class="col-md-12">
            <form method="post" action="/allads" class="full-form" id="my-form">

                <div class="form-group">
                    <label for="tid">TID</label>
                    <input type="text" class="form-control" id="tid" placeholder="user tid" name="tid">
                </div>

                <div class="form-group">
                    <label for="tid">IP</label>
                    <input type="text" class="form-control" id="ip" placeholder="example : 46.8.98.104" name="ip">
                </div>

                <div class="form-group">
                    <label for="user_agent">UserAgent</label>
                    <input type="text" class="form-control" id="user_agent" placeholder="example : Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36" name="user_agent">
                </div>

                <div class="form-group">
                    <label for="public_id">Publisher Pub id</label>
                    <input type="text" class="form-control" id="public_id" name="public_id">
                </div>

                <div class="form-group">
                    <label for="ad_type">Ad Type</label>
                    <select name="ad_type" id="ad_type" class="form-control">
                        <option value="web">web</option>
                        <option value="app">app</option>
                        <option value="native">native</option>
                    </select>
                </div>

                <div class="native-container" style="display: none">
                    <div class="form-group ad_count-container">
                        <label for="ad_count">Ad count</label>
                        <input type="number" class="form-control" id="ad_count" placeholder="ad count" name="ad_count">
                    </div>
                </div>


                <div class="form-group">
                    <label for="pack">Publisher (domain/package)</label>
                    <input type="text" class="form-control" id="pack" placeholder="domain or package" name="pack">
                </div>
                <div class="slot-container" style="display: none">
                    <div class="row">
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 1</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width1" class="width">size</label>
                                    <select class="form-control width" id="width1" placeholder="width" name="width1">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count1" class="count">count</label>
                                    <select class="form-control count" id="count1" placeholder="count" name="count1">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>

                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 2</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width2" class="width">size</label>
                                    <select class="form-control width" id="width2" placeholder="width" name="width2">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count2" class="count">count</label>
                                    <select class="form-control count" id="count2" placeholder="count" name="count2">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 3</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width3" class="width">size</label>
                                    <select class="form-control width" id="width3" placeholder="width" name="width3">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count3" class="count">count</label>
                                    <select class="form-control count" id="count3" placeholder="count" name="count3">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 4</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width4" class="width">size</label>
                                    <select class="form-control width" id="width4" placeholder="width" name="width4">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count4" class="count">count</label>
                                    <select class="form-control count" id="count4" placeholder="count" name="count4">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 5</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width5" class="width">size</label>
                                    <select class="form-control width" id="width5" placeholder="width" name="width5">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count5" class="count">count</label>
                                    <select class="form-control count" id="count5" placeholder="count" name="count5">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 6</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width6" class="width">size</label>
                                    <select class="form-control width" id="width6" placeholder="width" name="width6">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count6" class="count">count</label>
                                    <select class="form-control count" id="count6" placeholder="count" name="count6">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="row" style="margin-top: 20px; margin-bottom: 20px;">
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 7</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width7" class="width">size</label>
                                    <select class="form-control width" id="width7" placeholder="width" name="width7">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count7" class="count">count</label>
                                    <select class="form-control count" id="count7" placeholder="count" name="count7">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 8</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width8" class="width">size</label>
                                    <select class="form-control width" id="width8" placeholder="width" name="width8">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count8" class="count">count</label>
                                    <select class="form-control count" id="count8" placeholder="count" name="count8">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 9</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width9" class="width">size</label>
                                    <select class="form-control width" id="width9" placeholder="width" name="width9">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count9" class="count">count</label>
                                    <select class="form-control count" id="count9" placeholder="count" name="count9">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 10</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width10" class="width">size</label>
                                    <select class="form-control width" id="width10" placeholder="width" name="width10">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count10" class="count">count</label>
                                    <select class="form-control count" id="count10" placeholder="count" name="count10">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 11</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width11" class="width">size</label>
                                    <select class="form-control width" id="width11" placeholder="width" name="width11">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count11" class="count">count</label>
                                    <select class="form-control count" id="count11" placeholder="count" name="count11">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                        <div class="col-md-2">
                            <div class="card panel-primary">
                                <div class="panel-heading">
                                    <p class="panel-title">Slot 12</p>
                                    <span class="pull-right clickable"><i class="glyphicon glyphicon-chevron-up"></i></span>
                                </div>
                                <div class="panel-body">
                                    <label for="width12" class="width">size</label>
                                    <select class="form-control width" id="width12" placeholder="width" name="width12">
                                        <option value=""></option>
                                        <option value="120x600">120x600</option>
                                        <option value="160x600">160x600</option>
                                        <option value="300x250">300x250</option>
                                        <option value="336x280">336x280</option>
                                        <option value="468x60">468x60</option>
                                        <option value="728x90">728x90</option>
                                        <option value="120x240">120x240</option>
                                        <option value="320x50">320x50</option>
                                        <option value="800x440">800x440</option>
                                        <option value="300x600">300x600</option>
                                        <option value="970x90">970x90</option>
                                        <option value="970x250">970x250</option>
                                        <option value="250x250">250x250</option>
                                        <option value="300x1050">300x1050</option>
                                        <option value="320x480">320x480</option>
                                        <option value="480x320">480x320</option>
                                        <option value="128x128">128x128</option>
                                    </select>
                                    <label for="count12" class="count">count</label>
                                    <select class="form-control count" id="count12" placeholder="count" name="count12">
                                        <option value=""></option>
                                        <option value=1>1</option>
                                        <option value=2">2</option>
                                        <option value="3">3</option>
                                        <option value="4">4</option>
                                        <option value="5">5</option>
                                        <option value="6">6</option>
                                        <option value="7">7</option>
                                        <option value="8">8</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <button type="submit" class="btn btn-primary btn-lg btn-block" style="margin-bottom: 20px;">FETCH<i></i></button>
            </form>

        </div>
    </div>

    <div class="row">
        <div class="col-md-12 result">
        </div>
    </div>
</div>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.11.0/umd/popper.min.js" integrity="sha384-b/U6ypiBEHpOf/4+1nzFpr53nxSS+GLCkfwBdFNTxtclqqenISfwAzpKaMNFNmj4" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-beta/js/bootstrap.min.js" integrity="sha384-h0AbiXch4ZDo7tp9hKZ4TsHbi047NrKGLO3SEJAg45jXxnGIfYzk4Si90RDIqNm1" crossorigin="anonymous"></script>
<script>
    $(document).ready(function () {
        $("#ad_type").change(function (handler) {
            if(handler.target.value=="native"){
                $(".slot-container").hide();
                $(".native-container").show();
            }else{
                $(".slot-container").show();
                $(".native-container").hide();
            }
        });

        if($("#ad_type").val()=="native"){
            $(".slot-container").hide();
            $(".native-container").show();
        }else{
            $(".slot-container").show();
            $(".native-container").hide();
        }

        $("#my-form").submit(function (e) {
            e.preventDefault();
            var $this = $(this);
            var data = $this.serializeArray();
            var action = $this.attr('action');
            var method = $this.attr('method');
            var current_text = $this.find('button[type=submit]').html();

            var resl={
                slots:[],
                ad_type:"",
                tid:"",
                ip:"",
                ad_count:0,
                publisher:"",
                user_agent:"",
                public_id:""
            };
            for(var i = 1;i<=12;i++){
                var w,h;
                data.forEach(function(d){
                    if(d.name=="ad_type" && d.value){
                        resl.ad_type=d.value;
                    }
                    if(d.name=="tid" && d.value){
                        resl.tid=parseInt(d.value);
                    }
                    if(d.name=="ip" && d.value){
                        resl.ip=d.value;
                    }

                    if(d.name=="ad_count" && d.value){
                        resl.ad_count=parseInt(d.value);
                    }

                    if(d.name=="pack" && d.value){
                        resl.publisher=d.value;
                    }

                    if(d.name=="user_agent" && d.value){
                        resl.user_agent=d.value;
                    }

                    if(d.name=="public_id" && d.value){
                        resl.public_id=d.value;
                    }

                    if (d.name == "width"+i ) {
                        w=d.value
                    }
                    if (d.name == "count"+i ) {
                        h=d.value
                    }
                    if (w&&h){


                        resl.slots.push({
                            size:w,
                            count:parseInt(h)
                        });
                        w,h=null, null;
                    }
                });

            }

            i=resl.publisher;
            typ=resl.ad_type;

            console.log(JSON.stringify(resl));
            //console.log(data);
            $.ajax({
                url : action+"?i="+i+"&type="+typ,
                type : method,
                data : JSON.stringify(resl),
                dataType: 'json',
                beforeSend: function(){
                    $this.find('button[type="submit"]').find('i').attr('class', '').addClass('fa fa-spinner fa-spin');
                },
                complete: function(){
                    if($this.find('button[type="submit"]').find('i').hasClass('fa fa-spinner fa-spin')){
                        $this.find('button[type="submit"]').html(current_text);
                    }
                },
                success: function(data){
                    $(".result").html("");

                    var myKeys=Object.keys(data);
                    myKeys.forEach(function (d) {
                        if(Array.isArray(data[d])){
                            data[d].forEach(function (t) {
                                $(".result").append("<div><img src="+t.ad_img+">"+d+"</div>");
                            });
                        }else{
                            $(".result").append("<img src="+data[d].ad_img+">");
                        }

                    });

                },
                error: function(xhr){
                    //alert("An error occured: " + xhr.status + " " + xhr.statusText);
                }
            });
        });




    });
</script>
</body>
</html>
`
