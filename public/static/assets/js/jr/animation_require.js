define(['jr/core'],function(jr) {
    jr.extend({
        animation: {
            timer: function (call, overCall, start, end, speed) {
                if (!call) return;
                var _start = start != null ? start : 0;
                var _end = end != null ? end : 100;

                var _setup = 0;
                var _interval = 0;

                //速度:1-5倍
                if (speed < 1 || speed > 5) {
                    _setup = (end - start) < 0 ? -speed : speed;
                    _interval = 20;
                } else {
                    //(end-start)/20,假定interval = 50
                    _setup = (_end - _start) / (4 * (6 - speed));
                    //参照100次的时间来计算步长
                    _setup *= (Math.abs(_end - _start) / 100);

                    //修正interval
                    _interval = 1000 / _setup;
                    _interval = _interval < 0 ? -Math.ceil(_interval) : Math.floor(_interval);
                    if (_interval < 30) _interval = 30;
                }
                // alert("start:"+_start+"/ end:"+_end+" / setup:"+_setup+'/ interval: '+_interval)

                var t = setInterval(function () {
                    _start = _start + _setup;
                    if (Math.abs(_start) >= Math.abs(_end)) {
                        _start = _end;
                        if (overCall instanceof Function) overCall();
                        clearInterval(t);
                    }
                    call(_start, _setup);
                }, _interval);
            },
            opacity: function (e, call, end, speed) {
                var ele = jr.$(e);
                var s = jr.style(ele);
                var start = s["opacity"];

                //获取ie的透明度
                if (start == undefined) {
                    if (ele.filters.alpha) {
                        start = ele.filters.alpha.opacity / 100;
                    } else {
                        start = 0;
                    }
                }

                this.timer((function (e) {
                    return function (t, setup) {
                        var isRev = setup < 0;
                        var _opacity = isRev ? (100 + t) / 100 : t / 100;
                        e.style.opacity = _opacity;
                        e.style.filter = 'alpha(opacity=' + (_opacity * 100) + ')';
                        if (t == -100 && isRev) {
                            e.style.display = 'none';
                        } else if (!isRev && t == setup) {
                            e.style.display = '';
                        }
                    }
                })(ele), call, parseInt(start), end, speed);
            },
            fade: function (e, call, speed) {
                var _speed = speed != null ? speed : 3;
                var _end = -100;
                this.opacity(e, call, _end, _speed);
            },
            show: function (e, call, speed) {
                var _speed = speed != null ? speed : 3;
                var _end = 100;
                this.opacity(e, call, _end, _speed);
            },
            toggle: function (e, call, speed) {
                this._toggle(e, 'wh', call, speed);
            },
            toggleWidth: function (e, call, speed) {
                this._toggle(e, 'w', call, speed);
            },
            toggleHeight: function (e, call, speed) {
                this._toggle(e, 'h', call, speed);
            },
            _toggle: function (e, direction, call, speed) {
                e = jr.$(e);
                var style = jr.style(e);
                var w = e.offsetWidth;
                var h = e.offsetHeight;
                var tw = parseInt(e.getAttribute("toggle-w") || 0);
                var th = parseInt(e.getAttribute("toggle-h") || 0);

                //init
                if (tw == 0 || th == 0) {
                    tw = jr.clientWidth(e);
                    th = jr.clientHeight(e);
                    if (w == 0 || h == 0) {
                        w = tw;
                        h = th;
                    }
                    e.setAttribute('toggle-w', w);
                    e.setAttribute('toggle-h', h);
                }

                var speedX = speed == null ? 2 : speed;
                var speedY = speedX * (th / tw);


                if (style["display"] == 'none') {
                    var css = {overflow: 'hidden', display: 'inherit'};
                    if (direction.indexOf('w') != -1) {
                        css.width = '0px';
                    }
                    if (direction.indexOf('h') != -1) {
                        css.height = '0px';
                    }
                    jr.style(e, css);
                    this._toggleShow(e, direction, call, w, tw, h, th, speedX, speedY);
                } else {
                    e.style.overflow = 'hidden';
                    this._toggleClose(e, direction, call, w, tw, h, th, speedX, speedY);
                }
            },
            _toggleShow: function (e, direction, call, w, tw, h, th, speedX, speedY) {
                if (direction.indexOf('w') != -1) {
                    this.timer((function (e, w) {
                        return function (t, setup) {
                            e.style.width = (setup > 0 ? t : (w + t)) + 'px';
                        };
                    })(e, w), call, 0, tw, speedX);
                }
                if (direction.indexOf('h') != -1) {
                    this.timer((function (e, h) {
                        return function (t, setup) {
                            e.style.height = (setup > 0 ? t : (h + t)) + 'px';
                        };
                    })(e, h), call, 0, th, speedY);
                }
            },
            _toggleClose: function (e, direction, call, w, tw, h, tw, speedX, speedY) {
                if (direction.indexOf('w') != -1) {
                    this.timer((function (e, w) {
                        return function (t, setup) {
                            var _w = (setup > 0 ? t : (w + t));
                            if (_w < 0) _w = 0;
                            e.style.width = _w + 'px';
                            if (_w == 0) {
                                e.style.display = 'none';
                            }
                        };
                    })(e, w), call, 0, -w, speedX);
                }
                if (direction.indexOf('h') != -1) {
                    this.timer((function (e, h) {
                        return function (t, setup) {
                            var _h = setup > 0 ? t : h + t;
                            if (_h < 0) _h = 0;
                            e.style.height = _h + 'px';
                            if (_h == 0) {
                                e.style.display = 'none';
                            }
                        };
                    })(e, h), call, 0, -h, speedY);
                }
            }
        }
    });
});


