    /**
     * 获得base64
     * @param {Object} obj
     * @param {Number} [obj.maxWH] 图片需要压缩的宽度，高度会跟随调整
     * @param {Number} [obj.quality=0.8] 压缩质量，不压缩为1
     * @param {Function} [obj.before(this, blob, file)] 处理前函数,this指向的是input:file
     * @param {Function} obj.success(obj) 处理后函数
     * @example
     *
     */
    var resizeImage = function(obj) {
        var URL = window.URL || window.webkitURL;
        var url = URL.createObjectURL(obj.file);
        _create(url);

        /**
         * 生成base64
         * @param blob 通过file获得的二进制
         */
        function _create(url) {
            var img = new Image();
            img.src = url;

            img.onload = function() {
                var that = this;
                $("#t").html(that.width + "--" + that.height);
                //生成比例
                var w = that.width,
                    h = that.height,
                    scale = w / h;

                if (obj.trim == undefined || obj.trim == true) {
                    if (that.width > that.height) {
                        if (obj.maxWH && w > obj.maxWH) {
                            w = obj.maxWH;
                            h = Math.ceil(w / scale);
                        }

                    } else {
                        if (obj.maxWH && h > obj.maxWH) {
                            h = obj.maxWH;
                            w = Math.ceil(h * scale); //如果这里是小数，安卓系统压缩后的图片右侧将会有1px黑边。
                        }
                    }
                }

                //生成canvas
                var canvas = document.createElement('canvas');
                var ctx = canvas.getContext('2d');

                $(canvas).attr({
                    width: w,
                    height: h
                });
                ctx.fillStyle = "#fff";
                ctx.fillRect(0, 0, w, h);

                ctx.drawImage(that, 0, 0, w, h);

                var  imageData = ctx.getImageData(0, 0, w, h);
                var  data = imageData.data;

                for (let i = 0; i < data.length; i += 4) {
                    let r = data[i],
                        g = data[i + 1],
                        b = data[i + 2];

                    if ([r, g, b].every(v => v < 256 && v > 250)) {
                        data[i + 3] = 0;
                    }
                }

                ctx.putImageData(imageData, 0, 0);

                var base64 = canvas.toDataURL(obj.suffix, obj.quality || 0.8);
                var clearBase64 = base64.substr(base64.indexOf(',') + 1);
                var blob = b64toBlob(clearBase64, obj.suffix);

                // 生成结果
                var result = {
                    base64: base64,
                    clearBase64: clearBase64,
                    blob: blob,
                    width: w,
                    height: h
                };

                // 执行后函数
                obj.success(result);
            };
        }

        function b64toBlob(b64Data, contentType, sliceSize) {
            contentType = contentType || '';
            sliceSize = sliceSize || 512;

            var byteCharacters = atob(b64Data);
            var byteArrays = [];

            for (var offset = 0; offset < byteCharacters.length; offset += sliceSize) {
                var slice = byteCharacters.slice(offset, offset + sliceSize);

                var byteNumbers = new Array(slice.length);
                for (var i = 0; i < slice.length; i++) {
                    byteNumbers[i] = slice.charCodeAt(i);
                }

                var byteArray = new Uint8Array(byteNumbers);

                byteArrays.push(byteArray);
            }

            var blob = new Blob(byteArrays, {type: contentType});
            return blob;
        }
    };