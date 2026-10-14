(() => {
  var __create = Object.create;
  var __defProp = Object.defineProperty;
  var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
  var __getOwnPropNames = Object.getOwnPropertyNames;
  var __getProtoOf = Object.getPrototypeOf;
  var __hasOwnProp = Object.prototype.hasOwnProperty;
  var __defNormalProp = (obj, key, value) => key in obj ? __defProp(obj, key, { enumerable: true, configurable: true, writable: true, value }) : obj[key] = value;
  var __esm = (fn, res) => function __init() {
    return fn && (res = (0, fn[__getOwnPropNames(fn)[0]])(fn = 0)), res;
  };
  var __commonJS = (cb, mod) => function __require() {
    return mod || (0, cb[__getOwnPropNames(cb)[0]])((mod = { exports: {} }).exports, mod), mod.exports;
  };
  var __export = (target, all2) => {
    for (var name2 in all2)
      __defProp(target, name2, { get: all2[name2], enumerable: true });
  };
  var __copyProps = (to, from, except, desc) => {
    if (from && typeof from === "object" || typeof from === "function") {
      for (let key of __getOwnPropNames(from))
        if (!__hasOwnProp.call(to, key) && key !== except)
          __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
    }
    return to;
  };
  var __toESM = (mod, isNodeMode, target) => (target = mod != null ? __create(__getProtoOf(mod)) : {}, __copyProps(
    // If the importer is in node compatibility mode or this is not an ESM
    // file that has been converted to a CommonJS file using a Babel-
    // compatible transform (i.e. "__esModule" has not been set), then set
    // "default" to the CommonJS "module.exports" for node compatibility.
    isNodeMode || !mod || !mod.__esModule ? __defProp(target, "default", { value: mod, enumerable: true }) : target,
    mod
  ));
  var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);
  var __publicField = (obj, key, value) => __defNormalProp(obj, typeof key !== "symbol" ? key + "" : key, value);

  // ../../../ui/node_modules/.pnpm/extend@3.0.2/node_modules/extend/index.js
  var require_extend = __commonJS({
    "../../../ui/node_modules/.pnpm/extend@3.0.2/node_modules/extend/index.js"(exports, module) {
      "use strict";
      var hasOwn2 = Object.prototype.hasOwnProperty;
      var toStr = Object.prototype.toString;
      var defineProperty = Object.defineProperty;
      var gOPD = Object.getOwnPropertyDescriptor;
      var isArray2 = function isArray3(arr) {
        if (typeof Array.isArray === "function") {
          return Array.isArray(arr);
        }
        return toStr.call(arr) === "[object Array]";
      };
      var isPlainObject3 = function isPlainObject4(obj) {
        if (!obj || toStr.call(obj) !== "[object Object]") {
          return false;
        }
        var hasOwnConstructor = hasOwn2.call(obj, "constructor");
        var hasIsPrototypeOf = obj.constructor && obj.constructor.prototype && hasOwn2.call(obj.constructor.prototype, "isPrototypeOf");
        if (obj.constructor && !hasOwnConstructor && !hasIsPrototypeOf) {
          return false;
        }
        var key;
        for (key in obj) {
        }
        return typeof key === "undefined" || hasOwn2.call(obj, key);
      };
      var setProperty = function setProperty2(target, options) {
        if (defineProperty && options.name === "__proto__") {
          defineProperty(target, options.name, {
            enumerable: true,
            configurable: true,
            value: options.newValue,
            writable: true
          });
        } else {
          target[options.name] = options.newValue;
        }
      };
      var getProperty = function getProperty2(obj, name2) {
        if (name2 === "__proto__") {
          if (!hasOwn2.call(obj, name2)) {
            return void 0;
          } else if (gOPD) {
            return gOPD(obj, name2).value;
          }
        }
        return obj[name2];
      };
      module.exports = function extend3() {
        var options, name2, src, copy, copyIsArray, clone;
        var target = arguments[0];
        var i = 1;
        var length = arguments.length;
        var deep = false;
        if (typeof target === "boolean") {
          deep = target;
          target = arguments[1] || {};
          i = 2;
        }
        if (target == null || typeof target !== "object" && typeof target !== "function") {
          target = {};
        }
        for (; i < length; ++i) {
          options = arguments[i];
          if (options != null) {
            for (name2 in options) {
              src = getProperty(target, name2);
              copy = getProperty(options, name2);
              if (target !== copy) {
                if (deep && copy && (isPlainObject3(copy) || (copyIsArray = isArray2(copy)))) {
                  if (copyIsArray) {
                    copyIsArray = false;
                    clone = src && isArray2(src) ? src : [];
                  } else {
                    clone = src && isPlainObject3(src) ? src : {};
                  }
                  setProperty(target, { name: name2, newValue: extend3(deep, clone, copy) });
                } else if (typeof copy !== "undefined") {
                  setProperty(target, { name: name2, newValue: copy });
                }
              }
            }
          }
        }
        return target;
      };
    }
  });

  // ../../../ui/node_modules/.pnpm/acorn@8.15.0/node_modules/acorn/dist/acorn.mjs
  var acorn_exports = {};
  __export(acorn_exports, {
    Node: () => Node,
    Parser: () => Parser,
    Position: () => Position,
    SourceLocation: () => SourceLocation,
    TokContext: () => TokContext,
    Token: () => Token,
    TokenType: () => TokenType,
    defaultOptions: () => defaultOptions,
    getLineInfo: () => getLineInfo,
    isIdentifierChar: () => isIdentifierChar,
    isIdentifierStart: () => isIdentifierStart,
    isNewLine: () => isNewLine,
    keywordTypes: () => keywords,
    lineBreak: () => lineBreak,
    lineBreakG: () => lineBreakG,
    nonASCIIwhitespace: () => nonASCIIwhitespace,
    parse: () => parse3,
    parseExpressionAt: () => parseExpressionAt2,
    tokContexts: () => types,
    tokTypes: () => types$1,
    tokenizer: () => tokenizer2,
    version: () => version
  });
  function isInAstralSet(code3, set) {
    var pos = 65536;
    for (var i = 0; i < set.length; i += 2) {
      pos += set[i];
      if (pos > code3) {
        return false;
      }
      pos += set[i + 1];
      if (pos >= code3) {
        return true;
      }
    }
    return false;
  }
  function isIdentifierStart(code3, astral) {
    if (code3 < 65) {
      return code3 === 36;
    }
    if (code3 < 91) {
      return true;
    }
    if (code3 < 97) {
      return code3 === 95;
    }
    if (code3 < 123) {
      return true;
    }
    if (code3 <= 65535) {
      return code3 >= 170 && nonASCIIidentifierStart.test(String.fromCharCode(code3));
    }
    if (astral === false) {
      return false;
    }
    return isInAstralSet(code3, astralIdentifierStartCodes);
  }
  function isIdentifierChar(code3, astral) {
    if (code3 < 48) {
      return code3 === 36;
    }
    if (code3 < 58) {
      return true;
    }
    if (code3 < 65) {
      return false;
    }
    if (code3 < 91) {
      return true;
    }
    if (code3 < 97) {
      return code3 === 95;
    }
    if (code3 < 123) {
      return true;
    }
    if (code3 <= 65535) {
      return code3 >= 170 && nonASCIIidentifier.test(String.fromCharCode(code3));
    }
    if (astral === false) {
      return false;
    }
    return isInAstralSet(code3, astralIdentifierStartCodes) || isInAstralSet(code3, astralIdentifierCodes);
  }
  function binop(name2, prec) {
    return new TokenType(name2, { beforeExpr: true, binop: prec });
  }
  function kw(name2, options) {
    if (options === void 0) options = {};
    options.keyword = name2;
    return keywords[name2] = new TokenType(name2, options);
  }
  function isNewLine(code3) {
    return code3 === 10 || code3 === 13 || code3 === 8232 || code3 === 8233;
  }
  function nextLineBreak(code3, from, end) {
    if (end === void 0) end = code3.length;
    for (var i = from; i < end; i++) {
      var next = code3.charCodeAt(i);
      if (isNewLine(next)) {
        return i < end - 1 && next === 13 && code3.charCodeAt(i + 1) === 10 ? i + 2 : i + 1;
      }
    }
    return -1;
  }
  function wordsRegexp(words) {
    return regexpCache[words] || (regexpCache[words] = new RegExp("^(?:" + words.replace(/ /g, "|") + ")$"));
  }
  function codePointToString(code3) {
    if (code3 <= 65535) {
      return String.fromCharCode(code3);
    }
    code3 -= 65536;
    return String.fromCharCode((code3 >> 10) + 55296, (code3 & 1023) + 56320);
  }
  function getLineInfo(input, offset2) {
    for (var line = 1, cur = 0; ; ) {
      var nextBreak = nextLineBreak(input, cur, offset2);
      if (nextBreak < 0) {
        return new Position(line, offset2 - cur);
      }
      ++line;
      cur = nextBreak;
    }
  }
  function getOptions(opts) {
    var options = {};
    for (var opt in defaultOptions) {
      options[opt] = opts && hasOwn(opts, opt) ? opts[opt] : defaultOptions[opt];
    }
    if (options.ecmaVersion === "latest") {
      options.ecmaVersion = 1e8;
    } else if (options.ecmaVersion == null) {
      if (!warnedAboutEcmaVersion && typeof console === "object" && console.warn) {
        warnedAboutEcmaVersion = true;
        console.warn("Since Acorn 8.0.0, options.ecmaVersion is required.\nDefaulting to 2020, but this will stop working in the future.");
      }
      options.ecmaVersion = 11;
    } else if (options.ecmaVersion >= 2015) {
      options.ecmaVersion -= 2009;
    }
    if (options.allowReserved == null) {
      options.allowReserved = options.ecmaVersion < 5;
    }
    if (!opts || opts.allowHashBang == null) {
      options.allowHashBang = options.ecmaVersion >= 14;
    }
    if (isArray(options.onToken)) {
      var tokens = options.onToken;
      options.onToken = function(token) {
        return tokens.push(token);
      };
    }
    if (isArray(options.onComment)) {
      options.onComment = pushComment(options, options.onComment);
    }
    return options;
  }
  function pushComment(options, array) {
    return function(block, text5, start2, end, startLoc, endLoc) {
      var comment = {
        type: block ? "Block" : "Line",
        value: text5,
        start: start2,
        end
      };
      if (options.locations) {
        comment.loc = new SourceLocation(this, startLoc, endLoc);
      }
      if (options.ranges) {
        comment.range = [start2, end];
      }
      array.push(comment);
    };
  }
  function functionFlags(async, generator) {
    return SCOPE_FUNCTION | (async ? SCOPE_ASYNC : 0) | (generator ? SCOPE_GENERATOR : 0);
  }
  function isPrivateNameConflicted(privateNameMap, element) {
    var name2 = element.key.name;
    var curr = privateNameMap[name2];
    var next = "true";
    if (element.type === "MethodDefinition" && (element.kind === "get" || element.kind === "set")) {
      next = (element.static ? "s" : "i") + element.kind;
    }
    if (curr === "iget" && next === "iset" || curr === "iset" && next === "iget" || curr === "sget" && next === "sset" || curr === "sset" && next === "sget") {
      privateNameMap[name2] = "true";
      return false;
    } else if (!curr) {
      privateNameMap[name2] = next;
      return false;
    } else {
      return true;
    }
  }
  function checkKeyName(node2, name2) {
    var computed = node2.computed;
    var key = node2.key;
    return !computed && (key.type === "Identifier" && key.name === name2 || key.type === "Literal" && key.value === name2);
  }
  function isLocalVariableAccess(node2) {
    return node2.type === "Identifier" || node2.type === "ParenthesizedExpression" && isLocalVariableAccess(node2.expression);
  }
  function isPrivateFieldAccess(node2) {
    return node2.type === "MemberExpression" && node2.property.type === "PrivateIdentifier" || node2.type === "ChainExpression" && isPrivateFieldAccess(node2.expression) || node2.type === "ParenthesizedExpression" && isPrivateFieldAccess(node2.expression);
  }
  function finishNodeAt(node2, type, pos, loc) {
    node2.type = type;
    node2.end = pos;
    if (this.options.locations) {
      node2.loc.end = loc;
    }
    if (this.options.ranges) {
      node2.range[1] = pos;
    }
    return node2;
  }
  function buildUnicodeData(ecmaVersion) {
    var d = data[ecmaVersion] = {
      binary: wordsRegexp(unicodeBinaryProperties[ecmaVersion] + " " + unicodeGeneralCategoryValues),
      binaryOfStrings: wordsRegexp(unicodeBinaryPropertiesOfStrings[ecmaVersion]),
      nonBinary: {
        General_Category: wordsRegexp(unicodeGeneralCategoryValues),
        Script: wordsRegexp(unicodeScriptValues[ecmaVersion])
      }
    };
    d.nonBinary.Script_Extensions = d.nonBinary.Script;
    d.nonBinary.gc = d.nonBinary.General_Category;
    d.nonBinary.sc = d.nonBinary.Script;
    d.nonBinary.scx = d.nonBinary.Script_Extensions;
  }
  function hasProp(obj) {
    for (var _ in obj) {
      return true;
    }
    return false;
  }
  function isRegularExpressionModifier(ch) {
    return ch === 105 || ch === 109 || ch === 115;
  }
  function isSyntaxCharacter(ch) {
    return ch === 36 || ch >= 40 && ch <= 43 || ch === 46 || ch === 63 || ch >= 91 && ch <= 94 || ch >= 123 && ch <= 125;
  }
  function isRegExpIdentifierStart(ch) {
    return isIdentifierStart(ch, true) || ch === 36 || ch === 95;
  }
  function isRegExpIdentifierPart(ch) {
    return isIdentifierChar(ch, true) || ch === 36 || ch === 95 || ch === 8204 || ch === 8205;
  }
  function isControlLetter(ch) {
    return ch >= 65 && ch <= 90 || ch >= 97 && ch <= 122;
  }
  function isValidUnicode(ch) {
    return ch >= 0 && ch <= 1114111;
  }
  function isCharacterClassEscape(ch) {
    return ch === 100 || ch === 68 || ch === 115 || ch === 83 || ch === 119 || ch === 87;
  }
  function isUnicodePropertyNameCharacter(ch) {
    return isControlLetter(ch) || ch === 95;
  }
  function isUnicodePropertyValueCharacter(ch) {
    return isUnicodePropertyNameCharacter(ch) || isDecimalDigit(ch);
  }
  function isClassSetReservedDoublePunctuatorCharacter(ch) {
    return ch === 33 || ch >= 35 && ch <= 38 || ch >= 42 && ch <= 44 || ch === 46 || ch >= 58 && ch <= 64 || ch === 94 || ch === 96 || ch === 126;
  }
  function isClassSetSyntaxCharacter(ch) {
    return ch === 40 || ch === 41 || ch === 45 || ch === 47 || ch >= 91 && ch <= 93 || ch >= 123 && ch <= 125;
  }
  function isClassSetReservedPunctuator(ch) {
    return ch === 33 || ch === 35 || ch === 37 || ch === 38 || ch === 44 || ch === 45 || ch >= 58 && ch <= 62 || ch === 64 || ch === 96 || ch === 126;
  }
  function isDecimalDigit(ch) {
    return ch >= 48 && ch <= 57;
  }
  function isHexDigit(ch) {
    return ch >= 48 && ch <= 57 || ch >= 65 && ch <= 70 || ch >= 97 && ch <= 102;
  }
  function hexToInt(ch) {
    if (ch >= 65 && ch <= 70) {
      return 10 + (ch - 65);
    }
    if (ch >= 97 && ch <= 102) {
      return 10 + (ch - 97);
    }
    return ch - 48;
  }
  function isOctalDigit(ch) {
    return ch >= 48 && ch <= 55;
  }
  function stringToNumber(str, isLegacyOctalNumericLiteral) {
    if (isLegacyOctalNumericLiteral) {
      return parseInt(str, 8);
    }
    return parseFloat(str.replace(/_/g, ""));
  }
  function stringToBigInt(str) {
    if (typeof BigInt !== "function") {
      return null;
    }
    return BigInt(str.replace(/_/g, ""));
  }
  function parse3(input, options) {
    return Parser.parse(input, options);
  }
  function parseExpressionAt2(input, pos, options) {
    return Parser.parseExpressionAt(input, pos, options);
  }
  function tokenizer2(input, options) {
    return Parser.tokenizer(input, options);
  }
  var astralIdentifierCodes, astralIdentifierStartCodes, nonASCIIidentifierChars, nonASCIIidentifierStartChars, reservedWords, ecma5AndLessKeywords, keywords$1, keywordRelationalOperator, nonASCIIidentifierStart, nonASCIIidentifier, TokenType, beforeExpr, startsExpr, keywords, types$1, lineBreak, lineBreakG, nonASCIIwhitespace, skipWhiteSpace, ref, hasOwnProperty, toString, hasOwn, isArray, regexpCache, loneSurrogate, Position, SourceLocation, defaultOptions, warnedAboutEcmaVersion, SCOPE_TOP, SCOPE_FUNCTION, SCOPE_ASYNC, SCOPE_GENERATOR, SCOPE_ARROW, SCOPE_SIMPLE_CATCH, SCOPE_SUPER, SCOPE_DIRECT_SUPER, SCOPE_CLASS_STATIC_BLOCK, SCOPE_CLASS_FIELD_INIT, SCOPE_VAR, BIND_NONE, BIND_VAR, BIND_LEXICAL, BIND_FUNCTION, BIND_SIMPLE_CATCH, BIND_OUTSIDE, Parser, prototypeAccessors, pp$9, literal, DestructuringErrors, pp$8, loopLabel, switchLabel, empty$1, FUNC_STATEMENT, FUNC_HANGING_STATEMENT, FUNC_NULLABLE_ID, pp$7, TokContext, types, pp$6, pp$5, empty, pp$4, pp$3, Scope, Node, pp$2, scriptValuesAddedInUnicode, ecma9BinaryProperties, ecma10BinaryProperties, ecma11BinaryProperties, ecma12BinaryProperties, ecma13BinaryProperties, ecma14BinaryProperties, unicodeBinaryProperties, ecma14BinaryPropertiesOfStrings, unicodeBinaryPropertiesOfStrings, unicodeGeneralCategoryValues, ecma9ScriptValues, ecma10ScriptValues, ecma11ScriptValues, ecma12ScriptValues, ecma13ScriptValues, ecma14ScriptValues, unicodeScriptValues, data, ecmaVersion, i, list4, pp$1, BranchID, RegExpValidationState, CharSetNone, CharSetOk, CharSetString, Token, pp, INVALID_TEMPLATE_ESCAPE_ERROR, version;
  var init_acorn = __esm({
    "../../../ui/node_modules/.pnpm/acorn@8.15.0/node_modules/acorn/dist/acorn.mjs"() {
      astralIdentifierCodes = [509, 0, 227, 0, 150, 4, 294, 9, 1368, 2, 2, 1, 6, 3, 41, 2, 5, 0, 166, 1, 574, 3, 9, 9, 7, 9, 32, 4, 318, 1, 80, 3, 71, 10, 50, 3, 123, 2, 54, 14, 32, 10, 3, 1, 11, 3, 46, 10, 8, 0, 46, 9, 7, 2, 37, 13, 2, 9, 6, 1, 45, 0, 13, 2, 49, 13, 9, 3, 2, 11, 83, 11, 7, 0, 3, 0, 158, 11, 6, 9, 7, 3, 56, 1, 2, 6, 3, 1, 3, 2, 10, 0, 11, 1, 3, 6, 4, 4, 68, 8, 2, 0, 3, 0, 2, 3, 2, 4, 2, 0, 15, 1, 83, 17, 10, 9, 5, 0, 82, 19, 13, 9, 214, 6, 3, 8, 28, 1, 83, 16, 16, 9, 82, 12, 9, 9, 7, 19, 58, 14, 5, 9, 243, 14, 166, 9, 71, 5, 2, 1, 3, 3, 2, 0, 2, 1, 13, 9, 120, 6, 3, 6, 4, 0, 29, 9, 41, 6, 2, 3, 9, 0, 10, 10, 47, 15, 343, 9, 54, 7, 2, 7, 17, 9, 57, 21, 2, 13, 123, 5, 4, 0, 2, 1, 2, 6, 2, 0, 9, 9, 49, 4, 2, 1, 2, 4, 9, 9, 330, 3, 10, 1, 2, 0, 49, 6, 4, 4, 14, 10, 5350, 0, 7, 14, 11465, 27, 2343, 9, 87, 9, 39, 4, 60, 6, 26, 9, 535, 9, 470, 0, 2, 54, 8, 3, 82, 0, 12, 1, 19628, 1, 4178, 9, 519, 45, 3, 22, 543, 4, 4, 5, 9, 7, 3, 6, 31, 3, 149, 2, 1418, 49, 513, 54, 5, 49, 9, 0, 15, 0, 23, 4, 2, 14, 1361, 6, 2, 16, 3, 6, 2, 1, 2, 4, 101, 0, 161, 6, 10, 9, 357, 0, 62, 13, 499, 13, 245, 1, 2, 9, 726, 6, 110, 6, 6, 9, 4759, 9, 787719, 239];
      astralIdentifierStartCodes = [0, 11, 2, 25, 2, 18, 2, 1, 2, 14, 3, 13, 35, 122, 70, 52, 268, 28, 4, 48, 48, 31, 14, 29, 6, 37, 11, 29, 3, 35, 5, 7, 2, 4, 43, 157, 19, 35, 5, 35, 5, 39, 9, 51, 13, 10, 2, 14, 2, 6, 2, 1, 2, 10, 2, 14, 2, 6, 2, 1, 4, 51, 13, 310, 10, 21, 11, 7, 25, 5, 2, 41, 2, 8, 70, 5, 3, 0, 2, 43, 2, 1, 4, 0, 3, 22, 11, 22, 10, 30, 66, 18, 2, 1, 11, 21, 11, 25, 71, 55, 7, 1, 65, 0, 16, 3, 2, 2, 2, 28, 43, 28, 4, 28, 36, 7, 2, 27, 28, 53, 11, 21, 11, 18, 14, 17, 111, 72, 56, 50, 14, 50, 14, 35, 39, 27, 10, 22, 251, 41, 7, 1, 17, 2, 60, 28, 11, 0, 9, 21, 43, 17, 47, 20, 28, 22, 13, 52, 58, 1, 3, 0, 14, 44, 33, 24, 27, 35, 30, 0, 3, 0, 9, 34, 4, 0, 13, 47, 15, 3, 22, 0, 2, 0, 36, 17, 2, 24, 20, 1, 64, 6, 2, 0, 2, 3, 2, 14, 2, 9, 8, 46, 39, 7, 3, 1, 3, 21, 2, 6, 2, 1, 2, 4, 4, 0, 19, 0, 13, 4, 31, 9, 2, 0, 3, 0, 2, 37, 2, 0, 26, 0, 2, 0, 45, 52, 19, 3, 21, 2, 31, 47, 21, 1, 2, 0, 185, 46, 42, 3, 37, 47, 21, 0, 60, 42, 14, 0, 72, 26, 38, 6, 186, 43, 117, 63, 32, 7, 3, 0, 3, 7, 2, 1, 2, 23, 16, 0, 2, 0, 95, 7, 3, 38, 17, 0, 2, 0, 29, 0, 11, 39, 8, 0, 22, 0, 12, 45, 20, 0, 19, 72, 200, 32, 32, 8, 2, 36, 18, 0, 50, 29, 113, 6, 2, 1, 2, 37, 22, 0, 26, 5, 2, 1, 2, 31, 15, 0, 328, 18, 16, 0, 2, 12, 2, 33, 125, 0, 80, 921, 103, 110, 18, 195, 2637, 96, 16, 1071, 18, 5, 26, 3994, 6, 582, 6842, 29, 1763, 568, 8, 30, 18, 78, 18, 29, 19, 47, 17, 3, 32, 20, 6, 18, 433, 44, 212, 63, 129, 74, 6, 0, 67, 12, 65, 1, 2, 0, 29, 6135, 9, 1237, 42, 9, 8936, 3, 2, 6, 2, 1, 2, 290, 16, 0, 30, 2, 3, 0, 15, 3, 9, 395, 2309, 106, 6, 12, 4, 8, 8, 9, 5991, 84, 2, 70, 2, 1, 3, 0, 3, 1, 3, 3, 2, 11, 2, 0, 2, 6, 2, 64, 2, 3, 3, 7, 2, 6, 2, 27, 2, 3, 2, 4, 2, 0, 4, 6, 2, 339, 3, 24, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 30, 2, 24, 2, 7, 1845, 30, 7, 5, 262, 61, 147, 44, 11, 6, 17, 0, 322, 29, 19, 43, 485, 27, 229, 29, 3, 0, 496, 6, 2, 3, 2, 1, 2, 14, 2, 196, 60, 67, 8, 0, 1205, 3, 2, 26, 2, 1, 2, 0, 3, 0, 2, 9, 2, 3, 2, 0, 2, 0, 7, 0, 5, 0, 2, 0, 2, 0, 2, 2, 2, 1, 2, 0, 3, 0, 2, 0, 2, 0, 2, 0, 2, 0, 2, 1, 2, 0, 3, 3, 2, 6, 2, 3, 2, 3, 2, 0, 2, 9, 2, 16, 6, 2, 2, 4, 2, 16, 4421, 42719, 33, 4153, 7, 221, 3, 5761, 15, 7472, 16, 621, 2467, 541, 1507, 4938, 6, 4191];
      nonASCIIidentifierChars = "‌‍·̀-ͯ·҃-֑҇-ׇֽֿׁׂׅׄؐ-ًؚ-٩ٰۖ-ۜ۟-۪ۤۧۨ-ۭ۰-۹ܑܰ-݊ަ-ް߀-߉߫-߽߳ࠖ-࠙ࠛ-ࠣࠥ-ࠧࠩ-࡙࠭-࡛ࢗ-࢟࣊-ࣣ࣡-ःऺ-़ा-ॏ॑-ॗॢॣ०-९ঁ-ঃ়া-ৄেৈো-্ৗৢৣ০-৯৾ਁ-ਃ਼ਾ-ੂੇੈੋ-੍ੑ੦-ੱੵઁ-ઃ઼ા-ૅે-ૉો-્ૢૣ૦-૯ૺ-૿ଁ-ଃ଼ା-ୄେୈୋ-୍୕-ୗୢୣ୦-୯ஂா-ூெ-ைொ-்ௗ௦-௯ఀ-ఄ఼ా-ౄె-ైొ-్ౕౖౢౣ౦-౯ಁ-ಃ಼ಾ-ೄೆ-ೈೊ-್ೕೖೢೣ೦-೯ೳഀ-ഃ഻഼ാ-ൄെ-ൈൊ-്ൗൢൣ൦-൯ඁ-ඃ්ා-ුූෘ-ෟ෦-෯ෲෳัิ-ฺ็-๎๐-๙ັິ-ຼ່-໎໐-໙༘༙༠-༩༹༵༷༾༿ཱ-྄྆྇ྍ-ྗྙ-ྼ࿆ါ-ှ၀-၉ၖ-ၙၞ-ၠၢ-ၤၧ-ၭၱ-ၴႂ-ႍႏ-ႝ፝-፟፩-፱ᜒ-᜕ᜲ-᜴ᝒᝓᝲᝳ឴-៓៝០-៩᠋-᠍᠏-᠙ᢩᤠ-ᤫᤰ-᤻᥆-᥏᧐-᧚ᨗ-ᨛᩕ-ᩞ᩠-᩿᩼-᪉᪐-᪙᪰-᪽ᪿ-ᫎᬀ-ᬄ᬴-᭄᭐-᭙᭫-᭳ᮀ-ᮂᮡ-ᮭ᮰-᮹᯦-᯳ᰤ-᰷᱀-᱉᱐-᱙᳐-᳔᳒-᳨᳭᳴᳷-᳹᷀-᷿‌‍‿⁀⁔⃐-⃥⃜⃡-⃰⳯-⵿⳱ⷠ-〪ⷿ-゙゚〯・꘠-꘩꙯ꙴ-꙽ꚞꚟ꛰꛱ꠂ꠆ꠋꠣ-ꠧ꠬ꢀꢁꢴ-ꣅ꣐-꣙꣠-꣱ꣿ-꤉ꤦ-꤭ꥇ-꥓ꦀ-ꦃ꦳-꧀꧐-꧙ꧥ꧰-꧹ꨩ-ꨶꩃꩌꩍ꩐-꩙ꩻ-ꩽꪰꪲ-ꪴꪷꪸꪾ꪿꫁ꫫ-ꫯꫵ꫶ꯣ-ꯪ꯬꯭꯰-꯹ﬞ︀-️︠-︯︳︴﹍-﹏０-９＿･";
      nonASCIIidentifierStartChars = "ªµºÀ-ÖØ-öø-ˁˆ-ˑˠ-ˤˬˮͰ-ʹͶͷͺ-ͽͿΆΈ-ΊΌΎ-ΡΣ-ϵϷ-ҁҊ-ԯԱ-Ֆՙՠ-ֈא-תׯ-ײؠ-يٮٯٱ-ۓەۥۦۮۯۺ-ۼۿܐܒ-ܯݍ-ޥޱߊ-ߪߴߵߺࠀ-ࠕࠚࠤࠨࡀ-ࡘࡠ-ࡪࡰ-ࢇࢉ-ࢎࢠ-ࣉऄ-हऽॐक़-ॡॱ-ঀঅ-ঌএঐও-নপ-রলশ-হঽৎড়ঢ়য়-ৡৰৱৼਅ-ਊਏਐਓ-ਨਪ-ਰਲਲ਼ਵਸ਼ਸਹਖ਼-ੜਫ਼ੲ-ੴઅ-ઍએ-ઑઓ-નપ-રલળવ-હઽૐૠૡૹଅ-ଌଏଐଓ-ନପ-ରଲଳଵ-ହଽଡ଼ଢ଼ୟ-ୡୱஃஅ-ஊஎ-ஐஒ-கஙசஜஞடணதந-பம-ஹௐఅ-ఌఎ-ఐఒ-నప-హఽౘ-ౚౝౠౡಀಅ-ಌಎ-ಐಒ-ನಪ-ಳವ-ಹಽೝೞೠೡೱೲഄ-ഌഎ-ഐഒ-ഺഽൎൔ-ൖൟ-ൡൺ-ൿඅ-ඖක-නඳ-රලව-ෆก-ะาำเ-ๆກຂຄຆ-ຊຌ-ຣລວ-ະາຳຽເ-ໄໆໜ-ໟༀཀ-ཇཉ-ཬྈ-ྌက-ဪဿၐ-ၕၚ-ၝၡၥၦၮ-ၰၵ-ႁႎႠ-ჅჇჍა-ჺჼ-ቈቊ-ቍቐ-ቖቘቚ-ቝበ-ኈኊ-ኍነ-ኰኲ-ኵኸ-ኾዀዂ-ዅወ-ዖዘ-ጐጒ-ጕጘ-ፚᎀ-ᎏᎠ-Ᏽᏸ-ᏽᐁ-ᙬᙯ-ᙿᚁ-ᚚᚠ-ᛪᛮ-ᛸᜀ-ᜑᜟ-ᜱᝀ-ᝑᝠ-ᝬᝮ-ᝰក-ឳៗៜᠠ-ᡸᢀ-ᢨᢪᢰ-ᣵᤀ-ᤞᥐ-ᥭᥰ-ᥴᦀ-ᦫᦰ-ᧉᨀ-ᨖᨠ-ᩔᪧᬅ-ᬳᭅ-ᭌᮃ-ᮠᮮᮯᮺ-ᯥᰀ-ᰣᱍ-ᱏᱚ-ᱽᲀ-ᲊᲐ-ᲺᲽ-Ჿᳩ-ᳬᳮ-ᳳᳵᳶᳺᴀ-ᶿḀ-ἕἘ-Ἕἠ-ὅὈ-Ὅὐ-ὗὙὛὝὟ-ώᾀ-ᾴᾶ-ᾼιῂ-ῄῆ-ῌῐ-ΐῖ-Ίῠ-Ῥῲ-ῴῶ-ῼⁱⁿₐ-ₜℂℇℊ-ℓℕ℘-ℝℤΩℨK-ℹℼ-ℿⅅ-ⅉⅎⅠ-ↈⰀ-ⳤⳫ-ⳮⳲⳳⴀ-ⴥⴧⴭⴰ-ⵧⵯⶀ-ⶖⶠ-ⶦⶨ-ⶮⶰ-ⶶⶸ-ⶾⷀ-ⷆⷈ-ⷎⷐ-ⷖⷘ-ⷞ々-〇〡-〩〱-〵〸-〼ぁ-ゖ゛-ゟァ-ヺー-ヿㄅ-ㄯㄱ-ㆎㆠ-ㆿㇰ-ㇿ㐀-䶿一-ꒌꓐ-ꓽꔀ-ꘌꘐ-ꘟꘪꘫꙀ-ꙮꙿ-ꚝꚠ-ꛯꜗ-ꜟꜢ-ꞈꞋ-ꟍꟐꟑꟓꟕ-Ƛꟲ-ꠁꠃ-ꠅꠇ-ꠊꠌ-ꠢꡀ-ꡳꢂ-ꢳꣲ-ꣷꣻꣽꣾꤊ-ꤥꤰ-ꥆꥠ-ꥼꦄ-ꦲꧏꧠ-ꧤꧦ-ꧯꧺ-ꧾꨀ-ꨨꩀ-ꩂꩄ-ꩋꩠ-ꩶꩺꩾ-ꪯꪱꪵꪶꪹ-ꪽꫀꫂꫛ-ꫝꫠ-ꫪꫲ-ꫴꬁ-ꬆꬉ-ꬎꬑ-ꬖꬠ-ꬦꬨ-ꬮꬰ-ꭚꭜ-ꭩꭰ-ꯢ가-힣ힰ-ퟆퟋ-ퟻ豈-舘並-龎ﬀ-ﬆﬓ-ﬗיִײַ-ﬨשׁ-זּטּ-לּמּנּסּףּפּצּ-ﮱﯓ-ﴽﵐ-ﶏﶒ-ﷇﷰ-ﷻﹰ-ﹴﹶ-ﻼＡ-Ｚａ-ｚｦ-ﾾￂ-ￇￊ-ￏￒ-ￗￚ-ￜ";
      reservedWords = {
        3: "abstract boolean byte char class double enum export extends final float goto implements import int interface long native package private protected public short static super synchronized throws transient volatile",
        5: "class enum extends super const export import",
        6: "enum",
        strict: "implements interface let package private protected public static yield",
        strictBind: "eval arguments"
      };
      ecma5AndLessKeywords = "break case catch continue debugger default do else finally for function if return switch throw try var while with null true false instanceof typeof void delete new in this";
      keywords$1 = {
        5: ecma5AndLessKeywords,
        "5module": ecma5AndLessKeywords + " export import",
        6: ecma5AndLessKeywords + " const class extends export import super"
      };
      keywordRelationalOperator = /^in(stanceof)?$/;
      nonASCIIidentifierStart = new RegExp("[" + nonASCIIidentifierStartChars + "]");
      nonASCIIidentifier = new RegExp("[" + nonASCIIidentifierStartChars + nonASCIIidentifierChars + "]");
      TokenType = function TokenType2(label, conf) {
        if (conf === void 0) conf = {};
        this.label = label;
        this.keyword = conf.keyword;
        this.beforeExpr = !!conf.beforeExpr;
        this.startsExpr = !!conf.startsExpr;
        this.isLoop = !!conf.isLoop;
        this.isAssign = !!conf.isAssign;
        this.prefix = !!conf.prefix;
        this.postfix = !!conf.postfix;
        this.binop = conf.binop || null;
        this.updateContext = null;
      };
      beforeExpr = { beforeExpr: true };
      startsExpr = { startsExpr: true };
      keywords = {};
      types$1 = {
        num: new TokenType("num", startsExpr),
        regexp: new TokenType("regexp", startsExpr),
        string: new TokenType("string", startsExpr),
        name: new TokenType("name", startsExpr),
        privateId: new TokenType("privateId", startsExpr),
        eof: new TokenType("eof"),
        // Punctuation token types.
        bracketL: new TokenType("[", { beforeExpr: true, startsExpr: true }),
        bracketR: new TokenType("]"),
        braceL: new TokenType("{", { beforeExpr: true, startsExpr: true }),
        braceR: new TokenType("}"),
        parenL: new TokenType("(", { beforeExpr: true, startsExpr: true }),
        parenR: new TokenType(")"),
        comma: new TokenType(",", beforeExpr),
        semi: new TokenType(";", beforeExpr),
        colon: new TokenType(":", beforeExpr),
        dot: new TokenType("."),
        question: new TokenType("?", beforeExpr),
        questionDot: new TokenType("?."),
        arrow: new TokenType("=>", beforeExpr),
        template: new TokenType("template"),
        invalidTemplate: new TokenType("invalidTemplate"),
        ellipsis: new TokenType("...", beforeExpr),
        backQuote: new TokenType("`", startsExpr),
        dollarBraceL: new TokenType("${", { beforeExpr: true, startsExpr: true }),
        // Operators. These carry several kinds of properties to help the
        // parser use them properly (the presence of these properties is
        // what categorizes them as operators).
        //
        // `binop`, when present, specifies that this operator is a binary
        // operator, and will refer to its precedence.
        //
        // `prefix` and `postfix` mark the operator as a prefix or postfix
        // unary operator.
        //
        // `isAssign` marks all of `=`, `+=`, `-=` etcetera, which act as
        // binary operators with a very low precedence, that should result
        // in AssignmentExpression nodes.
        eq: new TokenType("=", { beforeExpr: true, isAssign: true }),
        assign: new TokenType("_=", { beforeExpr: true, isAssign: true }),
        incDec: new TokenType("++/--", { prefix: true, postfix: true, startsExpr: true }),
        prefix: new TokenType("!/~", { beforeExpr: true, prefix: true, startsExpr: true }),
        logicalOR: binop("||", 1),
        logicalAND: binop("&&", 2),
        bitwiseOR: binop("|", 3),
        bitwiseXOR: binop("^", 4),
        bitwiseAND: binop("&", 5),
        equality: binop("==/!=/===/!==", 6),
        relational: binop("</>/<=/>=", 7),
        bitShift: binop("<</>>/>>>", 8),
        plusMin: new TokenType("+/-", { beforeExpr: true, binop: 9, prefix: true, startsExpr: true }),
        modulo: binop("%", 10),
        star: binop("*", 10),
        slash: binop("/", 10),
        starstar: new TokenType("**", { beforeExpr: true }),
        coalesce: binop("??", 1),
        // Keyword token types.
        _break: kw("break"),
        _case: kw("case", beforeExpr),
        _catch: kw("catch"),
        _continue: kw("continue"),
        _debugger: kw("debugger"),
        _default: kw("default", beforeExpr),
        _do: kw("do", { isLoop: true, beforeExpr: true }),
        _else: kw("else", beforeExpr),
        _finally: kw("finally"),
        _for: kw("for", { isLoop: true }),
        _function: kw("function", startsExpr),
        _if: kw("if"),
        _return: kw("return", beforeExpr),
        _switch: kw("switch"),
        _throw: kw("throw", beforeExpr),
        _try: kw("try"),
        _var: kw("var"),
        _const: kw("const"),
        _while: kw("while", { isLoop: true }),
        _with: kw("with"),
        _new: kw("new", { beforeExpr: true, startsExpr: true }),
        _this: kw("this", startsExpr),
        _super: kw("super", startsExpr),
        _class: kw("class", startsExpr),
        _extends: kw("extends", beforeExpr),
        _export: kw("export"),
        _import: kw("import", startsExpr),
        _null: kw("null", startsExpr),
        _true: kw("true", startsExpr),
        _false: kw("false", startsExpr),
        _in: kw("in", { beforeExpr: true, binop: 7 }),
        _instanceof: kw("instanceof", { beforeExpr: true, binop: 7 }),
        _typeof: kw("typeof", { beforeExpr: true, prefix: true, startsExpr: true }),
        _void: kw("void", { beforeExpr: true, prefix: true, startsExpr: true }),
        _delete: kw("delete", { beforeExpr: true, prefix: true, startsExpr: true })
      };
      lineBreak = /\r\n?|\n|\u2028|\u2029/;
      lineBreakG = new RegExp(lineBreak.source, "g");
      nonASCIIwhitespace = /[\u1680\u2000-\u200a\u202f\u205f\u3000\ufeff]/;
      skipWhiteSpace = /(?:\s|\/\/.*|\/\*[^]*?\*\/)*/g;
      ref = Object.prototype;
      hasOwnProperty = ref.hasOwnProperty;
      toString = ref.toString;
      hasOwn = Object.hasOwn || (function(obj, propName) {
        return hasOwnProperty.call(obj, propName);
      });
      isArray = Array.isArray || (function(obj) {
        return toString.call(obj) === "[object Array]";
      });
      regexpCache = /* @__PURE__ */ Object.create(null);
      loneSurrogate = /(?:[\uD800-\uDBFF](?![\uDC00-\uDFFF])|(?:[^\uD800-\uDBFF]|^)[\uDC00-\uDFFF])/;
      Position = function Position2(line, col) {
        this.line = line;
        this.column = col;
      };
      Position.prototype.offset = function offset(n) {
        return new Position(this.line, this.column + n);
      };
      SourceLocation = function SourceLocation2(p, start2, end) {
        this.start = start2;
        this.end = end;
        if (p.sourceFile !== null) {
          this.source = p.sourceFile;
        }
      };
      defaultOptions = {
        // `ecmaVersion` indicates the ECMAScript version to parse. Must be
        // either 3, 5, 6 (or 2015), 7 (2016), 8 (2017), 9 (2018), 10
        // (2019), 11 (2020), 12 (2021), 13 (2022), 14 (2023), or `"latest"`
        // (the latest version the library supports). This influences
        // support for strict mode, the set of reserved words, and support
        // for new syntax features.
        ecmaVersion: null,
        // `sourceType` indicates the mode the code should be parsed in.
        // Can be either `"script"` or `"module"`. This influences global
        // strict mode and parsing of `import` and `export` declarations.
        sourceType: "script",
        // `onInsertedSemicolon` can be a callback that will be called when
        // a semicolon is automatically inserted. It will be passed the
        // position of the inserted semicolon as an offset, and if
        // `locations` is enabled, it is given the location as a `{line,
        // column}` object as second argument.
        onInsertedSemicolon: null,
        // `onTrailingComma` is similar to `onInsertedSemicolon`, but for
        // trailing commas.
        onTrailingComma: null,
        // By default, reserved words are only enforced if ecmaVersion >= 5.
        // Set `allowReserved` to a boolean value to explicitly turn this on
        // an off. When this option has the value "never", reserved words
        // and keywords can also not be used as property names.
        allowReserved: null,
        // When enabled, a return at the top level is not considered an
        // error.
        allowReturnOutsideFunction: false,
        // When enabled, import/export statements are not constrained to
        // appearing at the top of the program, and an import.meta expression
        // in a script isn't considered an error.
        allowImportExportEverywhere: false,
        // By default, await identifiers are allowed to appear at the top-level scope only if ecmaVersion >= 2022.
        // When enabled, await identifiers are allowed to appear at the top-level scope,
        // but they are still not allowed in non-async functions.
        allowAwaitOutsideFunction: null,
        // When enabled, super identifiers are not constrained to
        // appearing in methods and do not raise an error when they appear elsewhere.
        allowSuperOutsideMethod: null,
        // When enabled, hashbang directive in the beginning of file is
        // allowed and treated as a line comment. Enabled by default when
        // `ecmaVersion` >= 2023.
        allowHashBang: false,
        // By default, the parser will verify that private properties are
        // only used in places where they are valid and have been declared.
        // Set this to false to turn such checks off.
        checkPrivateFields: true,
        // When `locations` is on, `loc` properties holding objects with
        // `start` and `end` properties in `{line, column}` form (with
        // line being 1-based and column 0-based) will be attached to the
        // nodes.
        locations: false,
        // A function can be passed as `onToken` option, which will
        // cause Acorn to call that function with object in the same
        // format as tokens returned from `tokenizer().getToken()`. Note
        // that you are not allowed to call the parser from the
        // callback—that will corrupt its internal state.
        onToken: null,
        // A function can be passed as `onComment` option, which will
        // cause Acorn to call that function with `(block, text, start,
        // end)` parameters whenever a comment is skipped. `block` is a
        // boolean indicating whether this is a block (`/* */`) comment,
        // `text` is the content of the comment, and `start` and `end` are
        // character offsets that denote the start and end of the comment.
        // When the `locations` option is on, two more parameters are
        // passed, the full `{line, column}` locations of the start and
        // end of the comments. Note that you are not allowed to call the
        // parser from the callback—that will corrupt its internal state.
        // When this option has an array as value, objects representing the
        // comments are pushed to it.
        onComment: null,
        // Nodes have their start and end characters offsets recorded in
        // `start` and `end` properties (directly on the node, rather than
        // the `loc` object, which holds line/column data. To also add a
        // [semi-standardized][range] `range` property holding a `[start,
        // end]` array with the same numbers, set the `ranges` option to
        // `true`.
        //
        // [range]: https://bugzilla.mozilla.org/show_bug.cgi?id=745678
        ranges: false,
        // It is possible to parse multiple files into a single AST by
        // passing the tree produced by parsing the first file as
        // `program` option in subsequent parses. This will add the
        // toplevel forms of the parsed file to the `Program` (top) node
        // of an existing parse tree.
        program: null,
        // When `locations` is on, you can pass this to record the source
        // file in every node's `loc` object.
        sourceFile: null,
        // This value, if given, is stored in every node, whether
        // `locations` is on or off.
        directSourceFile: null,
        // When enabled, parenthesized expressions are represented by
        // (non-standard) ParenthesizedExpression nodes
        preserveParens: false
      };
      warnedAboutEcmaVersion = false;
      SCOPE_TOP = 1;
      SCOPE_FUNCTION = 2;
      SCOPE_ASYNC = 4;
      SCOPE_GENERATOR = 8;
      SCOPE_ARROW = 16;
      SCOPE_SIMPLE_CATCH = 32;
      SCOPE_SUPER = 64;
      SCOPE_DIRECT_SUPER = 128;
      SCOPE_CLASS_STATIC_BLOCK = 256;
      SCOPE_CLASS_FIELD_INIT = 512;
      SCOPE_VAR = SCOPE_TOP | SCOPE_FUNCTION | SCOPE_CLASS_STATIC_BLOCK;
      BIND_NONE = 0;
      BIND_VAR = 1;
      BIND_LEXICAL = 2;
      BIND_FUNCTION = 3;
      BIND_SIMPLE_CATCH = 4;
      BIND_OUTSIDE = 5;
      Parser = function Parser2(options, input, startPos) {
        this.options = options = getOptions(options);
        this.sourceFile = options.sourceFile;
        this.keywords = wordsRegexp(keywords$1[options.ecmaVersion >= 6 ? 6 : options.sourceType === "module" ? "5module" : 5]);
        var reserved = "";
        if (options.allowReserved !== true) {
          reserved = reservedWords[options.ecmaVersion >= 6 ? 6 : options.ecmaVersion === 5 ? 5 : 3];
          if (options.sourceType === "module") {
            reserved += " await";
          }
        }
        this.reservedWords = wordsRegexp(reserved);
        var reservedStrict = (reserved ? reserved + " " : "") + reservedWords.strict;
        this.reservedWordsStrict = wordsRegexp(reservedStrict);
        this.reservedWordsStrictBind = wordsRegexp(reservedStrict + " " + reservedWords.strictBind);
        this.input = String(input);
        this.containsEsc = false;
        if (startPos) {
          this.pos = startPos;
          this.lineStart = this.input.lastIndexOf("\n", startPos - 1) + 1;
          this.curLine = this.input.slice(0, this.lineStart).split(lineBreak).length;
        } else {
          this.pos = this.lineStart = 0;
          this.curLine = 1;
        }
        this.type = types$1.eof;
        this.value = null;
        this.start = this.end = this.pos;
        this.startLoc = this.endLoc = this.curPosition();
        this.lastTokEndLoc = this.lastTokStartLoc = null;
        this.lastTokStart = this.lastTokEnd = this.pos;
        this.context = this.initialContext();
        this.exprAllowed = true;
        this.inModule = options.sourceType === "module";
        this.strict = this.inModule || this.strictDirective(this.pos);
        this.potentialArrowAt = -1;
        this.potentialArrowInForAwait = false;
        this.yieldPos = this.awaitPos = this.awaitIdentPos = 0;
        this.labels = [];
        this.undefinedExports = /* @__PURE__ */ Object.create(null);
        if (this.pos === 0 && options.allowHashBang && this.input.slice(0, 2) === "#!") {
          this.skipLineComment(2);
        }
        this.scopeStack = [];
        this.enterScope(SCOPE_TOP);
        this.regexpState = null;
        this.privateNameStack = [];
      };
      prototypeAccessors = { inFunction: { configurable: true }, inGenerator: { configurable: true }, inAsync: { configurable: true }, canAwait: { configurable: true }, allowSuper: { configurable: true }, allowDirectSuper: { configurable: true }, treatFunctionsAsVar: { configurable: true }, allowNewDotTarget: { configurable: true }, inClassStaticBlock: { configurable: true } };
      Parser.prototype.parse = function parse() {
        var node2 = this.options.program || this.startNode();
        this.nextToken();
        return this.parseTopLevel(node2);
      };
      prototypeAccessors.inFunction.get = function() {
        return (this.currentVarScope().flags & SCOPE_FUNCTION) > 0;
      };
      prototypeAccessors.inGenerator.get = function() {
        return (this.currentVarScope().flags & SCOPE_GENERATOR) > 0;
      };
      prototypeAccessors.inAsync.get = function() {
        return (this.currentVarScope().flags & SCOPE_ASYNC) > 0;
      };
      prototypeAccessors.canAwait.get = function() {
        for (var i = this.scopeStack.length - 1; i >= 0; i--) {
          var ref2 = this.scopeStack[i];
          var flags = ref2.flags;
          if (flags & (SCOPE_CLASS_STATIC_BLOCK | SCOPE_CLASS_FIELD_INIT)) {
            return false;
          }
          if (flags & SCOPE_FUNCTION) {
            return (flags & SCOPE_ASYNC) > 0;
          }
        }
        return this.inModule && this.options.ecmaVersion >= 13 || this.options.allowAwaitOutsideFunction;
      };
      prototypeAccessors.allowSuper.get = function() {
        var ref2 = this.currentThisScope();
        var flags = ref2.flags;
        return (flags & SCOPE_SUPER) > 0 || this.options.allowSuperOutsideMethod;
      };
      prototypeAccessors.allowDirectSuper.get = function() {
        return (this.currentThisScope().flags & SCOPE_DIRECT_SUPER) > 0;
      };
      prototypeAccessors.treatFunctionsAsVar.get = function() {
        return this.treatFunctionsAsVarInScope(this.currentScope());
      };
      prototypeAccessors.allowNewDotTarget.get = function() {
        for (var i = this.scopeStack.length - 1; i >= 0; i--) {
          var ref2 = this.scopeStack[i];
          var flags = ref2.flags;
          if (flags & (SCOPE_CLASS_STATIC_BLOCK | SCOPE_CLASS_FIELD_INIT) || flags & SCOPE_FUNCTION && !(flags & SCOPE_ARROW)) {
            return true;
          }
        }
        return false;
      };
      prototypeAccessors.inClassStaticBlock.get = function() {
        return (this.currentVarScope().flags & SCOPE_CLASS_STATIC_BLOCK) > 0;
      };
      Parser.extend = function extend2() {
        var plugins = [], len = arguments.length;
        while (len--) plugins[len] = arguments[len];
        var cls = this;
        for (var i = 0; i < plugins.length; i++) {
          cls = plugins[i](cls);
        }
        return cls;
      };
      Parser.parse = function parse2(input, options) {
        return new this(options, input).parse();
      };
      Parser.parseExpressionAt = function parseExpressionAt(input, pos, options) {
        var parser = new this(options, input, pos);
        parser.nextToken();
        return parser.parseExpression();
      };
      Parser.tokenizer = function tokenizer(input, options) {
        return new this(options, input);
      };
      Object.defineProperties(Parser.prototype, prototypeAccessors);
      pp$9 = Parser.prototype;
      literal = /^(?:'((?:\\[^]|[^'\\])*?)'|"((?:\\[^]|[^"\\])*?)")/;
      pp$9.strictDirective = function(start2) {
        if (this.options.ecmaVersion < 5) {
          return false;
        }
        for (; ; ) {
          skipWhiteSpace.lastIndex = start2;
          start2 += skipWhiteSpace.exec(this.input)[0].length;
          var match = literal.exec(this.input.slice(start2));
          if (!match) {
            return false;
          }
          if ((match[1] || match[2]) === "use strict") {
            skipWhiteSpace.lastIndex = start2 + match[0].length;
            var spaceAfter = skipWhiteSpace.exec(this.input), end = spaceAfter.index + spaceAfter[0].length;
            var next = this.input.charAt(end);
            return next === ";" || next === "}" || lineBreak.test(spaceAfter[0]) && !(/[(`.[+\-/*%<>=,?^&]/.test(next) || next === "!" && this.input.charAt(end + 1) === "=");
          }
          start2 += match[0].length;
          skipWhiteSpace.lastIndex = start2;
          start2 += skipWhiteSpace.exec(this.input)[0].length;
          if (this.input[start2] === ";") {
            start2++;
          }
        }
      };
      pp$9.eat = function(type) {
        if (this.type === type) {
          this.next();
          return true;
        } else {
          return false;
        }
      };
      pp$9.isContextual = function(name2) {
        return this.type === types$1.name && this.value === name2 && !this.containsEsc;
      };
      pp$9.eatContextual = function(name2) {
        if (!this.isContextual(name2)) {
          return false;
        }
        this.next();
        return true;
      };
      pp$9.expectContextual = function(name2) {
        if (!this.eatContextual(name2)) {
          this.unexpected();
        }
      };
      pp$9.canInsertSemicolon = function() {
        return this.type === types$1.eof || this.type === types$1.braceR || lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
      };
      pp$9.insertSemicolon = function() {
        if (this.canInsertSemicolon()) {
          if (this.options.onInsertedSemicolon) {
            this.options.onInsertedSemicolon(this.lastTokEnd, this.lastTokEndLoc);
          }
          return true;
        }
      };
      pp$9.semicolon = function() {
        if (!this.eat(types$1.semi) && !this.insertSemicolon()) {
          this.unexpected();
        }
      };
      pp$9.afterTrailingComma = function(tokType, notNext) {
        if (this.type === tokType) {
          if (this.options.onTrailingComma) {
            this.options.onTrailingComma(this.lastTokStart, this.lastTokStartLoc);
          }
          if (!notNext) {
            this.next();
          }
          return true;
        }
      };
      pp$9.expect = function(type) {
        this.eat(type) || this.unexpected();
      };
      pp$9.unexpected = function(pos) {
        this.raise(pos != null ? pos : this.start, "Unexpected token");
      };
      DestructuringErrors = function DestructuringErrors2() {
        this.shorthandAssign = this.trailingComma = this.parenthesizedAssign = this.parenthesizedBind = this.doubleProto = -1;
      };
      pp$9.checkPatternErrors = function(refDestructuringErrors, isAssign) {
        if (!refDestructuringErrors) {
          return;
        }
        if (refDestructuringErrors.trailingComma > -1) {
          this.raiseRecoverable(refDestructuringErrors.trailingComma, "Comma is not permitted after the rest element");
        }
        var parens = isAssign ? refDestructuringErrors.parenthesizedAssign : refDestructuringErrors.parenthesizedBind;
        if (parens > -1) {
          this.raiseRecoverable(parens, isAssign ? "Assigning to rvalue" : "Parenthesized pattern");
        }
      };
      pp$9.checkExpressionErrors = function(refDestructuringErrors, andThrow) {
        if (!refDestructuringErrors) {
          return false;
        }
        var shorthandAssign = refDestructuringErrors.shorthandAssign;
        var doubleProto = refDestructuringErrors.doubleProto;
        if (!andThrow) {
          return shorthandAssign >= 0 || doubleProto >= 0;
        }
        if (shorthandAssign >= 0) {
          this.raise(shorthandAssign, "Shorthand property assignments are valid only in destructuring patterns");
        }
        if (doubleProto >= 0) {
          this.raiseRecoverable(doubleProto, "Redefinition of __proto__ property");
        }
      };
      pp$9.checkYieldAwaitInDefaultParams = function() {
        if (this.yieldPos && (!this.awaitPos || this.yieldPos < this.awaitPos)) {
          this.raise(this.yieldPos, "Yield expression cannot be a default value");
        }
        if (this.awaitPos) {
          this.raise(this.awaitPos, "Await expression cannot be a default value");
        }
      };
      pp$9.isSimpleAssignTarget = function(expr) {
        if (expr.type === "ParenthesizedExpression") {
          return this.isSimpleAssignTarget(expr.expression);
        }
        return expr.type === "Identifier" || expr.type === "MemberExpression";
      };
      pp$8 = Parser.prototype;
      pp$8.parseTopLevel = function(node2) {
        var exports = /* @__PURE__ */ Object.create(null);
        if (!node2.body) {
          node2.body = [];
        }
        while (this.type !== types$1.eof) {
          var stmt = this.parseStatement(null, true, exports);
          node2.body.push(stmt);
        }
        if (this.inModule) {
          for (var i = 0, list4 = Object.keys(this.undefinedExports); i < list4.length; i += 1) {
            var name2 = list4[i];
            this.raiseRecoverable(this.undefinedExports[name2].start, "Export '" + name2 + "' is not defined");
          }
        }
        this.adaptDirectivePrologue(node2.body);
        this.next();
        node2.sourceType = this.options.sourceType;
        return this.finishNode(node2, "Program");
      };
      loopLabel = { kind: "loop" };
      switchLabel = { kind: "switch" };
      pp$8.isLet = function(context) {
        if (this.options.ecmaVersion < 6 || !this.isContextual("let")) {
          return false;
        }
        skipWhiteSpace.lastIndex = this.pos;
        var skip = skipWhiteSpace.exec(this.input);
        var next = this.pos + skip[0].length, nextCh = this.input.charCodeAt(next);
        if (nextCh === 91 || nextCh === 92) {
          return true;
        }
        if (context) {
          return false;
        }
        if (nextCh === 123 || nextCh > 55295 && nextCh < 56320) {
          return true;
        }
        if (isIdentifierStart(nextCh, true)) {
          var pos = next + 1;
          while (isIdentifierChar(nextCh = this.input.charCodeAt(pos), true)) {
            ++pos;
          }
          if (nextCh === 92 || nextCh > 55295 && nextCh < 56320) {
            return true;
          }
          var ident = this.input.slice(next, pos);
          if (!keywordRelationalOperator.test(ident)) {
            return true;
          }
        }
        return false;
      };
      pp$8.isAsyncFunction = function() {
        if (this.options.ecmaVersion < 8 || !this.isContextual("async")) {
          return false;
        }
        skipWhiteSpace.lastIndex = this.pos;
        var skip = skipWhiteSpace.exec(this.input);
        var next = this.pos + skip[0].length, after;
        return !lineBreak.test(this.input.slice(this.pos, next)) && this.input.slice(next, next + 8) === "function" && (next + 8 === this.input.length || !(isIdentifierChar(after = this.input.charCodeAt(next + 8)) || after > 55295 && after < 56320));
      };
      pp$8.isUsingKeyword = function(isAwaitUsing, isFor) {
        if (this.options.ecmaVersion < 17 || !this.isContextual(isAwaitUsing ? "await" : "using")) {
          return false;
        }
        skipWhiteSpace.lastIndex = this.pos;
        var skip = skipWhiteSpace.exec(this.input);
        var next = this.pos + skip[0].length;
        if (lineBreak.test(this.input.slice(this.pos, next))) {
          return false;
        }
        if (isAwaitUsing) {
          var awaitEndPos = next + 5, after;
          if (this.input.slice(next, awaitEndPos) !== "using" || awaitEndPos === this.input.length || isIdentifierChar(after = this.input.charCodeAt(awaitEndPos)) || after > 55295 && after < 56320) {
            return false;
          }
          skipWhiteSpace.lastIndex = awaitEndPos;
          var skipAfterUsing = skipWhiteSpace.exec(this.input);
          if (skipAfterUsing && lineBreak.test(this.input.slice(awaitEndPos, awaitEndPos + skipAfterUsing[0].length))) {
            return false;
          }
        }
        if (isFor) {
          var ofEndPos = next + 2, after$1;
          if (this.input.slice(next, ofEndPos) === "of") {
            if (ofEndPos === this.input.length || !isIdentifierChar(after$1 = this.input.charCodeAt(ofEndPos)) && !(after$1 > 55295 && after$1 < 56320)) {
              return false;
            }
          }
        }
        var ch = this.input.charCodeAt(next);
        return isIdentifierStart(ch, true) || ch === 92;
      };
      pp$8.isAwaitUsing = function(isFor) {
        return this.isUsingKeyword(true, isFor);
      };
      pp$8.isUsing = function(isFor) {
        return this.isUsingKeyword(false, isFor);
      };
      pp$8.parseStatement = function(context, topLevel, exports) {
        var starttype = this.type, node2 = this.startNode(), kind;
        if (this.isLet(context)) {
          starttype = types$1._var;
          kind = "let";
        }
        switch (starttype) {
          case types$1._break:
          case types$1._continue:
            return this.parseBreakContinueStatement(node2, starttype.keyword);
          case types$1._debugger:
            return this.parseDebuggerStatement(node2);
          case types$1._do:
            return this.parseDoStatement(node2);
          case types$1._for:
            return this.parseForStatement(node2);
          case types$1._function:
            if (context && (this.strict || context !== "if" && context !== "label") && this.options.ecmaVersion >= 6) {
              this.unexpected();
            }
            return this.parseFunctionStatement(node2, false, !context);
          case types$1._class:
            if (context) {
              this.unexpected();
            }
            return this.parseClass(node2, true);
          case types$1._if:
            return this.parseIfStatement(node2);
          case types$1._return:
            return this.parseReturnStatement(node2);
          case types$1._switch:
            return this.parseSwitchStatement(node2);
          case types$1._throw:
            return this.parseThrowStatement(node2);
          case types$1._try:
            return this.parseTryStatement(node2);
          case types$1._const:
          case types$1._var:
            kind = kind || this.value;
            if (context && kind !== "var") {
              this.unexpected();
            }
            return this.parseVarStatement(node2, kind);
          case types$1._while:
            return this.parseWhileStatement(node2);
          case types$1._with:
            return this.parseWithStatement(node2);
          case types$1.braceL:
            return this.parseBlock(true, node2);
          case types$1.semi:
            return this.parseEmptyStatement(node2);
          case types$1._export:
          case types$1._import:
            if (this.options.ecmaVersion > 10 && starttype === types$1._import) {
              skipWhiteSpace.lastIndex = this.pos;
              var skip = skipWhiteSpace.exec(this.input);
              var next = this.pos + skip[0].length, nextCh = this.input.charCodeAt(next);
              if (nextCh === 40 || nextCh === 46) {
                return this.parseExpressionStatement(node2, this.parseExpression());
              }
            }
            if (!this.options.allowImportExportEverywhere) {
              if (!topLevel) {
                this.raise(this.start, "'import' and 'export' may only appear at the top level");
              }
              if (!this.inModule) {
                this.raise(this.start, "'import' and 'export' may appear only with 'sourceType: module'");
              }
            }
            return starttype === types$1._import ? this.parseImport(node2) : this.parseExport(node2, exports);
          // If the statement does not start with a statement keyword or a
          // brace, it's an ExpressionStatement or LabeledStatement. We
          // simply start parsing an expression, and afterwards, if the
          // next token is a colon and the expression was a simple
          // Identifier node, we switch to interpreting it as a label.
          default:
            if (this.isAsyncFunction()) {
              if (context) {
                this.unexpected();
              }
              this.next();
              return this.parseFunctionStatement(node2, true, !context);
            }
            var usingKind = this.isAwaitUsing(false) ? "await using" : this.isUsing(false) ? "using" : null;
            if (usingKind) {
              if (topLevel && this.options.sourceType === "script") {
                this.raise(this.start, "Using declaration cannot appear in the top level when source type is `script`");
              }
              if (usingKind === "await using") {
                if (!this.canAwait) {
                  this.raise(this.start, "Await using cannot appear outside of async function");
                }
                this.next();
              }
              this.next();
              this.parseVar(node2, false, usingKind);
              this.semicolon();
              return this.finishNode(node2, "VariableDeclaration");
            }
            var maybeName = this.value, expr = this.parseExpression();
            if (starttype === types$1.name && expr.type === "Identifier" && this.eat(types$1.colon)) {
              return this.parseLabeledStatement(node2, maybeName, expr, context);
            } else {
              return this.parseExpressionStatement(node2, expr);
            }
        }
      };
      pp$8.parseBreakContinueStatement = function(node2, keyword) {
        var isBreak = keyword === "break";
        this.next();
        if (this.eat(types$1.semi) || this.insertSemicolon()) {
          node2.label = null;
        } else if (this.type !== types$1.name) {
          this.unexpected();
        } else {
          node2.label = this.parseIdent();
          this.semicolon();
        }
        var i = 0;
        for (; i < this.labels.length; ++i) {
          var lab = this.labels[i];
          if (node2.label == null || lab.name === node2.label.name) {
            if (lab.kind != null && (isBreak || lab.kind === "loop")) {
              break;
            }
            if (node2.label && isBreak) {
              break;
            }
          }
        }
        if (i === this.labels.length) {
          this.raise(node2.start, "Unsyntactic " + keyword);
        }
        return this.finishNode(node2, isBreak ? "BreakStatement" : "ContinueStatement");
      };
      pp$8.parseDebuggerStatement = function(node2) {
        this.next();
        this.semicolon();
        return this.finishNode(node2, "DebuggerStatement");
      };
      pp$8.parseDoStatement = function(node2) {
        this.next();
        this.labels.push(loopLabel);
        node2.body = this.parseStatement("do");
        this.labels.pop();
        this.expect(types$1._while);
        node2.test = this.parseParenExpression();
        if (this.options.ecmaVersion >= 6) {
          this.eat(types$1.semi);
        } else {
          this.semicolon();
        }
        return this.finishNode(node2, "DoWhileStatement");
      };
      pp$8.parseForStatement = function(node2) {
        this.next();
        var awaitAt = this.options.ecmaVersion >= 9 && this.canAwait && this.eatContextual("await") ? this.lastTokStart : -1;
        this.labels.push(loopLabel);
        this.enterScope(0);
        this.expect(types$1.parenL);
        if (this.type === types$1.semi) {
          if (awaitAt > -1) {
            this.unexpected(awaitAt);
          }
          return this.parseFor(node2, null);
        }
        var isLet = this.isLet();
        if (this.type === types$1._var || this.type === types$1._const || isLet) {
          var init$1 = this.startNode(), kind = isLet ? "let" : this.value;
          this.next();
          this.parseVar(init$1, true, kind);
          this.finishNode(init$1, "VariableDeclaration");
          return this.parseForAfterInit(node2, init$1, awaitAt);
        }
        var startsWithLet = this.isContextual("let"), isForOf = false;
        var usingKind = this.isUsing(true) ? "using" : this.isAwaitUsing(true) ? "await using" : null;
        if (usingKind) {
          var init$2 = this.startNode();
          this.next();
          if (usingKind === "await using") {
            this.next();
          }
          this.parseVar(init$2, true, usingKind);
          this.finishNode(init$2, "VariableDeclaration");
          return this.parseForAfterInit(node2, init$2, awaitAt);
        }
        var containsEsc = this.containsEsc;
        var refDestructuringErrors = new DestructuringErrors();
        var initPos = this.start;
        var init = awaitAt > -1 ? this.parseExprSubscripts(refDestructuringErrors, "await") : this.parseExpression(true, refDestructuringErrors);
        if (this.type === types$1._in || (isForOf = this.options.ecmaVersion >= 6 && this.isContextual("of"))) {
          if (awaitAt > -1) {
            if (this.type === types$1._in) {
              this.unexpected(awaitAt);
            }
            node2.await = true;
          } else if (isForOf && this.options.ecmaVersion >= 8) {
            if (init.start === initPos && !containsEsc && init.type === "Identifier" && init.name === "async") {
              this.unexpected();
            } else if (this.options.ecmaVersion >= 9) {
              node2.await = false;
            }
          }
          if (startsWithLet && isForOf) {
            this.raise(init.start, "The left-hand side of a for-of loop may not start with 'let'.");
          }
          this.toAssignable(init, false, refDestructuringErrors);
          this.checkLValPattern(init);
          return this.parseForIn(node2, init);
        } else {
          this.checkExpressionErrors(refDestructuringErrors, true);
        }
        if (awaitAt > -1) {
          this.unexpected(awaitAt);
        }
        return this.parseFor(node2, init);
      };
      pp$8.parseForAfterInit = function(node2, init, awaitAt) {
        if ((this.type === types$1._in || this.options.ecmaVersion >= 6 && this.isContextual("of")) && init.declarations.length === 1) {
          if (this.options.ecmaVersion >= 9) {
            if (this.type === types$1._in) {
              if (awaitAt > -1) {
                this.unexpected(awaitAt);
              }
            } else {
              node2.await = awaitAt > -1;
            }
          }
          return this.parseForIn(node2, init);
        }
        if (awaitAt > -1) {
          this.unexpected(awaitAt);
        }
        return this.parseFor(node2, init);
      };
      pp$8.parseFunctionStatement = function(node2, isAsync, declarationPosition) {
        this.next();
        return this.parseFunction(node2, FUNC_STATEMENT | (declarationPosition ? 0 : FUNC_HANGING_STATEMENT), false, isAsync);
      };
      pp$8.parseIfStatement = function(node2) {
        this.next();
        node2.test = this.parseParenExpression();
        node2.consequent = this.parseStatement("if");
        node2.alternate = this.eat(types$1._else) ? this.parseStatement("if") : null;
        return this.finishNode(node2, "IfStatement");
      };
      pp$8.parseReturnStatement = function(node2) {
        if (!this.inFunction && !this.options.allowReturnOutsideFunction) {
          this.raise(this.start, "'return' outside of function");
        }
        this.next();
        if (this.eat(types$1.semi) || this.insertSemicolon()) {
          node2.argument = null;
        } else {
          node2.argument = this.parseExpression();
          this.semicolon();
        }
        return this.finishNode(node2, "ReturnStatement");
      };
      pp$8.parseSwitchStatement = function(node2) {
        this.next();
        node2.discriminant = this.parseParenExpression();
        node2.cases = [];
        this.expect(types$1.braceL);
        this.labels.push(switchLabel);
        this.enterScope(0);
        var cur;
        for (var sawDefault = false; this.type !== types$1.braceR; ) {
          if (this.type === types$1._case || this.type === types$1._default) {
            var isCase = this.type === types$1._case;
            if (cur) {
              this.finishNode(cur, "SwitchCase");
            }
            node2.cases.push(cur = this.startNode());
            cur.consequent = [];
            this.next();
            if (isCase) {
              cur.test = this.parseExpression();
            } else {
              if (sawDefault) {
                this.raiseRecoverable(this.lastTokStart, "Multiple default clauses");
              }
              sawDefault = true;
              cur.test = null;
            }
            this.expect(types$1.colon);
          } else {
            if (!cur) {
              this.unexpected();
            }
            cur.consequent.push(this.parseStatement(null));
          }
        }
        this.exitScope();
        if (cur) {
          this.finishNode(cur, "SwitchCase");
        }
        this.next();
        this.labels.pop();
        return this.finishNode(node2, "SwitchStatement");
      };
      pp$8.parseThrowStatement = function(node2) {
        this.next();
        if (lineBreak.test(this.input.slice(this.lastTokEnd, this.start))) {
          this.raise(this.lastTokEnd, "Illegal newline after throw");
        }
        node2.argument = this.parseExpression();
        this.semicolon();
        return this.finishNode(node2, "ThrowStatement");
      };
      empty$1 = [];
      pp$8.parseCatchClauseParam = function() {
        var param = this.parseBindingAtom();
        var simple = param.type === "Identifier";
        this.enterScope(simple ? SCOPE_SIMPLE_CATCH : 0);
        this.checkLValPattern(param, simple ? BIND_SIMPLE_CATCH : BIND_LEXICAL);
        this.expect(types$1.parenR);
        return param;
      };
      pp$8.parseTryStatement = function(node2) {
        this.next();
        node2.block = this.parseBlock();
        node2.handler = null;
        if (this.type === types$1._catch) {
          var clause = this.startNode();
          this.next();
          if (this.eat(types$1.parenL)) {
            clause.param = this.parseCatchClauseParam();
          } else {
            if (this.options.ecmaVersion < 10) {
              this.unexpected();
            }
            clause.param = null;
            this.enterScope(0);
          }
          clause.body = this.parseBlock(false);
          this.exitScope();
          node2.handler = this.finishNode(clause, "CatchClause");
        }
        node2.finalizer = this.eat(types$1._finally) ? this.parseBlock() : null;
        if (!node2.handler && !node2.finalizer) {
          this.raise(node2.start, "Missing catch or finally clause");
        }
        return this.finishNode(node2, "TryStatement");
      };
      pp$8.parseVarStatement = function(node2, kind, allowMissingInitializer) {
        this.next();
        this.parseVar(node2, false, kind, allowMissingInitializer);
        this.semicolon();
        return this.finishNode(node2, "VariableDeclaration");
      };
      pp$8.parseWhileStatement = function(node2) {
        this.next();
        node2.test = this.parseParenExpression();
        this.labels.push(loopLabel);
        node2.body = this.parseStatement("while");
        this.labels.pop();
        return this.finishNode(node2, "WhileStatement");
      };
      pp$8.parseWithStatement = function(node2) {
        if (this.strict) {
          this.raise(this.start, "'with' in strict mode");
        }
        this.next();
        node2.object = this.parseParenExpression();
        node2.body = this.parseStatement("with");
        return this.finishNode(node2, "WithStatement");
      };
      pp$8.parseEmptyStatement = function(node2) {
        this.next();
        return this.finishNode(node2, "EmptyStatement");
      };
      pp$8.parseLabeledStatement = function(node2, maybeName, expr, context) {
        for (var i$1 = 0, list4 = this.labels; i$1 < list4.length; i$1 += 1) {
          var label = list4[i$1];
          if (label.name === maybeName) {
            this.raise(expr.start, "Label '" + maybeName + "' is already declared");
          }
        }
        var kind = this.type.isLoop ? "loop" : this.type === types$1._switch ? "switch" : null;
        for (var i = this.labels.length - 1; i >= 0; i--) {
          var label$1 = this.labels[i];
          if (label$1.statementStart === node2.start) {
            label$1.statementStart = this.start;
            label$1.kind = kind;
          } else {
            break;
          }
        }
        this.labels.push({ name: maybeName, kind, statementStart: this.start });
        node2.body = this.parseStatement(context ? context.indexOf("label") === -1 ? context + "label" : context : "label");
        this.labels.pop();
        node2.label = expr;
        return this.finishNode(node2, "LabeledStatement");
      };
      pp$8.parseExpressionStatement = function(node2, expr) {
        node2.expression = expr;
        this.semicolon();
        return this.finishNode(node2, "ExpressionStatement");
      };
      pp$8.parseBlock = function(createNewLexicalScope, node2, exitStrict) {
        if (createNewLexicalScope === void 0) createNewLexicalScope = true;
        if (node2 === void 0) node2 = this.startNode();
        node2.body = [];
        this.expect(types$1.braceL);
        if (createNewLexicalScope) {
          this.enterScope(0);
        }
        while (this.type !== types$1.braceR) {
          var stmt = this.parseStatement(null);
          node2.body.push(stmt);
        }
        if (exitStrict) {
          this.strict = false;
        }
        this.next();
        if (createNewLexicalScope) {
          this.exitScope();
        }
        return this.finishNode(node2, "BlockStatement");
      };
      pp$8.parseFor = function(node2, init) {
        node2.init = init;
        this.expect(types$1.semi);
        node2.test = this.type === types$1.semi ? null : this.parseExpression();
        this.expect(types$1.semi);
        node2.update = this.type === types$1.parenR ? null : this.parseExpression();
        this.expect(types$1.parenR);
        node2.body = this.parseStatement("for");
        this.exitScope();
        this.labels.pop();
        return this.finishNode(node2, "ForStatement");
      };
      pp$8.parseForIn = function(node2, init) {
        var isForIn = this.type === types$1._in;
        this.next();
        if (init.type === "VariableDeclaration" && init.declarations[0].init != null && (!isForIn || this.options.ecmaVersion < 8 || this.strict || init.kind !== "var" || init.declarations[0].id.type !== "Identifier")) {
          this.raise(
            init.start,
            (isForIn ? "for-in" : "for-of") + " loop variable declaration may not have an initializer"
          );
        }
        node2.left = init;
        node2.right = isForIn ? this.parseExpression() : this.parseMaybeAssign();
        this.expect(types$1.parenR);
        node2.body = this.parseStatement("for");
        this.exitScope();
        this.labels.pop();
        return this.finishNode(node2, isForIn ? "ForInStatement" : "ForOfStatement");
      };
      pp$8.parseVar = function(node2, isFor, kind, allowMissingInitializer) {
        node2.declarations = [];
        node2.kind = kind;
        for (; ; ) {
          var decl = this.startNode();
          this.parseVarId(decl, kind);
          if (this.eat(types$1.eq)) {
            decl.init = this.parseMaybeAssign(isFor);
          } else if (!allowMissingInitializer && kind === "const" && !(this.type === types$1._in || this.options.ecmaVersion >= 6 && this.isContextual("of"))) {
            this.unexpected();
          } else if (!allowMissingInitializer && (kind === "using" || kind === "await using") && this.options.ecmaVersion >= 17 && this.type !== types$1._in && !this.isContextual("of")) {
            this.raise(this.lastTokEnd, "Missing initializer in " + kind + " declaration");
          } else if (!allowMissingInitializer && decl.id.type !== "Identifier" && !(isFor && (this.type === types$1._in || this.isContextual("of")))) {
            this.raise(this.lastTokEnd, "Complex binding patterns require an initialization value");
          } else {
            decl.init = null;
          }
          node2.declarations.push(this.finishNode(decl, "VariableDeclarator"));
          if (!this.eat(types$1.comma)) {
            break;
          }
        }
        return node2;
      };
      pp$8.parseVarId = function(decl, kind) {
        decl.id = kind === "using" || kind === "await using" ? this.parseIdent() : this.parseBindingAtom();
        this.checkLValPattern(decl.id, kind === "var" ? BIND_VAR : BIND_LEXICAL, false);
      };
      FUNC_STATEMENT = 1;
      FUNC_HANGING_STATEMENT = 2;
      FUNC_NULLABLE_ID = 4;
      pp$8.parseFunction = function(node2, statement, allowExpressionBody, isAsync, forInit) {
        this.initFunction(node2);
        if (this.options.ecmaVersion >= 9 || this.options.ecmaVersion >= 6 && !isAsync) {
          if (this.type === types$1.star && statement & FUNC_HANGING_STATEMENT) {
            this.unexpected();
          }
          node2.generator = this.eat(types$1.star);
        }
        if (this.options.ecmaVersion >= 8) {
          node2.async = !!isAsync;
        }
        if (statement & FUNC_STATEMENT) {
          node2.id = statement & FUNC_NULLABLE_ID && this.type !== types$1.name ? null : this.parseIdent();
          if (node2.id && !(statement & FUNC_HANGING_STATEMENT)) {
            this.checkLValSimple(node2.id, this.strict || node2.generator || node2.async ? this.treatFunctionsAsVar ? BIND_VAR : BIND_LEXICAL : BIND_FUNCTION);
          }
        }
        var oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
        this.yieldPos = 0;
        this.awaitPos = 0;
        this.awaitIdentPos = 0;
        this.enterScope(functionFlags(node2.async, node2.generator));
        if (!(statement & FUNC_STATEMENT)) {
          node2.id = this.type === types$1.name ? this.parseIdent() : null;
        }
        this.parseFunctionParams(node2);
        this.parseFunctionBody(node2, allowExpressionBody, false, forInit);
        this.yieldPos = oldYieldPos;
        this.awaitPos = oldAwaitPos;
        this.awaitIdentPos = oldAwaitIdentPos;
        return this.finishNode(node2, statement & FUNC_STATEMENT ? "FunctionDeclaration" : "FunctionExpression");
      };
      pp$8.parseFunctionParams = function(node2) {
        this.expect(types$1.parenL);
        node2.params = this.parseBindingList(types$1.parenR, false, this.options.ecmaVersion >= 8);
        this.checkYieldAwaitInDefaultParams();
      };
      pp$8.parseClass = function(node2, isStatement) {
        this.next();
        var oldStrict = this.strict;
        this.strict = true;
        this.parseClassId(node2, isStatement);
        this.parseClassSuper(node2);
        var privateNameMap = this.enterClassBody();
        var classBody = this.startNode();
        var hadConstructor = false;
        classBody.body = [];
        this.expect(types$1.braceL);
        while (this.type !== types$1.braceR) {
          var element = this.parseClassElement(node2.superClass !== null);
          if (element) {
            classBody.body.push(element);
            if (element.type === "MethodDefinition" && element.kind === "constructor") {
              if (hadConstructor) {
                this.raiseRecoverable(element.start, "Duplicate constructor in the same class");
              }
              hadConstructor = true;
            } else if (element.key && element.key.type === "PrivateIdentifier" && isPrivateNameConflicted(privateNameMap, element)) {
              this.raiseRecoverable(element.key.start, "Identifier '#" + element.key.name + "' has already been declared");
            }
          }
        }
        this.strict = oldStrict;
        this.next();
        node2.body = this.finishNode(classBody, "ClassBody");
        this.exitClassBody();
        return this.finishNode(node2, isStatement ? "ClassDeclaration" : "ClassExpression");
      };
      pp$8.parseClassElement = function(constructorAllowsSuper) {
        if (this.eat(types$1.semi)) {
          return null;
        }
        var ecmaVersion = this.options.ecmaVersion;
        var node2 = this.startNode();
        var keyName = "";
        var isGenerator = false;
        var isAsync = false;
        var kind = "method";
        var isStatic = false;
        if (this.eatContextual("static")) {
          if (ecmaVersion >= 13 && this.eat(types$1.braceL)) {
            this.parseClassStaticBlock(node2);
            return node2;
          }
          if (this.isClassElementNameStart() || this.type === types$1.star) {
            isStatic = true;
          } else {
            keyName = "static";
          }
        }
        node2.static = isStatic;
        if (!keyName && ecmaVersion >= 8 && this.eatContextual("async")) {
          if ((this.isClassElementNameStart() || this.type === types$1.star) && !this.canInsertSemicolon()) {
            isAsync = true;
          } else {
            keyName = "async";
          }
        }
        if (!keyName && (ecmaVersion >= 9 || !isAsync) && this.eat(types$1.star)) {
          isGenerator = true;
        }
        if (!keyName && !isAsync && !isGenerator) {
          var lastValue = this.value;
          if (this.eatContextual("get") || this.eatContextual("set")) {
            if (this.isClassElementNameStart()) {
              kind = lastValue;
            } else {
              keyName = lastValue;
            }
          }
        }
        if (keyName) {
          node2.computed = false;
          node2.key = this.startNodeAt(this.lastTokStart, this.lastTokStartLoc);
          node2.key.name = keyName;
          this.finishNode(node2.key, "Identifier");
        } else {
          this.parseClassElementName(node2);
        }
        if (ecmaVersion < 13 || this.type === types$1.parenL || kind !== "method" || isGenerator || isAsync) {
          var isConstructor = !node2.static && checkKeyName(node2, "constructor");
          var allowsDirectSuper = isConstructor && constructorAllowsSuper;
          if (isConstructor && kind !== "method") {
            this.raise(node2.key.start, "Constructor can't have get/set modifier");
          }
          node2.kind = isConstructor ? "constructor" : kind;
          this.parseClassMethod(node2, isGenerator, isAsync, allowsDirectSuper);
        } else {
          this.parseClassField(node2);
        }
        return node2;
      };
      pp$8.isClassElementNameStart = function() {
        return this.type === types$1.name || this.type === types$1.privateId || this.type === types$1.num || this.type === types$1.string || this.type === types$1.bracketL || this.type.keyword;
      };
      pp$8.parseClassElementName = function(element) {
        if (this.type === types$1.privateId) {
          if (this.value === "constructor") {
            this.raise(this.start, "Classes can't have an element named '#constructor'");
          }
          element.computed = false;
          element.key = this.parsePrivateIdent();
        } else {
          this.parsePropertyName(element);
        }
      };
      pp$8.parseClassMethod = function(method, isGenerator, isAsync, allowsDirectSuper) {
        var key = method.key;
        if (method.kind === "constructor") {
          if (isGenerator) {
            this.raise(key.start, "Constructor can't be a generator");
          }
          if (isAsync) {
            this.raise(key.start, "Constructor can't be an async method");
          }
        } else if (method.static && checkKeyName(method, "prototype")) {
          this.raise(key.start, "Classes may not have a static property named prototype");
        }
        var value = method.value = this.parseMethod(isGenerator, isAsync, allowsDirectSuper);
        if (method.kind === "get" && value.params.length !== 0) {
          this.raiseRecoverable(value.start, "getter should have no params");
        }
        if (method.kind === "set" && value.params.length !== 1) {
          this.raiseRecoverable(value.start, "setter should have exactly one param");
        }
        if (method.kind === "set" && value.params[0].type === "RestElement") {
          this.raiseRecoverable(value.params[0].start, "Setter cannot use rest params");
        }
        return this.finishNode(method, "MethodDefinition");
      };
      pp$8.parseClassField = function(field) {
        if (checkKeyName(field, "constructor")) {
          this.raise(field.key.start, "Classes can't have a field named 'constructor'");
        } else if (field.static && checkKeyName(field, "prototype")) {
          this.raise(field.key.start, "Classes can't have a static field named 'prototype'");
        }
        if (this.eat(types$1.eq)) {
          this.enterScope(SCOPE_CLASS_FIELD_INIT | SCOPE_SUPER);
          field.value = this.parseMaybeAssign();
          this.exitScope();
        } else {
          field.value = null;
        }
        this.semicolon();
        return this.finishNode(field, "PropertyDefinition");
      };
      pp$8.parseClassStaticBlock = function(node2) {
        node2.body = [];
        var oldLabels = this.labels;
        this.labels = [];
        this.enterScope(SCOPE_CLASS_STATIC_BLOCK | SCOPE_SUPER);
        while (this.type !== types$1.braceR) {
          var stmt = this.parseStatement(null);
          node2.body.push(stmt);
        }
        this.next();
        this.exitScope();
        this.labels = oldLabels;
        return this.finishNode(node2, "StaticBlock");
      };
      pp$8.parseClassId = function(node2, isStatement) {
        if (this.type === types$1.name) {
          node2.id = this.parseIdent();
          if (isStatement) {
            this.checkLValSimple(node2.id, BIND_LEXICAL, false);
          }
        } else {
          if (isStatement === true) {
            this.unexpected();
          }
          node2.id = null;
        }
      };
      pp$8.parseClassSuper = function(node2) {
        node2.superClass = this.eat(types$1._extends) ? this.parseExprSubscripts(null, false) : null;
      };
      pp$8.enterClassBody = function() {
        var element = { declared: /* @__PURE__ */ Object.create(null), used: [] };
        this.privateNameStack.push(element);
        return element.declared;
      };
      pp$8.exitClassBody = function() {
        var ref2 = this.privateNameStack.pop();
        var declared = ref2.declared;
        var used = ref2.used;
        if (!this.options.checkPrivateFields) {
          return;
        }
        var len = this.privateNameStack.length;
        var parent = len === 0 ? null : this.privateNameStack[len - 1];
        for (var i = 0; i < used.length; ++i) {
          var id = used[i];
          if (!hasOwn(declared, id.name)) {
            if (parent) {
              parent.used.push(id);
            } else {
              this.raiseRecoverable(id.start, "Private field '#" + id.name + "' must be declared in an enclosing class");
            }
          }
        }
      };
      pp$8.parseExportAllDeclaration = function(node2, exports) {
        if (this.options.ecmaVersion >= 11) {
          if (this.eatContextual("as")) {
            node2.exported = this.parseModuleExportName();
            this.checkExport(exports, node2.exported, this.lastTokStart);
          } else {
            node2.exported = null;
          }
        }
        this.expectContextual("from");
        if (this.type !== types$1.string) {
          this.unexpected();
        }
        node2.source = this.parseExprAtom();
        if (this.options.ecmaVersion >= 16) {
          node2.attributes = this.parseWithClause();
        }
        this.semicolon();
        return this.finishNode(node2, "ExportAllDeclaration");
      };
      pp$8.parseExport = function(node2, exports) {
        this.next();
        if (this.eat(types$1.star)) {
          return this.parseExportAllDeclaration(node2, exports);
        }
        if (this.eat(types$1._default)) {
          this.checkExport(exports, "default", this.lastTokStart);
          node2.declaration = this.parseExportDefaultDeclaration();
          return this.finishNode(node2, "ExportDefaultDeclaration");
        }
        if (this.shouldParseExportStatement()) {
          node2.declaration = this.parseExportDeclaration(node2);
          if (node2.declaration.type === "VariableDeclaration") {
            this.checkVariableExport(exports, node2.declaration.declarations);
          } else {
            this.checkExport(exports, node2.declaration.id, node2.declaration.id.start);
          }
          node2.specifiers = [];
          node2.source = null;
          if (this.options.ecmaVersion >= 16) {
            node2.attributes = [];
          }
        } else {
          node2.declaration = null;
          node2.specifiers = this.parseExportSpecifiers(exports);
          if (this.eatContextual("from")) {
            if (this.type !== types$1.string) {
              this.unexpected();
            }
            node2.source = this.parseExprAtom();
            if (this.options.ecmaVersion >= 16) {
              node2.attributes = this.parseWithClause();
            }
          } else {
            for (var i = 0, list4 = node2.specifiers; i < list4.length; i += 1) {
              var spec = list4[i];
              this.checkUnreserved(spec.local);
              this.checkLocalExport(spec.local);
              if (spec.local.type === "Literal") {
                this.raise(spec.local.start, "A string literal cannot be used as an exported binding without `from`.");
              }
            }
            node2.source = null;
            if (this.options.ecmaVersion >= 16) {
              node2.attributes = [];
            }
          }
          this.semicolon();
        }
        return this.finishNode(node2, "ExportNamedDeclaration");
      };
      pp$8.parseExportDeclaration = function(node2) {
        return this.parseStatement(null);
      };
      pp$8.parseExportDefaultDeclaration = function() {
        var isAsync;
        if (this.type === types$1._function || (isAsync = this.isAsyncFunction())) {
          var fNode = this.startNode();
          this.next();
          if (isAsync) {
            this.next();
          }
          return this.parseFunction(fNode, FUNC_STATEMENT | FUNC_NULLABLE_ID, false, isAsync);
        } else if (this.type === types$1._class) {
          var cNode = this.startNode();
          return this.parseClass(cNode, "nullableID");
        } else {
          var declaration = this.parseMaybeAssign();
          this.semicolon();
          return declaration;
        }
      };
      pp$8.checkExport = function(exports, name2, pos) {
        if (!exports) {
          return;
        }
        if (typeof name2 !== "string") {
          name2 = name2.type === "Identifier" ? name2.name : name2.value;
        }
        if (hasOwn(exports, name2)) {
          this.raiseRecoverable(pos, "Duplicate export '" + name2 + "'");
        }
        exports[name2] = true;
      };
      pp$8.checkPatternExport = function(exports, pat) {
        var type = pat.type;
        if (type === "Identifier") {
          this.checkExport(exports, pat, pat.start);
        } else if (type === "ObjectPattern") {
          for (var i = 0, list4 = pat.properties; i < list4.length; i += 1) {
            var prop = list4[i];
            this.checkPatternExport(exports, prop);
          }
        } else if (type === "ArrayPattern") {
          for (var i$1 = 0, list$1 = pat.elements; i$1 < list$1.length; i$1 += 1) {
            var elt = list$1[i$1];
            if (elt) {
              this.checkPatternExport(exports, elt);
            }
          }
        } else if (type === "Property") {
          this.checkPatternExport(exports, pat.value);
        } else if (type === "AssignmentPattern") {
          this.checkPatternExport(exports, pat.left);
        } else if (type === "RestElement") {
          this.checkPatternExport(exports, pat.argument);
        }
      };
      pp$8.checkVariableExport = function(exports, decls) {
        if (!exports) {
          return;
        }
        for (var i = 0, list4 = decls; i < list4.length; i += 1) {
          var decl = list4[i];
          this.checkPatternExport(exports, decl.id);
        }
      };
      pp$8.shouldParseExportStatement = function() {
        return this.type.keyword === "var" || this.type.keyword === "const" || this.type.keyword === "class" || this.type.keyword === "function" || this.isLet() || this.isAsyncFunction();
      };
      pp$8.parseExportSpecifier = function(exports) {
        var node2 = this.startNode();
        node2.local = this.parseModuleExportName();
        node2.exported = this.eatContextual("as") ? this.parseModuleExportName() : node2.local;
        this.checkExport(
          exports,
          node2.exported,
          node2.exported.start
        );
        return this.finishNode(node2, "ExportSpecifier");
      };
      pp$8.parseExportSpecifiers = function(exports) {
        var nodes = [], first = true;
        this.expect(types$1.braceL);
        while (!this.eat(types$1.braceR)) {
          if (!first) {
            this.expect(types$1.comma);
            if (this.afterTrailingComma(types$1.braceR)) {
              break;
            }
          } else {
            first = false;
          }
          nodes.push(this.parseExportSpecifier(exports));
        }
        return nodes;
      };
      pp$8.parseImport = function(node2) {
        this.next();
        if (this.type === types$1.string) {
          node2.specifiers = empty$1;
          node2.source = this.parseExprAtom();
        } else {
          node2.specifiers = this.parseImportSpecifiers();
          this.expectContextual("from");
          node2.source = this.type === types$1.string ? this.parseExprAtom() : this.unexpected();
        }
        if (this.options.ecmaVersion >= 16) {
          node2.attributes = this.parseWithClause();
        }
        this.semicolon();
        return this.finishNode(node2, "ImportDeclaration");
      };
      pp$8.parseImportSpecifier = function() {
        var node2 = this.startNode();
        node2.imported = this.parseModuleExportName();
        if (this.eatContextual("as")) {
          node2.local = this.parseIdent();
        } else {
          this.checkUnreserved(node2.imported);
          node2.local = node2.imported;
        }
        this.checkLValSimple(node2.local, BIND_LEXICAL);
        return this.finishNode(node2, "ImportSpecifier");
      };
      pp$8.parseImportDefaultSpecifier = function() {
        var node2 = this.startNode();
        node2.local = this.parseIdent();
        this.checkLValSimple(node2.local, BIND_LEXICAL);
        return this.finishNode(node2, "ImportDefaultSpecifier");
      };
      pp$8.parseImportNamespaceSpecifier = function() {
        var node2 = this.startNode();
        this.next();
        this.expectContextual("as");
        node2.local = this.parseIdent();
        this.checkLValSimple(node2.local, BIND_LEXICAL);
        return this.finishNode(node2, "ImportNamespaceSpecifier");
      };
      pp$8.parseImportSpecifiers = function() {
        var nodes = [], first = true;
        if (this.type === types$1.name) {
          nodes.push(this.parseImportDefaultSpecifier());
          if (!this.eat(types$1.comma)) {
            return nodes;
          }
        }
        if (this.type === types$1.star) {
          nodes.push(this.parseImportNamespaceSpecifier());
          return nodes;
        }
        this.expect(types$1.braceL);
        while (!this.eat(types$1.braceR)) {
          if (!first) {
            this.expect(types$1.comma);
            if (this.afterTrailingComma(types$1.braceR)) {
              break;
            }
          } else {
            first = false;
          }
          nodes.push(this.parseImportSpecifier());
        }
        return nodes;
      };
      pp$8.parseWithClause = function() {
        var nodes = [];
        if (!this.eat(types$1._with)) {
          return nodes;
        }
        this.expect(types$1.braceL);
        var attributeKeys = {};
        var first = true;
        while (!this.eat(types$1.braceR)) {
          if (!first) {
            this.expect(types$1.comma);
            if (this.afterTrailingComma(types$1.braceR)) {
              break;
            }
          } else {
            first = false;
          }
          var attr = this.parseImportAttribute();
          var keyName = attr.key.type === "Identifier" ? attr.key.name : attr.key.value;
          if (hasOwn(attributeKeys, keyName)) {
            this.raiseRecoverable(attr.key.start, "Duplicate attribute key '" + keyName + "'");
          }
          attributeKeys[keyName] = true;
          nodes.push(attr);
        }
        return nodes;
      };
      pp$8.parseImportAttribute = function() {
        var node2 = this.startNode();
        node2.key = this.type === types$1.string ? this.parseExprAtom() : this.parseIdent(this.options.allowReserved !== "never");
        this.expect(types$1.colon);
        if (this.type !== types$1.string) {
          this.unexpected();
        }
        node2.value = this.parseExprAtom();
        return this.finishNode(node2, "ImportAttribute");
      };
      pp$8.parseModuleExportName = function() {
        if (this.options.ecmaVersion >= 13 && this.type === types$1.string) {
          var stringLiteral = this.parseLiteral(this.value);
          if (loneSurrogate.test(stringLiteral.value)) {
            this.raise(stringLiteral.start, "An export name cannot include a lone surrogate.");
          }
          return stringLiteral;
        }
        return this.parseIdent(true);
      };
      pp$8.adaptDirectivePrologue = function(statements) {
        for (var i = 0; i < statements.length && this.isDirectiveCandidate(statements[i]); ++i) {
          statements[i].directive = statements[i].expression.raw.slice(1, -1);
        }
      };
      pp$8.isDirectiveCandidate = function(statement) {
        return this.options.ecmaVersion >= 5 && statement.type === "ExpressionStatement" && statement.expression.type === "Literal" && typeof statement.expression.value === "string" && // Reject parenthesized strings.
        (this.input[statement.start] === '"' || this.input[statement.start] === "'");
      };
      pp$7 = Parser.prototype;
      pp$7.toAssignable = function(node2, isBinding, refDestructuringErrors) {
        if (this.options.ecmaVersion >= 6 && node2) {
          switch (node2.type) {
            case "Identifier":
              if (this.inAsync && node2.name === "await") {
                this.raise(node2.start, "Cannot use 'await' as identifier inside an async function");
              }
              break;
            case "ObjectPattern":
            case "ArrayPattern":
            case "AssignmentPattern":
            case "RestElement":
              break;
            case "ObjectExpression":
              node2.type = "ObjectPattern";
              if (refDestructuringErrors) {
                this.checkPatternErrors(refDestructuringErrors, true);
              }
              for (var i = 0, list4 = node2.properties; i < list4.length; i += 1) {
                var prop = list4[i];
                this.toAssignable(prop, isBinding);
                if (prop.type === "RestElement" && (prop.argument.type === "ArrayPattern" || prop.argument.type === "ObjectPattern")) {
                  this.raise(prop.argument.start, "Unexpected token");
                }
              }
              break;
            case "Property":
              if (node2.kind !== "init") {
                this.raise(node2.key.start, "Object pattern can't contain getter or setter");
              }
              this.toAssignable(node2.value, isBinding);
              break;
            case "ArrayExpression":
              node2.type = "ArrayPattern";
              if (refDestructuringErrors) {
                this.checkPatternErrors(refDestructuringErrors, true);
              }
              this.toAssignableList(node2.elements, isBinding);
              break;
            case "SpreadElement":
              node2.type = "RestElement";
              this.toAssignable(node2.argument, isBinding);
              if (node2.argument.type === "AssignmentPattern") {
                this.raise(node2.argument.start, "Rest elements cannot have a default value");
              }
              break;
            case "AssignmentExpression":
              if (node2.operator !== "=") {
                this.raise(node2.left.end, "Only '=' operator can be used for specifying default value.");
              }
              node2.type = "AssignmentPattern";
              delete node2.operator;
              this.toAssignable(node2.left, isBinding);
              break;
            case "ParenthesizedExpression":
              this.toAssignable(node2.expression, isBinding, refDestructuringErrors);
              break;
            case "ChainExpression":
              this.raiseRecoverable(node2.start, "Optional chaining cannot appear in left-hand side");
              break;
            case "MemberExpression":
              if (!isBinding) {
                break;
              }
            default:
              this.raise(node2.start, "Assigning to rvalue");
          }
        } else if (refDestructuringErrors) {
          this.checkPatternErrors(refDestructuringErrors, true);
        }
        return node2;
      };
      pp$7.toAssignableList = function(exprList, isBinding) {
        var end = exprList.length;
        for (var i = 0; i < end; i++) {
          var elt = exprList[i];
          if (elt) {
            this.toAssignable(elt, isBinding);
          }
        }
        if (end) {
          var last = exprList[end - 1];
          if (this.options.ecmaVersion === 6 && isBinding && last && last.type === "RestElement" && last.argument.type !== "Identifier") {
            this.unexpected(last.argument.start);
          }
        }
        return exprList;
      };
      pp$7.parseSpread = function(refDestructuringErrors) {
        var node2 = this.startNode();
        this.next();
        node2.argument = this.parseMaybeAssign(false, refDestructuringErrors);
        return this.finishNode(node2, "SpreadElement");
      };
      pp$7.parseRestBinding = function() {
        var node2 = this.startNode();
        this.next();
        if (this.options.ecmaVersion === 6 && this.type !== types$1.name) {
          this.unexpected();
        }
        node2.argument = this.parseBindingAtom();
        return this.finishNode(node2, "RestElement");
      };
      pp$7.parseBindingAtom = function() {
        if (this.options.ecmaVersion >= 6) {
          switch (this.type) {
            case types$1.bracketL:
              var node2 = this.startNode();
              this.next();
              node2.elements = this.parseBindingList(types$1.bracketR, true, true);
              return this.finishNode(node2, "ArrayPattern");
            case types$1.braceL:
              return this.parseObj(true);
          }
        }
        return this.parseIdent();
      };
      pp$7.parseBindingList = function(close, allowEmpty, allowTrailingComma, allowModifiers) {
        var elts = [], first = true;
        while (!this.eat(close)) {
          if (first) {
            first = false;
          } else {
            this.expect(types$1.comma);
          }
          if (allowEmpty && this.type === types$1.comma) {
            elts.push(null);
          } else if (allowTrailingComma && this.afterTrailingComma(close)) {
            break;
          } else if (this.type === types$1.ellipsis) {
            var rest = this.parseRestBinding();
            this.parseBindingListItem(rest);
            elts.push(rest);
            if (this.type === types$1.comma) {
              this.raiseRecoverable(this.start, "Comma is not permitted after the rest element");
            }
            this.expect(close);
            break;
          } else {
            elts.push(this.parseAssignableListItem(allowModifiers));
          }
        }
        return elts;
      };
      pp$7.parseAssignableListItem = function(allowModifiers) {
        var elem = this.parseMaybeDefault(this.start, this.startLoc);
        this.parseBindingListItem(elem);
        return elem;
      };
      pp$7.parseBindingListItem = function(param) {
        return param;
      };
      pp$7.parseMaybeDefault = function(startPos, startLoc, left) {
        left = left || this.parseBindingAtom();
        if (this.options.ecmaVersion < 6 || !this.eat(types$1.eq)) {
          return left;
        }
        var node2 = this.startNodeAt(startPos, startLoc);
        node2.left = left;
        node2.right = this.parseMaybeAssign();
        return this.finishNode(node2, "AssignmentPattern");
      };
      pp$7.checkLValSimple = function(expr, bindingType, checkClashes) {
        if (bindingType === void 0) bindingType = BIND_NONE;
        var isBind = bindingType !== BIND_NONE;
        switch (expr.type) {
          case "Identifier":
            if (this.strict && this.reservedWordsStrictBind.test(expr.name)) {
              this.raiseRecoverable(expr.start, (isBind ? "Binding " : "Assigning to ") + expr.name + " in strict mode");
            }
            if (isBind) {
              if (bindingType === BIND_LEXICAL && expr.name === "let") {
                this.raiseRecoverable(expr.start, "let is disallowed as a lexically bound name");
              }
              if (checkClashes) {
                if (hasOwn(checkClashes, expr.name)) {
                  this.raiseRecoverable(expr.start, "Argument name clash");
                }
                checkClashes[expr.name] = true;
              }
              if (bindingType !== BIND_OUTSIDE) {
                this.declareName(expr.name, bindingType, expr.start);
              }
            }
            break;
          case "ChainExpression":
            this.raiseRecoverable(expr.start, "Optional chaining cannot appear in left-hand side");
            break;
          case "MemberExpression":
            if (isBind) {
              this.raiseRecoverable(expr.start, "Binding member expression");
            }
            break;
          case "ParenthesizedExpression":
            if (isBind) {
              this.raiseRecoverable(expr.start, "Binding parenthesized expression");
            }
            return this.checkLValSimple(expr.expression, bindingType, checkClashes);
          default:
            this.raise(expr.start, (isBind ? "Binding" : "Assigning to") + " rvalue");
        }
      };
      pp$7.checkLValPattern = function(expr, bindingType, checkClashes) {
        if (bindingType === void 0) bindingType = BIND_NONE;
        switch (expr.type) {
          case "ObjectPattern":
            for (var i = 0, list4 = expr.properties; i < list4.length; i += 1) {
              var prop = list4[i];
              this.checkLValInnerPattern(prop, bindingType, checkClashes);
            }
            break;
          case "ArrayPattern":
            for (var i$1 = 0, list$1 = expr.elements; i$1 < list$1.length; i$1 += 1) {
              var elem = list$1[i$1];
              if (elem) {
                this.checkLValInnerPattern(elem, bindingType, checkClashes);
              }
            }
            break;
          default:
            this.checkLValSimple(expr, bindingType, checkClashes);
        }
      };
      pp$7.checkLValInnerPattern = function(expr, bindingType, checkClashes) {
        if (bindingType === void 0) bindingType = BIND_NONE;
        switch (expr.type) {
          case "Property":
            this.checkLValInnerPattern(expr.value, bindingType, checkClashes);
            break;
          case "AssignmentPattern":
            this.checkLValPattern(expr.left, bindingType, checkClashes);
            break;
          case "RestElement":
            this.checkLValPattern(expr.argument, bindingType, checkClashes);
            break;
          default:
            this.checkLValPattern(expr, bindingType, checkClashes);
        }
      };
      TokContext = function TokContext2(token, isExpr, preserveSpace, override, generator) {
        this.token = token;
        this.isExpr = !!isExpr;
        this.preserveSpace = !!preserveSpace;
        this.override = override;
        this.generator = !!generator;
      };
      types = {
        b_stat: new TokContext("{", false),
        b_expr: new TokContext("{", true),
        b_tmpl: new TokContext("${", false),
        p_stat: new TokContext("(", false),
        p_expr: new TokContext("(", true),
        q_tmpl: new TokContext("`", true, true, function(p) {
          return p.tryReadTemplateToken();
        }),
        f_stat: new TokContext("function", false),
        f_expr: new TokContext("function", true),
        f_expr_gen: new TokContext("function", true, false, null, true),
        f_gen: new TokContext("function", false, false, null, true)
      };
      pp$6 = Parser.prototype;
      pp$6.initialContext = function() {
        return [types.b_stat];
      };
      pp$6.curContext = function() {
        return this.context[this.context.length - 1];
      };
      pp$6.braceIsBlock = function(prevType) {
        var parent = this.curContext();
        if (parent === types.f_expr || parent === types.f_stat) {
          return true;
        }
        if (prevType === types$1.colon && (parent === types.b_stat || parent === types.b_expr)) {
          return !parent.isExpr;
        }
        if (prevType === types$1._return || prevType === types$1.name && this.exprAllowed) {
          return lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
        }
        if (prevType === types$1._else || prevType === types$1.semi || prevType === types$1.eof || prevType === types$1.parenR || prevType === types$1.arrow) {
          return true;
        }
        if (prevType === types$1.braceL) {
          return parent === types.b_stat;
        }
        if (prevType === types$1._var || prevType === types$1._const || prevType === types$1.name) {
          return false;
        }
        return !this.exprAllowed;
      };
      pp$6.inGeneratorContext = function() {
        for (var i = this.context.length - 1; i >= 1; i--) {
          var context = this.context[i];
          if (context.token === "function") {
            return context.generator;
          }
        }
        return false;
      };
      pp$6.updateContext = function(prevType) {
        var update2, type = this.type;
        if (type.keyword && prevType === types$1.dot) {
          this.exprAllowed = false;
        } else if (update2 = type.updateContext) {
          update2.call(this, prevType);
        } else {
          this.exprAllowed = type.beforeExpr;
        }
      };
      pp$6.overrideContext = function(tokenCtx) {
        if (this.curContext() !== tokenCtx) {
          this.context[this.context.length - 1] = tokenCtx;
        }
      };
      types$1.parenR.updateContext = types$1.braceR.updateContext = function() {
        if (this.context.length === 1) {
          this.exprAllowed = true;
          return;
        }
        var out = this.context.pop();
        if (out === types.b_stat && this.curContext().token === "function") {
          out = this.context.pop();
        }
        this.exprAllowed = !out.isExpr;
      };
      types$1.braceL.updateContext = function(prevType) {
        this.context.push(this.braceIsBlock(prevType) ? types.b_stat : types.b_expr);
        this.exprAllowed = true;
      };
      types$1.dollarBraceL.updateContext = function() {
        this.context.push(types.b_tmpl);
        this.exprAllowed = true;
      };
      types$1.parenL.updateContext = function(prevType) {
        var statementParens = prevType === types$1._if || prevType === types$1._for || prevType === types$1._with || prevType === types$1._while;
        this.context.push(statementParens ? types.p_stat : types.p_expr);
        this.exprAllowed = true;
      };
      types$1.incDec.updateContext = function() {
      };
      types$1._function.updateContext = types$1._class.updateContext = function(prevType) {
        if (prevType.beforeExpr && prevType !== types$1._else && !(prevType === types$1.semi && this.curContext() !== types.p_stat) && !(prevType === types$1._return && lineBreak.test(this.input.slice(this.lastTokEnd, this.start))) && !((prevType === types$1.colon || prevType === types$1.braceL) && this.curContext() === types.b_stat)) {
          this.context.push(types.f_expr);
        } else {
          this.context.push(types.f_stat);
        }
        this.exprAllowed = false;
      };
      types$1.colon.updateContext = function() {
        if (this.curContext().token === "function") {
          this.context.pop();
        }
        this.exprAllowed = true;
      };
      types$1.backQuote.updateContext = function() {
        if (this.curContext() === types.q_tmpl) {
          this.context.pop();
        } else {
          this.context.push(types.q_tmpl);
        }
        this.exprAllowed = false;
      };
      types$1.star.updateContext = function(prevType) {
        if (prevType === types$1._function) {
          var index2 = this.context.length - 1;
          if (this.context[index2] === types.f_expr) {
            this.context[index2] = types.f_expr_gen;
          } else {
            this.context[index2] = types.f_gen;
          }
        }
        this.exprAllowed = true;
      };
      types$1.name.updateContext = function(prevType) {
        var allowed = false;
        if (this.options.ecmaVersion >= 6 && prevType !== types$1.dot) {
          if (this.value === "of" && !this.exprAllowed || this.value === "yield" && this.inGeneratorContext()) {
            allowed = true;
          }
        }
        this.exprAllowed = allowed;
      };
      pp$5 = Parser.prototype;
      pp$5.checkPropClash = function(prop, propHash, refDestructuringErrors) {
        if (this.options.ecmaVersion >= 9 && prop.type === "SpreadElement") {
          return;
        }
        if (this.options.ecmaVersion >= 6 && (prop.computed || prop.method || prop.shorthand)) {
          return;
        }
        var key = prop.key;
        var name2;
        switch (key.type) {
          case "Identifier":
            name2 = key.name;
            break;
          case "Literal":
            name2 = String(key.value);
            break;
          default:
            return;
        }
        var kind = prop.kind;
        if (this.options.ecmaVersion >= 6) {
          if (name2 === "__proto__" && kind === "init") {
            if (propHash.proto) {
              if (refDestructuringErrors) {
                if (refDestructuringErrors.doubleProto < 0) {
                  refDestructuringErrors.doubleProto = key.start;
                }
              } else {
                this.raiseRecoverable(key.start, "Redefinition of __proto__ property");
              }
            }
            propHash.proto = true;
          }
          return;
        }
        name2 = "$" + name2;
        var other = propHash[name2];
        if (other) {
          var redefinition;
          if (kind === "init") {
            redefinition = this.strict && other.init || other.get || other.set;
          } else {
            redefinition = other.init || other[kind];
          }
          if (redefinition) {
            this.raiseRecoverable(key.start, "Redefinition of property");
          }
        } else {
          other = propHash[name2] = {
            init: false,
            get: false,
            set: false
          };
        }
        other[kind] = true;
      };
      pp$5.parseExpression = function(forInit, refDestructuringErrors) {
        var startPos = this.start, startLoc = this.startLoc;
        var expr = this.parseMaybeAssign(forInit, refDestructuringErrors);
        if (this.type === types$1.comma) {
          var node2 = this.startNodeAt(startPos, startLoc);
          node2.expressions = [expr];
          while (this.eat(types$1.comma)) {
            node2.expressions.push(this.parseMaybeAssign(forInit, refDestructuringErrors));
          }
          return this.finishNode(node2, "SequenceExpression");
        }
        return expr;
      };
      pp$5.parseMaybeAssign = function(forInit, refDestructuringErrors, afterLeftParse) {
        if (this.isContextual("yield")) {
          if (this.inGenerator) {
            return this.parseYield(forInit);
          } else {
            this.exprAllowed = false;
          }
        }
        var ownDestructuringErrors = false, oldParenAssign = -1, oldTrailingComma = -1, oldDoubleProto = -1;
        if (refDestructuringErrors) {
          oldParenAssign = refDestructuringErrors.parenthesizedAssign;
          oldTrailingComma = refDestructuringErrors.trailingComma;
          oldDoubleProto = refDestructuringErrors.doubleProto;
          refDestructuringErrors.parenthesizedAssign = refDestructuringErrors.trailingComma = -1;
        } else {
          refDestructuringErrors = new DestructuringErrors();
          ownDestructuringErrors = true;
        }
        var startPos = this.start, startLoc = this.startLoc;
        if (this.type === types$1.parenL || this.type === types$1.name) {
          this.potentialArrowAt = this.start;
          this.potentialArrowInForAwait = forInit === "await";
        }
        var left = this.parseMaybeConditional(forInit, refDestructuringErrors);
        if (afterLeftParse) {
          left = afterLeftParse.call(this, left, startPos, startLoc);
        }
        if (this.type.isAssign) {
          var node2 = this.startNodeAt(startPos, startLoc);
          node2.operator = this.value;
          if (this.type === types$1.eq) {
            left = this.toAssignable(left, false, refDestructuringErrors);
          }
          if (!ownDestructuringErrors) {
            refDestructuringErrors.parenthesizedAssign = refDestructuringErrors.trailingComma = refDestructuringErrors.doubleProto = -1;
          }
          if (refDestructuringErrors.shorthandAssign >= left.start) {
            refDestructuringErrors.shorthandAssign = -1;
          }
          if (this.type === types$1.eq) {
            this.checkLValPattern(left);
          } else {
            this.checkLValSimple(left);
          }
          node2.left = left;
          this.next();
          node2.right = this.parseMaybeAssign(forInit);
          if (oldDoubleProto > -1) {
            refDestructuringErrors.doubleProto = oldDoubleProto;
          }
          return this.finishNode(node2, "AssignmentExpression");
        } else {
          if (ownDestructuringErrors) {
            this.checkExpressionErrors(refDestructuringErrors, true);
          }
        }
        if (oldParenAssign > -1) {
          refDestructuringErrors.parenthesizedAssign = oldParenAssign;
        }
        if (oldTrailingComma > -1) {
          refDestructuringErrors.trailingComma = oldTrailingComma;
        }
        return left;
      };
      pp$5.parseMaybeConditional = function(forInit, refDestructuringErrors) {
        var startPos = this.start, startLoc = this.startLoc;
        var expr = this.parseExprOps(forInit, refDestructuringErrors);
        if (this.checkExpressionErrors(refDestructuringErrors)) {
          return expr;
        }
        if (this.eat(types$1.question)) {
          var node2 = this.startNodeAt(startPos, startLoc);
          node2.test = expr;
          node2.consequent = this.parseMaybeAssign();
          this.expect(types$1.colon);
          node2.alternate = this.parseMaybeAssign(forInit);
          return this.finishNode(node2, "ConditionalExpression");
        }
        return expr;
      };
      pp$5.parseExprOps = function(forInit, refDestructuringErrors) {
        var startPos = this.start, startLoc = this.startLoc;
        var expr = this.parseMaybeUnary(refDestructuringErrors, false, false, forInit);
        if (this.checkExpressionErrors(refDestructuringErrors)) {
          return expr;
        }
        return expr.start === startPos && expr.type === "ArrowFunctionExpression" ? expr : this.parseExprOp(expr, startPos, startLoc, -1, forInit);
      };
      pp$5.parseExprOp = function(left, leftStartPos, leftStartLoc, minPrec, forInit) {
        var prec = this.type.binop;
        if (prec != null && (!forInit || this.type !== types$1._in)) {
          if (prec > minPrec) {
            var logical = this.type === types$1.logicalOR || this.type === types$1.logicalAND;
            var coalesce = this.type === types$1.coalesce;
            if (coalesce) {
              prec = types$1.logicalAND.binop;
            }
            var op = this.value;
            this.next();
            var startPos = this.start, startLoc = this.startLoc;
            var right = this.parseExprOp(this.parseMaybeUnary(null, false, false, forInit), startPos, startLoc, prec, forInit);
            var node2 = this.buildBinary(leftStartPos, leftStartLoc, left, right, op, logical || coalesce);
            if (logical && this.type === types$1.coalesce || coalesce && (this.type === types$1.logicalOR || this.type === types$1.logicalAND)) {
              this.raiseRecoverable(this.start, "Logical expressions and coalesce expressions cannot be mixed. Wrap either by parentheses");
            }
            return this.parseExprOp(node2, leftStartPos, leftStartLoc, minPrec, forInit);
          }
        }
        return left;
      };
      pp$5.buildBinary = function(startPos, startLoc, left, right, op, logical) {
        if (right.type === "PrivateIdentifier") {
          this.raise(right.start, "Private identifier can only be left side of binary expression");
        }
        var node2 = this.startNodeAt(startPos, startLoc);
        node2.left = left;
        node2.operator = op;
        node2.right = right;
        return this.finishNode(node2, logical ? "LogicalExpression" : "BinaryExpression");
      };
      pp$5.parseMaybeUnary = function(refDestructuringErrors, sawUnary, incDec, forInit) {
        var startPos = this.start, startLoc = this.startLoc, expr;
        if (this.isContextual("await") && this.canAwait) {
          expr = this.parseAwait(forInit);
          sawUnary = true;
        } else if (this.type.prefix) {
          var node2 = this.startNode(), update2 = this.type === types$1.incDec;
          node2.operator = this.value;
          node2.prefix = true;
          this.next();
          node2.argument = this.parseMaybeUnary(null, true, update2, forInit);
          this.checkExpressionErrors(refDestructuringErrors, true);
          if (update2) {
            this.checkLValSimple(node2.argument);
          } else if (this.strict && node2.operator === "delete" && isLocalVariableAccess(node2.argument)) {
            this.raiseRecoverable(node2.start, "Deleting local variable in strict mode");
          } else if (node2.operator === "delete" && isPrivateFieldAccess(node2.argument)) {
            this.raiseRecoverable(node2.start, "Private fields can not be deleted");
          } else {
            sawUnary = true;
          }
          expr = this.finishNode(node2, update2 ? "UpdateExpression" : "UnaryExpression");
        } else if (!sawUnary && this.type === types$1.privateId) {
          if ((forInit || this.privateNameStack.length === 0) && this.options.checkPrivateFields) {
            this.unexpected();
          }
          expr = this.parsePrivateIdent();
          if (this.type !== types$1._in) {
            this.unexpected();
          }
        } else {
          expr = this.parseExprSubscripts(refDestructuringErrors, forInit);
          if (this.checkExpressionErrors(refDestructuringErrors)) {
            return expr;
          }
          while (this.type.postfix && !this.canInsertSemicolon()) {
            var node$1 = this.startNodeAt(startPos, startLoc);
            node$1.operator = this.value;
            node$1.prefix = false;
            node$1.argument = expr;
            this.checkLValSimple(expr);
            this.next();
            expr = this.finishNode(node$1, "UpdateExpression");
          }
        }
        if (!incDec && this.eat(types$1.starstar)) {
          if (sawUnary) {
            this.unexpected(this.lastTokStart);
          } else {
            return this.buildBinary(startPos, startLoc, expr, this.parseMaybeUnary(null, false, false, forInit), "**", false);
          }
        } else {
          return expr;
        }
      };
      pp$5.parseExprSubscripts = function(refDestructuringErrors, forInit) {
        var startPos = this.start, startLoc = this.startLoc;
        var expr = this.parseExprAtom(refDestructuringErrors, forInit);
        if (expr.type === "ArrowFunctionExpression" && this.input.slice(this.lastTokStart, this.lastTokEnd) !== ")") {
          return expr;
        }
        var result = this.parseSubscripts(expr, startPos, startLoc, false, forInit);
        if (refDestructuringErrors && result.type === "MemberExpression") {
          if (refDestructuringErrors.parenthesizedAssign >= result.start) {
            refDestructuringErrors.parenthesizedAssign = -1;
          }
          if (refDestructuringErrors.parenthesizedBind >= result.start) {
            refDestructuringErrors.parenthesizedBind = -1;
          }
          if (refDestructuringErrors.trailingComma >= result.start) {
            refDestructuringErrors.trailingComma = -1;
          }
        }
        return result;
      };
      pp$5.parseSubscripts = function(base, startPos, startLoc, noCalls, forInit) {
        var maybeAsyncArrow = this.options.ecmaVersion >= 8 && base.type === "Identifier" && base.name === "async" && this.lastTokEnd === base.end && !this.canInsertSemicolon() && base.end - base.start === 5 && this.potentialArrowAt === base.start;
        var optionalChained = false;
        while (true) {
          var element = this.parseSubscript(base, startPos, startLoc, noCalls, maybeAsyncArrow, optionalChained, forInit);
          if (element.optional) {
            optionalChained = true;
          }
          if (element === base || element.type === "ArrowFunctionExpression") {
            if (optionalChained) {
              var chainNode = this.startNodeAt(startPos, startLoc);
              chainNode.expression = element;
              element = this.finishNode(chainNode, "ChainExpression");
            }
            return element;
          }
          base = element;
        }
      };
      pp$5.shouldParseAsyncArrow = function() {
        return !this.canInsertSemicolon() && this.eat(types$1.arrow);
      };
      pp$5.parseSubscriptAsyncArrow = function(startPos, startLoc, exprList, forInit) {
        return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), exprList, true, forInit);
      };
      pp$5.parseSubscript = function(base, startPos, startLoc, noCalls, maybeAsyncArrow, optionalChained, forInit) {
        var optionalSupported = this.options.ecmaVersion >= 11;
        var optional = optionalSupported && this.eat(types$1.questionDot);
        if (noCalls && optional) {
          this.raise(this.lastTokStart, "Optional chaining cannot appear in the callee of new expressions");
        }
        var computed = this.eat(types$1.bracketL);
        if (computed || optional && this.type !== types$1.parenL && this.type !== types$1.backQuote || this.eat(types$1.dot)) {
          var node2 = this.startNodeAt(startPos, startLoc);
          node2.object = base;
          if (computed) {
            node2.property = this.parseExpression();
            this.expect(types$1.bracketR);
          } else if (this.type === types$1.privateId && base.type !== "Super") {
            node2.property = this.parsePrivateIdent();
          } else {
            node2.property = this.parseIdent(this.options.allowReserved !== "never");
          }
          node2.computed = !!computed;
          if (optionalSupported) {
            node2.optional = optional;
          }
          base = this.finishNode(node2, "MemberExpression");
        } else if (!noCalls && this.eat(types$1.parenL)) {
          var refDestructuringErrors = new DestructuringErrors(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
          this.yieldPos = 0;
          this.awaitPos = 0;
          this.awaitIdentPos = 0;
          var exprList = this.parseExprList(types$1.parenR, this.options.ecmaVersion >= 8, false, refDestructuringErrors);
          if (maybeAsyncArrow && !optional && this.shouldParseAsyncArrow()) {
            this.checkPatternErrors(refDestructuringErrors, false);
            this.checkYieldAwaitInDefaultParams();
            if (this.awaitIdentPos > 0) {
              this.raise(this.awaitIdentPos, "Cannot use 'await' as identifier inside an async function");
            }
            this.yieldPos = oldYieldPos;
            this.awaitPos = oldAwaitPos;
            this.awaitIdentPos = oldAwaitIdentPos;
            return this.parseSubscriptAsyncArrow(startPos, startLoc, exprList, forInit);
          }
          this.checkExpressionErrors(refDestructuringErrors, true);
          this.yieldPos = oldYieldPos || this.yieldPos;
          this.awaitPos = oldAwaitPos || this.awaitPos;
          this.awaitIdentPos = oldAwaitIdentPos || this.awaitIdentPos;
          var node$1 = this.startNodeAt(startPos, startLoc);
          node$1.callee = base;
          node$1.arguments = exprList;
          if (optionalSupported) {
            node$1.optional = optional;
          }
          base = this.finishNode(node$1, "CallExpression");
        } else if (this.type === types$1.backQuote) {
          if (optional || optionalChained) {
            this.raise(this.start, "Optional chaining cannot appear in the tag of tagged template expressions");
          }
          var node$2 = this.startNodeAt(startPos, startLoc);
          node$2.tag = base;
          node$2.quasi = this.parseTemplate({ isTagged: true });
          base = this.finishNode(node$2, "TaggedTemplateExpression");
        }
        return base;
      };
      pp$5.parseExprAtom = function(refDestructuringErrors, forInit, forNew) {
        if (this.type === types$1.slash) {
          this.readRegexp();
        }
        var node2, canBeArrow = this.potentialArrowAt === this.start;
        switch (this.type) {
          case types$1._super:
            if (!this.allowSuper) {
              this.raise(this.start, "'super' keyword outside a method");
            }
            node2 = this.startNode();
            this.next();
            if (this.type === types$1.parenL && !this.allowDirectSuper) {
              this.raise(node2.start, "super() call outside constructor of a subclass");
            }
            if (this.type !== types$1.dot && this.type !== types$1.bracketL && this.type !== types$1.parenL) {
              this.unexpected();
            }
            return this.finishNode(node2, "Super");
          case types$1._this:
            node2 = this.startNode();
            this.next();
            return this.finishNode(node2, "ThisExpression");
          case types$1.name:
            var startPos = this.start, startLoc = this.startLoc, containsEsc = this.containsEsc;
            var id = this.parseIdent(false);
            if (this.options.ecmaVersion >= 8 && !containsEsc && id.name === "async" && !this.canInsertSemicolon() && this.eat(types$1._function)) {
              this.overrideContext(types.f_expr);
              return this.parseFunction(this.startNodeAt(startPos, startLoc), 0, false, true, forInit);
            }
            if (canBeArrow && !this.canInsertSemicolon()) {
              if (this.eat(types$1.arrow)) {
                return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), [id], false, forInit);
              }
              if (this.options.ecmaVersion >= 8 && id.name === "async" && this.type === types$1.name && !containsEsc && (!this.potentialArrowInForAwait || this.value !== "of" || this.containsEsc)) {
                id = this.parseIdent(false);
                if (this.canInsertSemicolon() || !this.eat(types$1.arrow)) {
                  this.unexpected();
                }
                return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), [id], true, forInit);
              }
            }
            return id;
          case types$1.regexp:
            var value = this.value;
            node2 = this.parseLiteral(value.value);
            node2.regex = { pattern: value.pattern, flags: value.flags };
            return node2;
          case types$1.num:
          case types$1.string:
            return this.parseLiteral(this.value);
          case types$1._null:
          case types$1._true:
          case types$1._false:
            node2 = this.startNode();
            node2.value = this.type === types$1._null ? null : this.type === types$1._true;
            node2.raw = this.type.keyword;
            this.next();
            return this.finishNode(node2, "Literal");
          case types$1.parenL:
            var start2 = this.start, expr = this.parseParenAndDistinguishExpression(canBeArrow, forInit);
            if (refDestructuringErrors) {
              if (refDestructuringErrors.parenthesizedAssign < 0 && !this.isSimpleAssignTarget(expr)) {
                refDestructuringErrors.parenthesizedAssign = start2;
              }
              if (refDestructuringErrors.parenthesizedBind < 0) {
                refDestructuringErrors.parenthesizedBind = start2;
              }
            }
            return expr;
          case types$1.bracketL:
            node2 = this.startNode();
            this.next();
            node2.elements = this.parseExprList(types$1.bracketR, true, true, refDestructuringErrors);
            return this.finishNode(node2, "ArrayExpression");
          case types$1.braceL:
            this.overrideContext(types.b_expr);
            return this.parseObj(false, refDestructuringErrors);
          case types$1._function:
            node2 = this.startNode();
            this.next();
            return this.parseFunction(node2, 0);
          case types$1._class:
            return this.parseClass(this.startNode(), false);
          case types$1._new:
            return this.parseNew();
          case types$1.backQuote:
            return this.parseTemplate();
          case types$1._import:
            if (this.options.ecmaVersion >= 11) {
              return this.parseExprImport(forNew);
            } else {
              return this.unexpected();
            }
          default:
            return this.parseExprAtomDefault();
        }
      };
      pp$5.parseExprAtomDefault = function() {
        this.unexpected();
      };
      pp$5.parseExprImport = function(forNew) {
        var node2 = this.startNode();
        if (this.containsEsc) {
          this.raiseRecoverable(this.start, "Escape sequence in keyword import");
        }
        this.next();
        if (this.type === types$1.parenL && !forNew) {
          return this.parseDynamicImport(node2);
        } else if (this.type === types$1.dot) {
          var meta = this.startNodeAt(node2.start, node2.loc && node2.loc.start);
          meta.name = "import";
          node2.meta = this.finishNode(meta, "Identifier");
          return this.parseImportMeta(node2);
        } else {
          this.unexpected();
        }
      };
      pp$5.parseDynamicImport = function(node2) {
        this.next();
        node2.source = this.parseMaybeAssign();
        if (this.options.ecmaVersion >= 16) {
          if (!this.eat(types$1.parenR)) {
            this.expect(types$1.comma);
            if (!this.afterTrailingComma(types$1.parenR)) {
              node2.options = this.parseMaybeAssign();
              if (!this.eat(types$1.parenR)) {
                this.expect(types$1.comma);
                if (!this.afterTrailingComma(types$1.parenR)) {
                  this.unexpected();
                }
              }
            } else {
              node2.options = null;
            }
          } else {
            node2.options = null;
          }
        } else {
          if (!this.eat(types$1.parenR)) {
            var errorPos = this.start;
            if (this.eat(types$1.comma) && this.eat(types$1.parenR)) {
              this.raiseRecoverable(errorPos, "Trailing comma is not allowed in import()");
            } else {
              this.unexpected(errorPos);
            }
          }
        }
        return this.finishNode(node2, "ImportExpression");
      };
      pp$5.parseImportMeta = function(node2) {
        this.next();
        var containsEsc = this.containsEsc;
        node2.property = this.parseIdent(true);
        if (node2.property.name !== "meta") {
          this.raiseRecoverable(node2.property.start, "The only valid meta property for import is 'import.meta'");
        }
        if (containsEsc) {
          this.raiseRecoverable(node2.start, "'import.meta' must not contain escaped characters");
        }
        if (this.options.sourceType !== "module" && !this.options.allowImportExportEverywhere) {
          this.raiseRecoverable(node2.start, "Cannot use 'import.meta' outside a module");
        }
        return this.finishNode(node2, "MetaProperty");
      };
      pp$5.parseLiteral = function(value) {
        var node2 = this.startNode();
        node2.value = value;
        node2.raw = this.input.slice(this.start, this.end);
        if (node2.raw.charCodeAt(node2.raw.length - 1) === 110) {
          node2.bigint = node2.value != null ? node2.value.toString() : node2.raw.slice(0, -1).replace(/_/g, "");
        }
        this.next();
        return this.finishNode(node2, "Literal");
      };
      pp$5.parseParenExpression = function() {
        this.expect(types$1.parenL);
        var val = this.parseExpression();
        this.expect(types$1.parenR);
        return val;
      };
      pp$5.shouldParseArrow = function(exprList) {
        return !this.canInsertSemicolon();
      };
      pp$5.parseParenAndDistinguishExpression = function(canBeArrow, forInit) {
        var startPos = this.start, startLoc = this.startLoc, val, allowTrailingComma = this.options.ecmaVersion >= 8;
        if (this.options.ecmaVersion >= 6) {
          this.next();
          var innerStartPos = this.start, innerStartLoc = this.startLoc;
          var exprList = [], first = true, lastIsComma = false;
          var refDestructuringErrors = new DestructuringErrors(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, spreadStart;
          this.yieldPos = 0;
          this.awaitPos = 0;
          while (this.type !== types$1.parenR) {
            first ? first = false : this.expect(types$1.comma);
            if (allowTrailingComma && this.afterTrailingComma(types$1.parenR, true)) {
              lastIsComma = true;
              break;
            } else if (this.type === types$1.ellipsis) {
              spreadStart = this.start;
              exprList.push(this.parseParenItem(this.parseRestBinding()));
              if (this.type === types$1.comma) {
                this.raiseRecoverable(
                  this.start,
                  "Comma is not permitted after the rest element"
                );
              }
              break;
            } else {
              exprList.push(this.parseMaybeAssign(false, refDestructuringErrors, this.parseParenItem));
            }
          }
          var innerEndPos = this.lastTokEnd, innerEndLoc = this.lastTokEndLoc;
          this.expect(types$1.parenR);
          if (canBeArrow && this.shouldParseArrow(exprList) && this.eat(types$1.arrow)) {
            this.checkPatternErrors(refDestructuringErrors, false);
            this.checkYieldAwaitInDefaultParams();
            this.yieldPos = oldYieldPos;
            this.awaitPos = oldAwaitPos;
            return this.parseParenArrowList(startPos, startLoc, exprList, forInit);
          }
          if (!exprList.length || lastIsComma) {
            this.unexpected(this.lastTokStart);
          }
          if (spreadStart) {
            this.unexpected(spreadStart);
          }
          this.checkExpressionErrors(refDestructuringErrors, true);
          this.yieldPos = oldYieldPos || this.yieldPos;
          this.awaitPos = oldAwaitPos || this.awaitPos;
          if (exprList.length > 1) {
            val = this.startNodeAt(innerStartPos, innerStartLoc);
            val.expressions = exprList;
            this.finishNodeAt(val, "SequenceExpression", innerEndPos, innerEndLoc);
          } else {
            val = exprList[0];
          }
        } else {
          val = this.parseParenExpression();
        }
        if (this.options.preserveParens) {
          var par = this.startNodeAt(startPos, startLoc);
          par.expression = val;
          return this.finishNode(par, "ParenthesizedExpression");
        } else {
          return val;
        }
      };
      pp$5.parseParenItem = function(item) {
        return item;
      };
      pp$5.parseParenArrowList = function(startPos, startLoc, exprList, forInit) {
        return this.parseArrowExpression(this.startNodeAt(startPos, startLoc), exprList, false, forInit);
      };
      empty = [];
      pp$5.parseNew = function() {
        if (this.containsEsc) {
          this.raiseRecoverable(this.start, "Escape sequence in keyword new");
        }
        var node2 = this.startNode();
        this.next();
        if (this.options.ecmaVersion >= 6 && this.type === types$1.dot) {
          var meta = this.startNodeAt(node2.start, node2.loc && node2.loc.start);
          meta.name = "new";
          node2.meta = this.finishNode(meta, "Identifier");
          this.next();
          var containsEsc = this.containsEsc;
          node2.property = this.parseIdent(true);
          if (node2.property.name !== "target") {
            this.raiseRecoverable(node2.property.start, "The only valid meta property for new is 'new.target'");
          }
          if (containsEsc) {
            this.raiseRecoverable(node2.start, "'new.target' must not contain escaped characters");
          }
          if (!this.allowNewDotTarget) {
            this.raiseRecoverable(node2.start, "'new.target' can only be used in functions and class static block");
          }
          return this.finishNode(node2, "MetaProperty");
        }
        var startPos = this.start, startLoc = this.startLoc;
        node2.callee = this.parseSubscripts(this.parseExprAtom(null, false, true), startPos, startLoc, true, false);
        if (this.eat(types$1.parenL)) {
          node2.arguments = this.parseExprList(types$1.parenR, this.options.ecmaVersion >= 8, false);
        } else {
          node2.arguments = empty;
        }
        return this.finishNode(node2, "NewExpression");
      };
      pp$5.parseTemplateElement = function(ref2) {
        var isTagged = ref2.isTagged;
        var elem = this.startNode();
        if (this.type === types$1.invalidTemplate) {
          if (!isTagged) {
            this.raiseRecoverable(this.start, "Bad escape sequence in untagged template literal");
          }
          elem.value = {
            raw: this.value.replace(/\r\n?/g, "\n"),
            cooked: null
          };
        } else {
          elem.value = {
            raw: this.input.slice(this.start, this.end).replace(/\r\n?/g, "\n"),
            cooked: this.value
          };
        }
        this.next();
        elem.tail = this.type === types$1.backQuote;
        return this.finishNode(elem, "TemplateElement");
      };
      pp$5.parseTemplate = function(ref2) {
        if (ref2 === void 0) ref2 = {};
        var isTagged = ref2.isTagged;
        if (isTagged === void 0) isTagged = false;
        var node2 = this.startNode();
        this.next();
        node2.expressions = [];
        var curElt = this.parseTemplateElement({ isTagged });
        node2.quasis = [curElt];
        while (!curElt.tail) {
          if (this.type === types$1.eof) {
            this.raise(this.pos, "Unterminated template literal");
          }
          this.expect(types$1.dollarBraceL);
          node2.expressions.push(this.parseExpression());
          this.expect(types$1.braceR);
          node2.quasis.push(curElt = this.parseTemplateElement({ isTagged }));
        }
        this.next();
        return this.finishNode(node2, "TemplateLiteral");
      };
      pp$5.isAsyncProp = function(prop) {
        return !prop.computed && prop.key.type === "Identifier" && prop.key.name === "async" && (this.type === types$1.name || this.type === types$1.num || this.type === types$1.string || this.type === types$1.bracketL || this.type.keyword || this.options.ecmaVersion >= 9 && this.type === types$1.star) && !lineBreak.test(this.input.slice(this.lastTokEnd, this.start));
      };
      pp$5.parseObj = function(isPattern, refDestructuringErrors) {
        var node2 = this.startNode(), first = true, propHash = {};
        node2.properties = [];
        this.next();
        while (!this.eat(types$1.braceR)) {
          if (!first) {
            this.expect(types$1.comma);
            if (this.options.ecmaVersion >= 5 && this.afterTrailingComma(types$1.braceR)) {
              break;
            }
          } else {
            first = false;
          }
          var prop = this.parseProperty(isPattern, refDestructuringErrors);
          if (!isPattern) {
            this.checkPropClash(prop, propHash, refDestructuringErrors);
          }
          node2.properties.push(prop);
        }
        return this.finishNode(node2, isPattern ? "ObjectPattern" : "ObjectExpression");
      };
      pp$5.parseProperty = function(isPattern, refDestructuringErrors) {
        var prop = this.startNode(), isGenerator, isAsync, startPos, startLoc;
        if (this.options.ecmaVersion >= 9 && this.eat(types$1.ellipsis)) {
          if (isPattern) {
            prop.argument = this.parseIdent(false);
            if (this.type === types$1.comma) {
              this.raiseRecoverable(this.start, "Comma is not permitted after the rest element");
            }
            return this.finishNode(prop, "RestElement");
          }
          prop.argument = this.parseMaybeAssign(false, refDestructuringErrors);
          if (this.type === types$1.comma && refDestructuringErrors && refDestructuringErrors.trailingComma < 0) {
            refDestructuringErrors.trailingComma = this.start;
          }
          return this.finishNode(prop, "SpreadElement");
        }
        if (this.options.ecmaVersion >= 6) {
          prop.method = false;
          prop.shorthand = false;
          if (isPattern || refDestructuringErrors) {
            startPos = this.start;
            startLoc = this.startLoc;
          }
          if (!isPattern) {
            isGenerator = this.eat(types$1.star);
          }
        }
        var containsEsc = this.containsEsc;
        this.parsePropertyName(prop);
        if (!isPattern && !containsEsc && this.options.ecmaVersion >= 8 && !isGenerator && this.isAsyncProp(prop)) {
          isAsync = true;
          isGenerator = this.options.ecmaVersion >= 9 && this.eat(types$1.star);
          this.parsePropertyName(prop);
        } else {
          isAsync = false;
        }
        this.parsePropertyValue(prop, isPattern, isGenerator, isAsync, startPos, startLoc, refDestructuringErrors, containsEsc);
        return this.finishNode(prop, "Property");
      };
      pp$5.parseGetterSetter = function(prop) {
        var kind = prop.key.name;
        this.parsePropertyName(prop);
        prop.value = this.parseMethod(false);
        prop.kind = kind;
        var paramCount = prop.kind === "get" ? 0 : 1;
        if (prop.value.params.length !== paramCount) {
          var start2 = prop.value.start;
          if (prop.kind === "get") {
            this.raiseRecoverable(start2, "getter should have no params");
          } else {
            this.raiseRecoverable(start2, "setter should have exactly one param");
          }
        } else {
          if (prop.kind === "set" && prop.value.params[0].type === "RestElement") {
            this.raiseRecoverable(prop.value.params[0].start, "Setter cannot use rest params");
          }
        }
      };
      pp$5.parsePropertyValue = function(prop, isPattern, isGenerator, isAsync, startPos, startLoc, refDestructuringErrors, containsEsc) {
        if ((isGenerator || isAsync) && this.type === types$1.colon) {
          this.unexpected();
        }
        if (this.eat(types$1.colon)) {
          prop.value = isPattern ? this.parseMaybeDefault(this.start, this.startLoc) : this.parseMaybeAssign(false, refDestructuringErrors);
          prop.kind = "init";
        } else if (this.options.ecmaVersion >= 6 && this.type === types$1.parenL) {
          if (isPattern) {
            this.unexpected();
          }
          prop.method = true;
          prop.value = this.parseMethod(isGenerator, isAsync);
          prop.kind = "init";
        } else if (!isPattern && !containsEsc && this.options.ecmaVersion >= 5 && !prop.computed && prop.key.type === "Identifier" && (prop.key.name === "get" || prop.key.name === "set") && (this.type !== types$1.comma && this.type !== types$1.braceR && this.type !== types$1.eq)) {
          if (isGenerator || isAsync) {
            this.unexpected();
          }
          this.parseGetterSetter(prop);
        } else if (this.options.ecmaVersion >= 6 && !prop.computed && prop.key.type === "Identifier") {
          if (isGenerator || isAsync) {
            this.unexpected();
          }
          this.checkUnreserved(prop.key);
          if (prop.key.name === "await" && !this.awaitIdentPos) {
            this.awaitIdentPos = startPos;
          }
          if (isPattern) {
            prop.value = this.parseMaybeDefault(startPos, startLoc, this.copyNode(prop.key));
          } else if (this.type === types$1.eq && refDestructuringErrors) {
            if (refDestructuringErrors.shorthandAssign < 0) {
              refDestructuringErrors.shorthandAssign = this.start;
            }
            prop.value = this.parseMaybeDefault(startPos, startLoc, this.copyNode(prop.key));
          } else {
            prop.value = this.copyNode(prop.key);
          }
          prop.kind = "init";
          prop.shorthand = true;
        } else {
          this.unexpected();
        }
      };
      pp$5.parsePropertyName = function(prop) {
        if (this.options.ecmaVersion >= 6) {
          if (this.eat(types$1.bracketL)) {
            prop.computed = true;
            prop.key = this.parseMaybeAssign();
            this.expect(types$1.bracketR);
            return prop.key;
          } else {
            prop.computed = false;
          }
        }
        return prop.key = this.type === types$1.num || this.type === types$1.string ? this.parseExprAtom() : this.parseIdent(this.options.allowReserved !== "never");
      };
      pp$5.initFunction = function(node2) {
        node2.id = null;
        if (this.options.ecmaVersion >= 6) {
          node2.generator = node2.expression = false;
        }
        if (this.options.ecmaVersion >= 8) {
          node2.async = false;
        }
      };
      pp$5.parseMethod = function(isGenerator, isAsync, allowDirectSuper) {
        var node2 = this.startNode(), oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
        this.initFunction(node2);
        if (this.options.ecmaVersion >= 6) {
          node2.generator = isGenerator;
        }
        if (this.options.ecmaVersion >= 8) {
          node2.async = !!isAsync;
        }
        this.yieldPos = 0;
        this.awaitPos = 0;
        this.awaitIdentPos = 0;
        this.enterScope(functionFlags(isAsync, node2.generator) | SCOPE_SUPER | (allowDirectSuper ? SCOPE_DIRECT_SUPER : 0));
        this.expect(types$1.parenL);
        node2.params = this.parseBindingList(types$1.parenR, false, this.options.ecmaVersion >= 8);
        this.checkYieldAwaitInDefaultParams();
        this.parseFunctionBody(node2, false, true, false);
        this.yieldPos = oldYieldPos;
        this.awaitPos = oldAwaitPos;
        this.awaitIdentPos = oldAwaitIdentPos;
        return this.finishNode(node2, "FunctionExpression");
      };
      pp$5.parseArrowExpression = function(node2, params, isAsync, forInit) {
        var oldYieldPos = this.yieldPos, oldAwaitPos = this.awaitPos, oldAwaitIdentPos = this.awaitIdentPos;
        this.enterScope(functionFlags(isAsync, false) | SCOPE_ARROW);
        this.initFunction(node2);
        if (this.options.ecmaVersion >= 8) {
          node2.async = !!isAsync;
        }
        this.yieldPos = 0;
        this.awaitPos = 0;
        this.awaitIdentPos = 0;
        node2.params = this.toAssignableList(params, true);
        this.parseFunctionBody(node2, true, false, forInit);
        this.yieldPos = oldYieldPos;
        this.awaitPos = oldAwaitPos;
        this.awaitIdentPos = oldAwaitIdentPos;
        return this.finishNode(node2, "ArrowFunctionExpression");
      };
      pp$5.parseFunctionBody = function(node2, isArrowFunction, isMethod, forInit) {
        var isExpression = isArrowFunction && this.type !== types$1.braceL;
        var oldStrict = this.strict, useStrict = false;
        if (isExpression) {
          node2.body = this.parseMaybeAssign(forInit);
          node2.expression = true;
          this.checkParams(node2, false);
        } else {
          var nonSimple = this.options.ecmaVersion >= 7 && !this.isSimpleParamList(node2.params);
          if (!oldStrict || nonSimple) {
            useStrict = this.strictDirective(this.end);
            if (useStrict && nonSimple) {
              this.raiseRecoverable(node2.start, "Illegal 'use strict' directive in function with non-simple parameter list");
            }
          }
          var oldLabels = this.labels;
          this.labels = [];
          if (useStrict) {
            this.strict = true;
          }
          this.checkParams(node2, !oldStrict && !useStrict && !isArrowFunction && !isMethod && this.isSimpleParamList(node2.params));
          if (this.strict && node2.id) {
            this.checkLValSimple(node2.id, BIND_OUTSIDE);
          }
          node2.body = this.parseBlock(false, void 0, useStrict && !oldStrict);
          node2.expression = false;
          this.adaptDirectivePrologue(node2.body.body);
          this.labels = oldLabels;
        }
        this.exitScope();
      };
      pp$5.isSimpleParamList = function(params) {
        for (var i = 0, list4 = params; i < list4.length; i += 1) {
          var param = list4[i];
          if (param.type !== "Identifier") {
            return false;
          }
        }
        return true;
      };
      pp$5.checkParams = function(node2, allowDuplicates) {
        var nameHash = /* @__PURE__ */ Object.create(null);
        for (var i = 0, list4 = node2.params; i < list4.length; i += 1) {
          var param = list4[i];
          this.checkLValInnerPattern(param, BIND_VAR, allowDuplicates ? null : nameHash);
        }
      };
      pp$5.parseExprList = function(close, allowTrailingComma, allowEmpty, refDestructuringErrors) {
        var elts = [], first = true;
        while (!this.eat(close)) {
          if (!first) {
            this.expect(types$1.comma);
            if (allowTrailingComma && this.afterTrailingComma(close)) {
              break;
            }
          } else {
            first = false;
          }
          var elt = void 0;
          if (allowEmpty && this.type === types$1.comma) {
            elt = null;
          } else if (this.type === types$1.ellipsis) {
            elt = this.parseSpread(refDestructuringErrors);
            if (refDestructuringErrors && this.type === types$1.comma && refDestructuringErrors.trailingComma < 0) {
              refDestructuringErrors.trailingComma = this.start;
            }
          } else {
            elt = this.parseMaybeAssign(false, refDestructuringErrors);
          }
          elts.push(elt);
        }
        return elts;
      };
      pp$5.checkUnreserved = function(ref2) {
        var start2 = ref2.start;
        var end = ref2.end;
        var name2 = ref2.name;
        if (this.inGenerator && name2 === "yield") {
          this.raiseRecoverable(start2, "Cannot use 'yield' as identifier inside a generator");
        }
        if (this.inAsync && name2 === "await") {
          this.raiseRecoverable(start2, "Cannot use 'await' as identifier inside an async function");
        }
        if (!(this.currentThisScope().flags & SCOPE_VAR) && name2 === "arguments") {
          this.raiseRecoverable(start2, "Cannot use 'arguments' in class field initializer");
        }
        if (this.inClassStaticBlock && (name2 === "arguments" || name2 === "await")) {
          this.raise(start2, "Cannot use " + name2 + " in class static initialization block");
        }
        if (this.keywords.test(name2)) {
          this.raise(start2, "Unexpected keyword '" + name2 + "'");
        }
        if (this.options.ecmaVersion < 6 && this.input.slice(start2, end).indexOf("\\") !== -1) {
          return;
        }
        var re = this.strict ? this.reservedWordsStrict : this.reservedWords;
        if (re.test(name2)) {
          if (!this.inAsync && name2 === "await") {
            this.raiseRecoverable(start2, "Cannot use keyword 'await' outside an async function");
          }
          this.raiseRecoverable(start2, "The keyword '" + name2 + "' is reserved");
        }
      };
      pp$5.parseIdent = function(liberal) {
        var node2 = this.parseIdentNode();
        this.next(!!liberal);
        this.finishNode(node2, "Identifier");
        if (!liberal) {
          this.checkUnreserved(node2);
          if (node2.name === "await" && !this.awaitIdentPos) {
            this.awaitIdentPos = node2.start;
          }
        }
        return node2;
      };
      pp$5.parseIdentNode = function() {
        var node2 = this.startNode();
        if (this.type === types$1.name) {
          node2.name = this.value;
        } else if (this.type.keyword) {
          node2.name = this.type.keyword;
          if ((node2.name === "class" || node2.name === "function") && (this.lastTokEnd !== this.lastTokStart + 1 || this.input.charCodeAt(this.lastTokStart) !== 46)) {
            this.context.pop();
          }
          this.type = types$1.name;
        } else {
          this.unexpected();
        }
        return node2;
      };
      pp$5.parsePrivateIdent = function() {
        var node2 = this.startNode();
        if (this.type === types$1.privateId) {
          node2.name = this.value;
        } else {
          this.unexpected();
        }
        this.next();
        this.finishNode(node2, "PrivateIdentifier");
        if (this.options.checkPrivateFields) {
          if (this.privateNameStack.length === 0) {
            this.raise(node2.start, "Private field '#" + node2.name + "' must be declared in an enclosing class");
          } else {
            this.privateNameStack[this.privateNameStack.length - 1].used.push(node2);
          }
        }
        return node2;
      };
      pp$5.parseYield = function(forInit) {
        if (!this.yieldPos) {
          this.yieldPos = this.start;
        }
        var node2 = this.startNode();
        this.next();
        if (this.type === types$1.semi || this.canInsertSemicolon() || this.type !== types$1.star && !this.type.startsExpr) {
          node2.delegate = false;
          node2.argument = null;
        } else {
          node2.delegate = this.eat(types$1.star);
          node2.argument = this.parseMaybeAssign(forInit);
        }
        return this.finishNode(node2, "YieldExpression");
      };
      pp$5.parseAwait = function(forInit) {
        if (!this.awaitPos) {
          this.awaitPos = this.start;
        }
        var node2 = this.startNode();
        this.next();
        node2.argument = this.parseMaybeUnary(null, true, false, forInit);
        return this.finishNode(node2, "AwaitExpression");
      };
      pp$4 = Parser.prototype;
      pp$4.raise = function(pos, message) {
        var loc = getLineInfo(this.input, pos);
        message += " (" + loc.line + ":" + loc.column + ")";
        if (this.sourceFile) {
          message += " in " + this.sourceFile;
        }
        var err = new SyntaxError(message);
        err.pos = pos;
        err.loc = loc;
        err.raisedAt = this.pos;
        throw err;
      };
      pp$4.raiseRecoverable = pp$4.raise;
      pp$4.curPosition = function() {
        if (this.options.locations) {
          return new Position(this.curLine, this.pos - this.lineStart);
        }
      };
      pp$3 = Parser.prototype;
      Scope = function Scope2(flags) {
        this.flags = flags;
        this.var = [];
        this.lexical = [];
        this.functions = [];
      };
      pp$3.enterScope = function(flags) {
        this.scopeStack.push(new Scope(flags));
      };
      pp$3.exitScope = function() {
        this.scopeStack.pop();
      };
      pp$3.treatFunctionsAsVarInScope = function(scope) {
        return scope.flags & SCOPE_FUNCTION || !this.inModule && scope.flags & SCOPE_TOP;
      };
      pp$3.declareName = function(name2, bindingType, pos) {
        var redeclared = false;
        if (bindingType === BIND_LEXICAL) {
          var scope = this.currentScope();
          redeclared = scope.lexical.indexOf(name2) > -1 || scope.functions.indexOf(name2) > -1 || scope.var.indexOf(name2) > -1;
          scope.lexical.push(name2);
          if (this.inModule && scope.flags & SCOPE_TOP) {
            delete this.undefinedExports[name2];
          }
        } else if (bindingType === BIND_SIMPLE_CATCH) {
          var scope$1 = this.currentScope();
          scope$1.lexical.push(name2);
        } else if (bindingType === BIND_FUNCTION) {
          var scope$2 = this.currentScope();
          if (this.treatFunctionsAsVar) {
            redeclared = scope$2.lexical.indexOf(name2) > -1;
          } else {
            redeclared = scope$2.lexical.indexOf(name2) > -1 || scope$2.var.indexOf(name2) > -1;
          }
          scope$2.functions.push(name2);
        } else {
          for (var i = this.scopeStack.length - 1; i >= 0; --i) {
            var scope$3 = this.scopeStack[i];
            if (scope$3.lexical.indexOf(name2) > -1 && !(scope$3.flags & SCOPE_SIMPLE_CATCH && scope$3.lexical[0] === name2) || !this.treatFunctionsAsVarInScope(scope$3) && scope$3.functions.indexOf(name2) > -1) {
              redeclared = true;
              break;
            }
            scope$3.var.push(name2);
            if (this.inModule && scope$3.flags & SCOPE_TOP) {
              delete this.undefinedExports[name2];
            }
            if (scope$3.flags & SCOPE_VAR) {
              break;
            }
          }
        }
        if (redeclared) {
          this.raiseRecoverable(pos, "Identifier '" + name2 + "' has already been declared");
        }
      };
      pp$3.checkLocalExport = function(id) {
        if (this.scopeStack[0].lexical.indexOf(id.name) === -1 && this.scopeStack[0].var.indexOf(id.name) === -1) {
          this.undefinedExports[id.name] = id;
        }
      };
      pp$3.currentScope = function() {
        return this.scopeStack[this.scopeStack.length - 1];
      };
      pp$3.currentVarScope = function() {
        for (var i = this.scopeStack.length - 1; ; i--) {
          var scope = this.scopeStack[i];
          if (scope.flags & (SCOPE_VAR | SCOPE_CLASS_FIELD_INIT | SCOPE_CLASS_STATIC_BLOCK)) {
            return scope;
          }
        }
      };
      pp$3.currentThisScope = function() {
        for (var i = this.scopeStack.length - 1; ; i--) {
          var scope = this.scopeStack[i];
          if (scope.flags & (SCOPE_VAR | SCOPE_CLASS_FIELD_INIT | SCOPE_CLASS_STATIC_BLOCK) && !(scope.flags & SCOPE_ARROW)) {
            return scope;
          }
        }
      };
      Node = function Node2(parser, pos, loc) {
        this.type = "";
        this.start = pos;
        this.end = 0;
        if (parser.options.locations) {
          this.loc = new SourceLocation(parser, loc);
        }
        if (parser.options.directSourceFile) {
          this.sourceFile = parser.options.directSourceFile;
        }
        if (parser.options.ranges) {
          this.range = [pos, 0];
        }
      };
      pp$2 = Parser.prototype;
      pp$2.startNode = function() {
        return new Node(this, this.start, this.startLoc);
      };
      pp$2.startNodeAt = function(pos, loc) {
        return new Node(this, pos, loc);
      };
      pp$2.finishNode = function(node2, type) {
        return finishNodeAt.call(this, node2, type, this.lastTokEnd, this.lastTokEndLoc);
      };
      pp$2.finishNodeAt = function(node2, type, pos, loc) {
        return finishNodeAt.call(this, node2, type, pos, loc);
      };
      pp$2.copyNode = function(node2) {
        var newNode = new Node(this, node2.start, this.startLoc);
        for (var prop in node2) {
          newNode[prop] = node2[prop];
        }
        return newNode;
      };
      scriptValuesAddedInUnicode = "Gara Garay Gukh Gurung_Khema Hrkt Katakana_Or_Hiragana Kawi Kirat_Rai Krai Nag_Mundari Nagm Ol_Onal Onao Sunu Sunuwar Todhri Todr Tulu_Tigalari Tutg Unknown Zzzz";
      ecma9BinaryProperties = "ASCII ASCII_Hex_Digit AHex Alphabetic Alpha Any Assigned Bidi_Control Bidi_C Bidi_Mirrored Bidi_M Case_Ignorable CI Cased Changes_When_Casefolded CWCF Changes_When_Casemapped CWCM Changes_When_Lowercased CWL Changes_When_NFKC_Casefolded CWKCF Changes_When_Titlecased CWT Changes_When_Uppercased CWU Dash Default_Ignorable_Code_Point DI Deprecated Dep Diacritic Dia Emoji Emoji_Component Emoji_Modifier Emoji_Modifier_Base Emoji_Presentation Extender Ext Grapheme_Base Gr_Base Grapheme_Extend Gr_Ext Hex_Digit Hex IDS_Binary_Operator IDSB IDS_Trinary_Operator IDST ID_Continue IDC ID_Start IDS Ideographic Ideo Join_Control Join_C Logical_Order_Exception LOE Lowercase Lower Math Noncharacter_Code_Point NChar Pattern_Syntax Pat_Syn Pattern_White_Space Pat_WS Quotation_Mark QMark Radical Regional_Indicator RI Sentence_Terminal STerm Soft_Dotted SD Terminal_Punctuation Term Unified_Ideograph UIdeo Uppercase Upper Variation_Selector VS White_Space space XID_Continue XIDC XID_Start XIDS";
      ecma10BinaryProperties = ecma9BinaryProperties + " Extended_Pictographic";
      ecma11BinaryProperties = ecma10BinaryProperties;
      ecma12BinaryProperties = ecma11BinaryProperties + " EBase EComp EMod EPres ExtPict";
      ecma13BinaryProperties = ecma12BinaryProperties;
      ecma14BinaryProperties = ecma13BinaryProperties;
      unicodeBinaryProperties = {
        9: ecma9BinaryProperties,
        10: ecma10BinaryProperties,
        11: ecma11BinaryProperties,
        12: ecma12BinaryProperties,
        13: ecma13BinaryProperties,
        14: ecma14BinaryProperties
      };
      ecma14BinaryPropertiesOfStrings = "Basic_Emoji Emoji_Keycap_Sequence RGI_Emoji_Modifier_Sequence RGI_Emoji_Flag_Sequence RGI_Emoji_Tag_Sequence RGI_Emoji_ZWJ_Sequence RGI_Emoji";
      unicodeBinaryPropertiesOfStrings = {
        9: "",
        10: "",
        11: "",
        12: "",
        13: "",
        14: ecma14BinaryPropertiesOfStrings
      };
      unicodeGeneralCategoryValues = "Cased_Letter LC Close_Punctuation Pe Connector_Punctuation Pc Control Cc cntrl Currency_Symbol Sc Dash_Punctuation Pd Decimal_Number Nd digit Enclosing_Mark Me Final_Punctuation Pf Format Cf Initial_Punctuation Pi Letter L Letter_Number Nl Line_Separator Zl Lowercase_Letter Ll Mark M Combining_Mark Math_Symbol Sm Modifier_Letter Lm Modifier_Symbol Sk Nonspacing_Mark Mn Number N Open_Punctuation Ps Other C Other_Letter Lo Other_Number No Other_Punctuation Po Other_Symbol So Paragraph_Separator Zp Private_Use Co Punctuation P punct Separator Z Space_Separator Zs Spacing_Mark Mc Surrogate Cs Symbol S Titlecase_Letter Lt Unassigned Cn Uppercase_Letter Lu";
      ecma9ScriptValues = "Adlam Adlm Ahom Anatolian_Hieroglyphs Hluw Arabic Arab Armenian Armn Avestan Avst Balinese Bali Bamum Bamu Bassa_Vah Bass Batak Batk Bengali Beng Bhaiksuki Bhks Bopomofo Bopo Brahmi Brah Braille Brai Buginese Bugi Buhid Buhd Canadian_Aboriginal Cans Carian Cari Caucasian_Albanian Aghb Chakma Cakm Cham Cham Cherokee Cher Common Zyyy Coptic Copt Qaac Cuneiform Xsux Cypriot Cprt Cyrillic Cyrl Deseret Dsrt Devanagari Deva Duployan Dupl Egyptian_Hieroglyphs Egyp Elbasan Elba Ethiopic Ethi Georgian Geor Glagolitic Glag Gothic Goth Grantha Gran Greek Grek Gujarati Gujr Gurmukhi Guru Han Hani Hangul Hang Hanunoo Hano Hatran Hatr Hebrew Hebr Hiragana Hira Imperial_Aramaic Armi Inherited Zinh Qaai Inscriptional_Pahlavi Phli Inscriptional_Parthian Prti Javanese Java Kaithi Kthi Kannada Knda Katakana Kana Kayah_Li Kali Kharoshthi Khar Khmer Khmr Khojki Khoj Khudawadi Sind Lao Laoo Latin Latn Lepcha Lepc Limbu Limb Linear_A Lina Linear_B Linb Lisu Lisu Lycian Lyci Lydian Lydi Mahajani Mahj Malayalam Mlym Mandaic Mand Manichaean Mani Marchen Marc Masaram_Gondi Gonm Meetei_Mayek Mtei Mende_Kikakui Mend Meroitic_Cursive Merc Meroitic_Hieroglyphs Mero Miao Plrd Modi Mongolian Mong Mro Mroo Multani Mult Myanmar Mymr Nabataean Nbat New_Tai_Lue Talu Newa Newa Nko Nkoo Nushu Nshu Ogham Ogam Ol_Chiki Olck Old_Hungarian Hung Old_Italic Ital Old_North_Arabian Narb Old_Permic Perm Old_Persian Xpeo Old_South_Arabian Sarb Old_Turkic Orkh Oriya Orya Osage Osge Osmanya Osma Pahawh_Hmong Hmng Palmyrene Palm Pau_Cin_Hau Pauc Phags_Pa Phag Phoenician Phnx Psalter_Pahlavi Phlp Rejang Rjng Runic Runr Samaritan Samr Saurashtra Saur Sharada Shrd Shavian Shaw Siddham Sidd SignWriting Sgnw Sinhala Sinh Sora_Sompeng Sora Soyombo Soyo Sundanese Sund Syloti_Nagri Sylo Syriac Syrc Tagalog Tglg Tagbanwa Tagb Tai_Le Tale Tai_Tham Lana Tai_Viet Tavt Takri Takr Tamil Taml Tangut Tang Telugu Telu Thaana Thaa Thai Thai Tibetan Tibt Tifinagh Tfng Tirhuta Tirh Ugaritic Ugar Vai Vaii Warang_Citi Wara Yi Yiii Zanabazar_Square Zanb";
      ecma10ScriptValues = ecma9ScriptValues + " Dogra Dogr Gunjala_Gondi Gong Hanifi_Rohingya Rohg Makasar Maka Medefaidrin Medf Old_Sogdian Sogo Sogdian Sogd";
      ecma11ScriptValues = ecma10ScriptValues + " Elymaic Elym Nandinagari Nand Nyiakeng_Puachue_Hmong Hmnp Wancho Wcho";
      ecma12ScriptValues = ecma11ScriptValues + " Chorasmian Chrs Diak Dives_Akuru Khitan_Small_Script Kits Yezi Yezidi";
      ecma13ScriptValues = ecma12ScriptValues + " Cypro_Minoan Cpmn Old_Uyghur Ougr Tangsa Tnsa Toto Vithkuqi Vith";
      ecma14ScriptValues = ecma13ScriptValues + " " + scriptValuesAddedInUnicode;
      unicodeScriptValues = {
        9: ecma9ScriptValues,
        10: ecma10ScriptValues,
        11: ecma11ScriptValues,
        12: ecma12ScriptValues,
        13: ecma13ScriptValues,
        14: ecma14ScriptValues
      };
      data = {};
      for (i = 0, list4 = [9, 10, 11, 12, 13, 14]; i < list4.length; i += 1) {
        ecmaVersion = list4[i];
        buildUnicodeData(ecmaVersion);
      }
      pp$1 = Parser.prototype;
      BranchID = function BranchID2(parent, base) {
        this.parent = parent;
        this.base = base || this;
      };
      BranchID.prototype.separatedFrom = function separatedFrom(alt) {
        for (var self2 = this; self2; self2 = self2.parent) {
          for (var other = alt; other; other = other.parent) {
            if (self2.base === other.base && self2 !== other) {
              return true;
            }
          }
        }
        return false;
      };
      BranchID.prototype.sibling = function sibling() {
        return new BranchID(this.parent, this.base);
      };
      RegExpValidationState = function RegExpValidationState2(parser) {
        this.parser = parser;
        this.validFlags = "gim" + (parser.options.ecmaVersion >= 6 ? "uy" : "") + (parser.options.ecmaVersion >= 9 ? "s" : "") + (parser.options.ecmaVersion >= 13 ? "d" : "") + (parser.options.ecmaVersion >= 15 ? "v" : "");
        this.unicodeProperties = data[parser.options.ecmaVersion >= 14 ? 14 : parser.options.ecmaVersion];
        this.source = "";
        this.flags = "";
        this.start = 0;
        this.switchU = false;
        this.switchV = false;
        this.switchN = false;
        this.pos = 0;
        this.lastIntValue = 0;
        this.lastStringValue = "";
        this.lastAssertionIsQuantifiable = false;
        this.numCapturingParens = 0;
        this.maxBackReference = 0;
        this.groupNames = /* @__PURE__ */ Object.create(null);
        this.backReferenceNames = [];
        this.branchID = null;
      };
      RegExpValidationState.prototype.reset = function reset(start2, pattern, flags) {
        var unicodeSets = flags.indexOf("v") !== -1;
        var unicode = flags.indexOf("u") !== -1;
        this.start = start2 | 0;
        this.source = pattern + "";
        this.flags = flags;
        if (unicodeSets && this.parser.options.ecmaVersion >= 15) {
          this.switchU = true;
          this.switchV = true;
          this.switchN = true;
        } else {
          this.switchU = unicode && this.parser.options.ecmaVersion >= 6;
          this.switchV = false;
          this.switchN = unicode && this.parser.options.ecmaVersion >= 9;
        }
      };
      RegExpValidationState.prototype.raise = function raise(message) {
        this.parser.raiseRecoverable(this.start, "Invalid regular expression: /" + this.source + "/: " + message);
      };
      RegExpValidationState.prototype.at = function at(i, forceU) {
        if (forceU === void 0) forceU = false;
        var s = this.source;
        var l = s.length;
        if (i >= l) {
          return -1;
        }
        var c2 = s.charCodeAt(i);
        if (!(forceU || this.switchU) || c2 <= 55295 || c2 >= 57344 || i + 1 >= l) {
          return c2;
        }
        var next = s.charCodeAt(i + 1);
        return next >= 56320 && next <= 57343 ? (c2 << 10) + next - 56613888 : c2;
      };
      RegExpValidationState.prototype.nextIndex = function nextIndex(i, forceU) {
        if (forceU === void 0) forceU = false;
        var s = this.source;
        var l = s.length;
        if (i >= l) {
          return l;
        }
        var c2 = s.charCodeAt(i), next;
        if (!(forceU || this.switchU) || c2 <= 55295 || c2 >= 57344 || i + 1 >= l || (next = s.charCodeAt(i + 1)) < 56320 || next > 57343) {
          return i + 1;
        }
        return i + 2;
      };
      RegExpValidationState.prototype.current = function current(forceU) {
        if (forceU === void 0) forceU = false;
        return this.at(this.pos, forceU);
      };
      RegExpValidationState.prototype.lookahead = function lookahead(forceU) {
        if (forceU === void 0) forceU = false;
        return this.at(this.nextIndex(this.pos, forceU), forceU);
      };
      RegExpValidationState.prototype.advance = function advance(forceU) {
        if (forceU === void 0) forceU = false;
        this.pos = this.nextIndex(this.pos, forceU);
      };
      RegExpValidationState.prototype.eat = function eat(ch, forceU) {
        if (forceU === void 0) forceU = false;
        if (this.current(forceU) === ch) {
          this.advance(forceU);
          return true;
        }
        return false;
      };
      RegExpValidationState.prototype.eatChars = function eatChars(chs, forceU) {
        if (forceU === void 0) forceU = false;
        var pos = this.pos;
        for (var i = 0, list4 = chs; i < list4.length; i += 1) {
          var ch = list4[i];
          var current2 = this.at(pos, forceU);
          if (current2 === -1 || current2 !== ch) {
            return false;
          }
          pos = this.nextIndex(pos, forceU);
        }
        this.pos = pos;
        return true;
      };
      pp$1.validateRegExpFlags = function(state) {
        var validFlags = state.validFlags;
        var flags = state.flags;
        var u = false;
        var v2 = false;
        for (var i = 0; i < flags.length; i++) {
          var flag = flags.charAt(i);
          if (validFlags.indexOf(flag) === -1) {
            this.raise(state.start, "Invalid regular expression flag");
          }
          if (flags.indexOf(flag, i + 1) > -1) {
            this.raise(state.start, "Duplicate regular expression flag");
          }
          if (flag === "u") {
            u = true;
          }
          if (flag === "v") {
            v2 = true;
          }
        }
        if (this.options.ecmaVersion >= 15 && u && v2) {
          this.raise(state.start, "Invalid regular expression flag");
        }
      };
      pp$1.validateRegExpPattern = function(state) {
        this.regexp_pattern(state);
        if (!state.switchN && this.options.ecmaVersion >= 9 && hasProp(state.groupNames)) {
          state.switchN = true;
          this.regexp_pattern(state);
        }
      };
      pp$1.regexp_pattern = function(state) {
        state.pos = 0;
        state.lastIntValue = 0;
        state.lastStringValue = "";
        state.lastAssertionIsQuantifiable = false;
        state.numCapturingParens = 0;
        state.maxBackReference = 0;
        state.groupNames = /* @__PURE__ */ Object.create(null);
        state.backReferenceNames.length = 0;
        state.branchID = null;
        this.regexp_disjunction(state);
        if (state.pos !== state.source.length) {
          if (state.eat(
            41
            /* ) */
          )) {
            state.raise("Unmatched ')'");
          }
          if (state.eat(
            93
            /* ] */
          ) || state.eat(
            125
            /* } */
          )) {
            state.raise("Lone quantifier brackets");
          }
        }
        if (state.maxBackReference > state.numCapturingParens) {
          state.raise("Invalid escape");
        }
        for (var i = 0, list4 = state.backReferenceNames; i < list4.length; i += 1) {
          var name2 = list4[i];
          if (!state.groupNames[name2]) {
            state.raise("Invalid named capture referenced");
          }
        }
      };
      pp$1.regexp_disjunction = function(state) {
        var trackDisjunction = this.options.ecmaVersion >= 16;
        if (trackDisjunction) {
          state.branchID = new BranchID(state.branchID, null);
        }
        this.regexp_alternative(state);
        while (state.eat(
          124
          /* | */
        )) {
          if (trackDisjunction) {
            state.branchID = state.branchID.sibling();
          }
          this.regexp_alternative(state);
        }
        if (trackDisjunction) {
          state.branchID = state.branchID.parent;
        }
        if (this.regexp_eatQuantifier(state, true)) {
          state.raise("Nothing to repeat");
        }
        if (state.eat(
          123
          /* { */
        )) {
          state.raise("Lone quantifier brackets");
        }
      };
      pp$1.regexp_alternative = function(state) {
        while (state.pos < state.source.length && this.regexp_eatTerm(state)) {
        }
      };
      pp$1.regexp_eatTerm = function(state) {
        if (this.regexp_eatAssertion(state)) {
          if (state.lastAssertionIsQuantifiable && this.regexp_eatQuantifier(state)) {
            if (state.switchU) {
              state.raise("Invalid quantifier");
            }
          }
          return true;
        }
        if (state.switchU ? this.regexp_eatAtom(state) : this.regexp_eatExtendedAtom(state)) {
          this.regexp_eatQuantifier(state);
          return true;
        }
        return false;
      };
      pp$1.regexp_eatAssertion = function(state) {
        var start2 = state.pos;
        state.lastAssertionIsQuantifiable = false;
        if (state.eat(
          94
          /* ^ */
        ) || state.eat(
          36
          /* $ */
        )) {
          return true;
        }
        if (state.eat(
          92
          /* \ */
        )) {
          if (state.eat(
            66
            /* B */
          ) || state.eat(
            98
            /* b */
          )) {
            return true;
          }
          state.pos = start2;
        }
        if (state.eat(
          40
          /* ( */
        ) && state.eat(
          63
          /* ? */
        )) {
          var lookbehind = false;
          if (this.options.ecmaVersion >= 9) {
            lookbehind = state.eat(
              60
              /* < */
            );
          }
          if (state.eat(
            61
            /* = */
          ) || state.eat(
            33
            /* ! */
          )) {
            this.regexp_disjunction(state);
            if (!state.eat(
              41
              /* ) */
            )) {
              state.raise("Unterminated group");
            }
            state.lastAssertionIsQuantifiable = !lookbehind;
            return true;
          }
        }
        state.pos = start2;
        return false;
      };
      pp$1.regexp_eatQuantifier = function(state, noError) {
        if (noError === void 0) noError = false;
        if (this.regexp_eatQuantifierPrefix(state, noError)) {
          state.eat(
            63
            /* ? */
          );
          return true;
        }
        return false;
      };
      pp$1.regexp_eatQuantifierPrefix = function(state, noError) {
        return state.eat(
          42
          /* * */
        ) || state.eat(
          43
          /* + */
        ) || state.eat(
          63
          /* ? */
        ) || this.regexp_eatBracedQuantifier(state, noError);
      };
      pp$1.regexp_eatBracedQuantifier = function(state, noError) {
        var start2 = state.pos;
        if (state.eat(
          123
          /* { */
        )) {
          var min = 0, max = -1;
          if (this.regexp_eatDecimalDigits(state)) {
            min = state.lastIntValue;
            if (state.eat(
              44
              /* , */
            ) && this.regexp_eatDecimalDigits(state)) {
              max = state.lastIntValue;
            }
            if (state.eat(
              125
              /* } */
            )) {
              if (max !== -1 && max < min && !noError) {
                state.raise("numbers out of order in {} quantifier");
              }
              return true;
            }
          }
          if (state.switchU && !noError) {
            state.raise("Incomplete quantifier");
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatAtom = function(state) {
        return this.regexp_eatPatternCharacters(state) || state.eat(
          46
          /* . */
        ) || this.regexp_eatReverseSolidusAtomEscape(state) || this.regexp_eatCharacterClass(state) || this.regexp_eatUncapturingGroup(state) || this.regexp_eatCapturingGroup(state);
      };
      pp$1.regexp_eatReverseSolidusAtomEscape = function(state) {
        var start2 = state.pos;
        if (state.eat(
          92
          /* \ */
        )) {
          if (this.regexp_eatAtomEscape(state)) {
            return true;
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatUncapturingGroup = function(state) {
        var start2 = state.pos;
        if (state.eat(
          40
          /* ( */
        )) {
          if (state.eat(
            63
            /* ? */
          )) {
            if (this.options.ecmaVersion >= 16) {
              var addModifiers = this.regexp_eatModifiers(state);
              var hasHyphen = state.eat(
                45
                /* - */
              );
              if (addModifiers || hasHyphen) {
                for (var i = 0; i < addModifiers.length; i++) {
                  var modifier = addModifiers.charAt(i);
                  if (addModifiers.indexOf(modifier, i + 1) > -1) {
                    state.raise("Duplicate regular expression modifiers");
                  }
                }
                if (hasHyphen) {
                  var removeModifiers = this.regexp_eatModifiers(state);
                  if (!addModifiers && !removeModifiers && state.current() === 58) {
                    state.raise("Invalid regular expression modifiers");
                  }
                  for (var i$1 = 0; i$1 < removeModifiers.length; i$1++) {
                    var modifier$1 = removeModifiers.charAt(i$1);
                    if (removeModifiers.indexOf(modifier$1, i$1 + 1) > -1 || addModifiers.indexOf(modifier$1) > -1) {
                      state.raise("Duplicate regular expression modifiers");
                    }
                  }
                }
              }
            }
            if (state.eat(
              58
              /* : */
            )) {
              this.regexp_disjunction(state);
              if (state.eat(
                41
                /* ) */
              )) {
                return true;
              }
              state.raise("Unterminated group");
            }
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatCapturingGroup = function(state) {
        if (state.eat(
          40
          /* ( */
        )) {
          if (this.options.ecmaVersion >= 9) {
            this.regexp_groupSpecifier(state);
          } else if (state.current() === 63) {
            state.raise("Invalid group");
          }
          this.regexp_disjunction(state);
          if (state.eat(
            41
            /* ) */
          )) {
            state.numCapturingParens += 1;
            return true;
          }
          state.raise("Unterminated group");
        }
        return false;
      };
      pp$1.regexp_eatModifiers = function(state) {
        var modifiers = "";
        var ch = 0;
        while ((ch = state.current()) !== -1 && isRegularExpressionModifier(ch)) {
          modifiers += codePointToString(ch);
          state.advance();
        }
        return modifiers;
      };
      pp$1.regexp_eatExtendedAtom = function(state) {
        return state.eat(
          46
          /* . */
        ) || this.regexp_eatReverseSolidusAtomEscape(state) || this.regexp_eatCharacterClass(state) || this.regexp_eatUncapturingGroup(state) || this.regexp_eatCapturingGroup(state) || this.regexp_eatInvalidBracedQuantifier(state) || this.regexp_eatExtendedPatternCharacter(state);
      };
      pp$1.regexp_eatInvalidBracedQuantifier = function(state) {
        if (this.regexp_eatBracedQuantifier(state, true)) {
          state.raise("Nothing to repeat");
        }
        return false;
      };
      pp$1.regexp_eatSyntaxCharacter = function(state) {
        var ch = state.current();
        if (isSyntaxCharacter(ch)) {
          state.lastIntValue = ch;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatPatternCharacters = function(state) {
        var start2 = state.pos;
        var ch = 0;
        while ((ch = state.current()) !== -1 && !isSyntaxCharacter(ch)) {
          state.advance();
        }
        return state.pos !== start2;
      };
      pp$1.regexp_eatExtendedPatternCharacter = function(state) {
        var ch = state.current();
        if (ch !== -1 && ch !== 36 && !(ch >= 40 && ch <= 43) && ch !== 46 && ch !== 63 && ch !== 91 && ch !== 94 && ch !== 124) {
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_groupSpecifier = function(state) {
        if (state.eat(
          63
          /* ? */
        )) {
          if (!this.regexp_eatGroupName(state)) {
            state.raise("Invalid group");
          }
          var trackDisjunction = this.options.ecmaVersion >= 16;
          var known = state.groupNames[state.lastStringValue];
          if (known) {
            if (trackDisjunction) {
              for (var i = 0, list4 = known; i < list4.length; i += 1) {
                var altID = list4[i];
                if (!altID.separatedFrom(state.branchID)) {
                  state.raise("Duplicate capture group name");
                }
              }
            } else {
              state.raise("Duplicate capture group name");
            }
          }
          if (trackDisjunction) {
            (known || (state.groupNames[state.lastStringValue] = [])).push(state.branchID);
          } else {
            state.groupNames[state.lastStringValue] = true;
          }
        }
      };
      pp$1.regexp_eatGroupName = function(state) {
        state.lastStringValue = "";
        if (state.eat(
          60
          /* < */
        )) {
          if (this.regexp_eatRegExpIdentifierName(state) && state.eat(
            62
            /* > */
          )) {
            return true;
          }
          state.raise("Invalid capture group name");
        }
        return false;
      };
      pp$1.regexp_eatRegExpIdentifierName = function(state) {
        state.lastStringValue = "";
        if (this.regexp_eatRegExpIdentifierStart(state)) {
          state.lastStringValue += codePointToString(state.lastIntValue);
          while (this.regexp_eatRegExpIdentifierPart(state)) {
            state.lastStringValue += codePointToString(state.lastIntValue);
          }
          return true;
        }
        return false;
      };
      pp$1.regexp_eatRegExpIdentifierStart = function(state) {
        var start2 = state.pos;
        var forceU = this.options.ecmaVersion >= 11;
        var ch = state.current(forceU);
        state.advance(forceU);
        if (ch === 92 && this.regexp_eatRegExpUnicodeEscapeSequence(state, forceU)) {
          ch = state.lastIntValue;
        }
        if (isRegExpIdentifierStart(ch)) {
          state.lastIntValue = ch;
          return true;
        }
        state.pos = start2;
        return false;
      };
      pp$1.regexp_eatRegExpIdentifierPart = function(state) {
        var start2 = state.pos;
        var forceU = this.options.ecmaVersion >= 11;
        var ch = state.current(forceU);
        state.advance(forceU);
        if (ch === 92 && this.regexp_eatRegExpUnicodeEscapeSequence(state, forceU)) {
          ch = state.lastIntValue;
        }
        if (isRegExpIdentifierPart(ch)) {
          state.lastIntValue = ch;
          return true;
        }
        state.pos = start2;
        return false;
      };
      pp$1.regexp_eatAtomEscape = function(state) {
        if (this.regexp_eatBackReference(state) || this.regexp_eatCharacterClassEscape(state) || this.regexp_eatCharacterEscape(state) || state.switchN && this.regexp_eatKGroupName(state)) {
          return true;
        }
        if (state.switchU) {
          if (state.current() === 99) {
            state.raise("Invalid unicode escape");
          }
          state.raise("Invalid escape");
        }
        return false;
      };
      pp$1.regexp_eatBackReference = function(state) {
        var start2 = state.pos;
        if (this.regexp_eatDecimalEscape(state)) {
          var n = state.lastIntValue;
          if (state.switchU) {
            if (n > state.maxBackReference) {
              state.maxBackReference = n;
            }
            return true;
          }
          if (n <= state.numCapturingParens) {
            return true;
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatKGroupName = function(state) {
        if (state.eat(
          107
          /* k */
        )) {
          if (this.regexp_eatGroupName(state)) {
            state.backReferenceNames.push(state.lastStringValue);
            return true;
          }
          state.raise("Invalid named reference");
        }
        return false;
      };
      pp$1.regexp_eatCharacterEscape = function(state) {
        return this.regexp_eatControlEscape(state) || this.regexp_eatCControlLetter(state) || this.regexp_eatZero(state) || this.regexp_eatHexEscapeSequence(state) || this.regexp_eatRegExpUnicodeEscapeSequence(state, false) || !state.switchU && this.regexp_eatLegacyOctalEscapeSequence(state) || this.regexp_eatIdentityEscape(state);
      };
      pp$1.regexp_eatCControlLetter = function(state) {
        var start2 = state.pos;
        if (state.eat(
          99
          /* c */
        )) {
          if (this.regexp_eatControlLetter(state)) {
            return true;
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatZero = function(state) {
        if (state.current() === 48 && !isDecimalDigit(state.lookahead())) {
          state.lastIntValue = 0;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatControlEscape = function(state) {
        var ch = state.current();
        if (ch === 116) {
          state.lastIntValue = 9;
          state.advance();
          return true;
        }
        if (ch === 110) {
          state.lastIntValue = 10;
          state.advance();
          return true;
        }
        if (ch === 118) {
          state.lastIntValue = 11;
          state.advance();
          return true;
        }
        if (ch === 102) {
          state.lastIntValue = 12;
          state.advance();
          return true;
        }
        if (ch === 114) {
          state.lastIntValue = 13;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatControlLetter = function(state) {
        var ch = state.current();
        if (isControlLetter(ch)) {
          state.lastIntValue = ch % 32;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatRegExpUnicodeEscapeSequence = function(state, forceU) {
        if (forceU === void 0) forceU = false;
        var start2 = state.pos;
        var switchU = forceU || state.switchU;
        if (state.eat(
          117
          /* u */
        )) {
          if (this.regexp_eatFixedHexDigits(state, 4)) {
            var lead = state.lastIntValue;
            if (switchU && lead >= 55296 && lead <= 56319) {
              var leadSurrogateEnd = state.pos;
              if (state.eat(
                92
                /* \ */
              ) && state.eat(
                117
                /* u */
              ) && this.regexp_eatFixedHexDigits(state, 4)) {
                var trail2 = state.lastIntValue;
                if (trail2 >= 56320 && trail2 <= 57343) {
                  state.lastIntValue = (lead - 55296) * 1024 + (trail2 - 56320) + 65536;
                  return true;
                }
              }
              state.pos = leadSurrogateEnd;
              state.lastIntValue = lead;
            }
            return true;
          }
          if (switchU && state.eat(
            123
            /* { */
          ) && this.regexp_eatHexDigits(state) && state.eat(
            125
            /* } */
          ) && isValidUnicode(state.lastIntValue)) {
            return true;
          }
          if (switchU) {
            state.raise("Invalid unicode escape");
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatIdentityEscape = function(state) {
        if (state.switchU) {
          if (this.regexp_eatSyntaxCharacter(state)) {
            return true;
          }
          if (state.eat(
            47
            /* / */
          )) {
            state.lastIntValue = 47;
            return true;
          }
          return false;
        }
        var ch = state.current();
        if (ch !== 99 && (!state.switchN || ch !== 107)) {
          state.lastIntValue = ch;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatDecimalEscape = function(state) {
        state.lastIntValue = 0;
        var ch = state.current();
        if (ch >= 49 && ch <= 57) {
          do {
            state.lastIntValue = 10 * state.lastIntValue + (ch - 48);
            state.advance();
          } while ((ch = state.current()) >= 48 && ch <= 57);
          return true;
        }
        return false;
      };
      CharSetNone = 0;
      CharSetOk = 1;
      CharSetString = 2;
      pp$1.regexp_eatCharacterClassEscape = function(state) {
        var ch = state.current();
        if (isCharacterClassEscape(ch)) {
          state.lastIntValue = -1;
          state.advance();
          return CharSetOk;
        }
        var negate = false;
        if (state.switchU && this.options.ecmaVersion >= 9 && ((negate = ch === 80) || ch === 112)) {
          state.lastIntValue = -1;
          state.advance();
          var result;
          if (state.eat(
            123
            /* { */
          ) && (result = this.regexp_eatUnicodePropertyValueExpression(state)) && state.eat(
            125
            /* } */
          )) {
            if (negate && result === CharSetString) {
              state.raise("Invalid property name");
            }
            return result;
          }
          state.raise("Invalid property name");
        }
        return CharSetNone;
      };
      pp$1.regexp_eatUnicodePropertyValueExpression = function(state) {
        var start2 = state.pos;
        if (this.regexp_eatUnicodePropertyName(state) && state.eat(
          61
          /* = */
        )) {
          var name2 = state.lastStringValue;
          if (this.regexp_eatUnicodePropertyValue(state)) {
            var value = state.lastStringValue;
            this.regexp_validateUnicodePropertyNameAndValue(state, name2, value);
            return CharSetOk;
          }
        }
        state.pos = start2;
        if (this.regexp_eatLoneUnicodePropertyNameOrValue(state)) {
          var nameOrValue = state.lastStringValue;
          return this.regexp_validateUnicodePropertyNameOrValue(state, nameOrValue);
        }
        return CharSetNone;
      };
      pp$1.regexp_validateUnicodePropertyNameAndValue = function(state, name2, value) {
        if (!hasOwn(state.unicodeProperties.nonBinary, name2)) {
          state.raise("Invalid property name");
        }
        if (!state.unicodeProperties.nonBinary[name2].test(value)) {
          state.raise("Invalid property value");
        }
      };
      pp$1.regexp_validateUnicodePropertyNameOrValue = function(state, nameOrValue) {
        if (state.unicodeProperties.binary.test(nameOrValue)) {
          return CharSetOk;
        }
        if (state.switchV && state.unicodeProperties.binaryOfStrings.test(nameOrValue)) {
          return CharSetString;
        }
        state.raise("Invalid property name");
      };
      pp$1.regexp_eatUnicodePropertyName = function(state) {
        var ch = 0;
        state.lastStringValue = "";
        while (isUnicodePropertyNameCharacter(ch = state.current())) {
          state.lastStringValue += codePointToString(ch);
          state.advance();
        }
        return state.lastStringValue !== "";
      };
      pp$1.regexp_eatUnicodePropertyValue = function(state) {
        var ch = 0;
        state.lastStringValue = "";
        while (isUnicodePropertyValueCharacter(ch = state.current())) {
          state.lastStringValue += codePointToString(ch);
          state.advance();
        }
        return state.lastStringValue !== "";
      };
      pp$1.regexp_eatLoneUnicodePropertyNameOrValue = function(state) {
        return this.regexp_eatUnicodePropertyValue(state);
      };
      pp$1.regexp_eatCharacterClass = function(state) {
        if (state.eat(
          91
          /* [ */
        )) {
          var negate = state.eat(
            94
            /* ^ */
          );
          var result = this.regexp_classContents(state);
          if (!state.eat(
            93
            /* ] */
          )) {
            state.raise("Unterminated character class");
          }
          if (negate && result === CharSetString) {
            state.raise("Negated character class may contain strings");
          }
          return true;
        }
        return false;
      };
      pp$1.regexp_classContents = function(state) {
        if (state.current() === 93) {
          return CharSetOk;
        }
        if (state.switchV) {
          return this.regexp_classSetExpression(state);
        }
        this.regexp_nonEmptyClassRanges(state);
        return CharSetOk;
      };
      pp$1.regexp_nonEmptyClassRanges = function(state) {
        while (this.regexp_eatClassAtom(state)) {
          var left = state.lastIntValue;
          if (state.eat(
            45
            /* - */
          ) && this.regexp_eatClassAtom(state)) {
            var right = state.lastIntValue;
            if (state.switchU && (left === -1 || right === -1)) {
              state.raise("Invalid character class");
            }
            if (left !== -1 && right !== -1 && left > right) {
              state.raise("Range out of order in character class");
            }
          }
        }
      };
      pp$1.regexp_eatClassAtom = function(state) {
        var start2 = state.pos;
        if (state.eat(
          92
          /* \ */
        )) {
          if (this.regexp_eatClassEscape(state)) {
            return true;
          }
          if (state.switchU) {
            var ch$1 = state.current();
            if (ch$1 === 99 || isOctalDigit(ch$1)) {
              state.raise("Invalid class escape");
            }
            state.raise("Invalid escape");
          }
          state.pos = start2;
        }
        var ch = state.current();
        if (ch !== 93) {
          state.lastIntValue = ch;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatClassEscape = function(state) {
        var start2 = state.pos;
        if (state.eat(
          98
          /* b */
        )) {
          state.lastIntValue = 8;
          return true;
        }
        if (state.switchU && state.eat(
          45
          /* - */
        )) {
          state.lastIntValue = 45;
          return true;
        }
        if (!state.switchU && state.eat(
          99
          /* c */
        )) {
          if (this.regexp_eatClassControlLetter(state)) {
            return true;
          }
          state.pos = start2;
        }
        return this.regexp_eatCharacterClassEscape(state) || this.regexp_eatCharacterEscape(state);
      };
      pp$1.regexp_classSetExpression = function(state) {
        var result = CharSetOk, subResult;
        if (this.regexp_eatClassSetRange(state)) ;
        else if (subResult = this.regexp_eatClassSetOperand(state)) {
          if (subResult === CharSetString) {
            result = CharSetString;
          }
          var start2 = state.pos;
          while (state.eatChars(
            [38, 38]
            /* && */
          )) {
            if (state.current() !== 38 && (subResult = this.regexp_eatClassSetOperand(state))) {
              if (subResult !== CharSetString) {
                result = CharSetOk;
              }
              continue;
            }
            state.raise("Invalid character in character class");
          }
          if (start2 !== state.pos) {
            return result;
          }
          while (state.eatChars(
            [45, 45]
            /* -- */
          )) {
            if (this.regexp_eatClassSetOperand(state)) {
              continue;
            }
            state.raise("Invalid character in character class");
          }
          if (start2 !== state.pos) {
            return result;
          }
        } else {
          state.raise("Invalid character in character class");
        }
        for (; ; ) {
          if (this.regexp_eatClassSetRange(state)) {
            continue;
          }
          subResult = this.regexp_eatClassSetOperand(state);
          if (!subResult) {
            return result;
          }
          if (subResult === CharSetString) {
            result = CharSetString;
          }
        }
      };
      pp$1.regexp_eatClassSetRange = function(state) {
        var start2 = state.pos;
        if (this.regexp_eatClassSetCharacter(state)) {
          var left = state.lastIntValue;
          if (state.eat(
            45
            /* - */
          ) && this.regexp_eatClassSetCharacter(state)) {
            var right = state.lastIntValue;
            if (left !== -1 && right !== -1 && left > right) {
              state.raise("Range out of order in character class");
            }
            return true;
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatClassSetOperand = function(state) {
        if (this.regexp_eatClassSetCharacter(state)) {
          return CharSetOk;
        }
        return this.regexp_eatClassStringDisjunction(state) || this.regexp_eatNestedClass(state);
      };
      pp$1.regexp_eatNestedClass = function(state) {
        var start2 = state.pos;
        if (state.eat(
          91
          /* [ */
        )) {
          var negate = state.eat(
            94
            /* ^ */
          );
          var result = this.regexp_classContents(state);
          if (state.eat(
            93
            /* ] */
          )) {
            if (negate && result === CharSetString) {
              state.raise("Negated character class may contain strings");
            }
            return result;
          }
          state.pos = start2;
        }
        if (state.eat(
          92
          /* \ */
        )) {
          var result$1 = this.regexp_eatCharacterClassEscape(state);
          if (result$1) {
            return result$1;
          }
          state.pos = start2;
        }
        return null;
      };
      pp$1.regexp_eatClassStringDisjunction = function(state) {
        var start2 = state.pos;
        if (state.eatChars(
          [92, 113]
          /* \q */
        )) {
          if (state.eat(
            123
            /* { */
          )) {
            var result = this.regexp_classStringDisjunctionContents(state);
            if (state.eat(
              125
              /* } */
            )) {
              return result;
            }
          } else {
            state.raise("Invalid escape");
          }
          state.pos = start2;
        }
        return null;
      };
      pp$1.regexp_classStringDisjunctionContents = function(state) {
        var result = this.regexp_classString(state);
        while (state.eat(
          124
          /* | */
        )) {
          if (this.regexp_classString(state) === CharSetString) {
            result = CharSetString;
          }
        }
        return result;
      };
      pp$1.regexp_classString = function(state) {
        var count3 = 0;
        while (this.regexp_eatClassSetCharacter(state)) {
          count3++;
        }
        return count3 === 1 ? CharSetOk : CharSetString;
      };
      pp$1.regexp_eatClassSetCharacter = function(state) {
        var start2 = state.pos;
        if (state.eat(
          92
          /* \ */
        )) {
          if (this.regexp_eatCharacterEscape(state) || this.regexp_eatClassSetReservedPunctuator(state)) {
            return true;
          }
          if (state.eat(
            98
            /* b */
          )) {
            state.lastIntValue = 8;
            return true;
          }
          state.pos = start2;
          return false;
        }
        var ch = state.current();
        if (ch < 0 || ch === state.lookahead() && isClassSetReservedDoublePunctuatorCharacter(ch)) {
          return false;
        }
        if (isClassSetSyntaxCharacter(ch)) {
          return false;
        }
        state.advance();
        state.lastIntValue = ch;
        return true;
      };
      pp$1.regexp_eatClassSetReservedPunctuator = function(state) {
        var ch = state.current();
        if (isClassSetReservedPunctuator(ch)) {
          state.lastIntValue = ch;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatClassControlLetter = function(state) {
        var ch = state.current();
        if (isDecimalDigit(ch) || ch === 95) {
          state.lastIntValue = ch % 32;
          state.advance();
          return true;
        }
        return false;
      };
      pp$1.regexp_eatHexEscapeSequence = function(state) {
        var start2 = state.pos;
        if (state.eat(
          120
          /* x */
        )) {
          if (this.regexp_eatFixedHexDigits(state, 2)) {
            return true;
          }
          if (state.switchU) {
            state.raise("Invalid escape");
          }
          state.pos = start2;
        }
        return false;
      };
      pp$1.regexp_eatDecimalDigits = function(state) {
        var start2 = state.pos;
        var ch = 0;
        state.lastIntValue = 0;
        while (isDecimalDigit(ch = state.current())) {
          state.lastIntValue = 10 * state.lastIntValue + (ch - 48);
          state.advance();
        }
        return state.pos !== start2;
      };
      pp$1.regexp_eatHexDigits = function(state) {
        var start2 = state.pos;
        var ch = 0;
        state.lastIntValue = 0;
        while (isHexDigit(ch = state.current())) {
          state.lastIntValue = 16 * state.lastIntValue + hexToInt(ch);
          state.advance();
        }
        return state.pos !== start2;
      };
      pp$1.regexp_eatLegacyOctalEscapeSequence = function(state) {
        if (this.regexp_eatOctalDigit(state)) {
          var n1 = state.lastIntValue;
          if (this.regexp_eatOctalDigit(state)) {
            var n2 = state.lastIntValue;
            if (n1 <= 3 && this.regexp_eatOctalDigit(state)) {
              state.lastIntValue = n1 * 64 + n2 * 8 + state.lastIntValue;
            } else {
              state.lastIntValue = n1 * 8 + n2;
            }
          } else {
            state.lastIntValue = n1;
          }
          return true;
        }
        return false;
      };
      pp$1.regexp_eatOctalDigit = function(state) {
        var ch = state.current();
        if (isOctalDigit(ch)) {
          state.lastIntValue = ch - 48;
          state.advance();
          return true;
        }
        state.lastIntValue = 0;
        return false;
      };
      pp$1.regexp_eatFixedHexDigits = function(state, length) {
        var start2 = state.pos;
        state.lastIntValue = 0;
        for (var i = 0; i < length; ++i) {
          var ch = state.current();
          if (!isHexDigit(ch)) {
            state.pos = start2;
            return false;
          }
          state.lastIntValue = 16 * state.lastIntValue + hexToInt(ch);
          state.advance();
        }
        return true;
      };
      Token = function Token2(p) {
        this.type = p.type;
        this.value = p.value;
        this.start = p.start;
        this.end = p.end;
        if (p.options.locations) {
          this.loc = new SourceLocation(p, p.startLoc, p.endLoc);
        }
        if (p.options.ranges) {
          this.range = [p.start, p.end];
        }
      };
      pp = Parser.prototype;
      pp.next = function(ignoreEscapeSequenceInKeyword) {
        if (!ignoreEscapeSequenceInKeyword && this.type.keyword && this.containsEsc) {
          this.raiseRecoverable(this.start, "Escape sequence in keyword " + this.type.keyword);
        }
        if (this.options.onToken) {
          this.options.onToken(new Token(this));
        }
        this.lastTokEnd = this.end;
        this.lastTokStart = this.start;
        this.lastTokEndLoc = this.endLoc;
        this.lastTokStartLoc = this.startLoc;
        this.nextToken();
      };
      pp.getToken = function() {
        this.next();
        return new Token(this);
      };
      if (typeof Symbol !== "undefined") {
        pp[Symbol.iterator] = function() {
          var this$1$1 = this;
          return {
            next: function() {
              var token = this$1$1.getToken();
              return {
                done: token.type === types$1.eof,
                value: token
              };
            }
          };
        };
      }
      pp.nextToken = function() {
        var curContext = this.curContext();
        if (!curContext || !curContext.preserveSpace) {
          this.skipSpace();
        }
        this.start = this.pos;
        if (this.options.locations) {
          this.startLoc = this.curPosition();
        }
        if (this.pos >= this.input.length) {
          return this.finishToken(types$1.eof);
        }
        if (curContext.override) {
          return curContext.override(this);
        } else {
          this.readToken(this.fullCharCodeAtPos());
        }
      };
      pp.readToken = function(code3) {
        if (isIdentifierStart(code3, this.options.ecmaVersion >= 6) || code3 === 92) {
          return this.readWord();
        }
        return this.getTokenFromCode(code3);
      };
      pp.fullCharCodeAtPos = function() {
        var code3 = this.input.charCodeAt(this.pos);
        if (code3 <= 55295 || code3 >= 56320) {
          return code3;
        }
        var next = this.input.charCodeAt(this.pos + 1);
        return next <= 56319 || next >= 57344 ? code3 : (code3 << 10) + next - 56613888;
      };
      pp.skipBlockComment = function() {
        var startLoc = this.options.onComment && this.curPosition();
        var start2 = this.pos, end = this.input.indexOf("*/", this.pos += 2);
        if (end === -1) {
          this.raise(this.pos - 2, "Unterminated comment");
        }
        this.pos = end + 2;
        if (this.options.locations) {
          for (var nextBreak = void 0, pos = start2; (nextBreak = nextLineBreak(this.input, pos, this.pos)) > -1; ) {
            ++this.curLine;
            pos = this.lineStart = nextBreak;
          }
        }
        if (this.options.onComment) {
          this.options.onComment(
            true,
            this.input.slice(start2 + 2, end),
            start2,
            this.pos,
            startLoc,
            this.curPosition()
          );
        }
      };
      pp.skipLineComment = function(startSkip) {
        var start2 = this.pos;
        var startLoc = this.options.onComment && this.curPosition();
        var ch = this.input.charCodeAt(this.pos += startSkip);
        while (this.pos < this.input.length && !isNewLine(ch)) {
          ch = this.input.charCodeAt(++this.pos);
        }
        if (this.options.onComment) {
          this.options.onComment(
            false,
            this.input.slice(start2 + startSkip, this.pos),
            start2,
            this.pos,
            startLoc,
            this.curPosition()
          );
        }
      };
      pp.skipSpace = function() {
        loop: while (this.pos < this.input.length) {
          var ch = this.input.charCodeAt(this.pos);
          switch (ch) {
            case 32:
            case 160:
              ++this.pos;
              break;
            case 13:
              if (this.input.charCodeAt(this.pos + 1) === 10) {
                ++this.pos;
              }
            case 10:
            case 8232:
            case 8233:
              ++this.pos;
              if (this.options.locations) {
                ++this.curLine;
                this.lineStart = this.pos;
              }
              break;
            case 47:
              switch (this.input.charCodeAt(this.pos + 1)) {
                case 42:
                  this.skipBlockComment();
                  break;
                case 47:
                  this.skipLineComment(2);
                  break;
                default:
                  break loop;
              }
              break;
            default:
              if (ch > 8 && ch < 14 || ch >= 5760 && nonASCIIwhitespace.test(String.fromCharCode(ch))) {
                ++this.pos;
              } else {
                break loop;
              }
          }
        }
      };
      pp.finishToken = function(type, val) {
        this.end = this.pos;
        if (this.options.locations) {
          this.endLoc = this.curPosition();
        }
        var prevType = this.type;
        this.type = type;
        this.value = val;
        this.updateContext(prevType);
      };
      pp.readToken_dot = function() {
        var next = this.input.charCodeAt(this.pos + 1);
        if (next >= 48 && next <= 57) {
          return this.readNumber(true);
        }
        var next2 = this.input.charCodeAt(this.pos + 2);
        if (this.options.ecmaVersion >= 6 && next === 46 && next2 === 46) {
          this.pos += 3;
          return this.finishToken(types$1.ellipsis);
        } else {
          ++this.pos;
          return this.finishToken(types$1.dot);
        }
      };
      pp.readToken_slash = function() {
        var next = this.input.charCodeAt(this.pos + 1);
        if (this.exprAllowed) {
          ++this.pos;
          return this.readRegexp();
        }
        if (next === 61) {
          return this.finishOp(types$1.assign, 2);
        }
        return this.finishOp(types$1.slash, 1);
      };
      pp.readToken_mult_modulo_exp = function(code3) {
        var next = this.input.charCodeAt(this.pos + 1);
        var size = 1;
        var tokentype = code3 === 42 ? types$1.star : types$1.modulo;
        if (this.options.ecmaVersion >= 7 && code3 === 42 && next === 42) {
          ++size;
          tokentype = types$1.starstar;
          next = this.input.charCodeAt(this.pos + 2);
        }
        if (next === 61) {
          return this.finishOp(types$1.assign, size + 1);
        }
        return this.finishOp(tokentype, size);
      };
      pp.readToken_pipe_amp = function(code3) {
        var next = this.input.charCodeAt(this.pos + 1);
        if (next === code3) {
          if (this.options.ecmaVersion >= 12) {
            var next2 = this.input.charCodeAt(this.pos + 2);
            if (next2 === 61) {
              return this.finishOp(types$1.assign, 3);
            }
          }
          return this.finishOp(code3 === 124 ? types$1.logicalOR : types$1.logicalAND, 2);
        }
        if (next === 61) {
          return this.finishOp(types$1.assign, 2);
        }
        return this.finishOp(code3 === 124 ? types$1.bitwiseOR : types$1.bitwiseAND, 1);
      };
      pp.readToken_caret = function() {
        var next = this.input.charCodeAt(this.pos + 1);
        if (next === 61) {
          return this.finishOp(types$1.assign, 2);
        }
        return this.finishOp(types$1.bitwiseXOR, 1);
      };
      pp.readToken_plus_min = function(code3) {
        var next = this.input.charCodeAt(this.pos + 1);
        if (next === code3) {
          if (next === 45 && !this.inModule && this.input.charCodeAt(this.pos + 2) === 62 && (this.lastTokEnd === 0 || lineBreak.test(this.input.slice(this.lastTokEnd, this.pos)))) {
            this.skipLineComment(3);
            this.skipSpace();
            return this.nextToken();
          }
          return this.finishOp(types$1.incDec, 2);
        }
        if (next === 61) {
          return this.finishOp(types$1.assign, 2);
        }
        return this.finishOp(types$1.plusMin, 1);
      };
      pp.readToken_lt_gt = function(code3) {
        var next = this.input.charCodeAt(this.pos + 1);
        var size = 1;
        if (next === code3) {
          size = code3 === 62 && this.input.charCodeAt(this.pos + 2) === 62 ? 3 : 2;
          if (this.input.charCodeAt(this.pos + size) === 61) {
            return this.finishOp(types$1.assign, size + 1);
          }
          return this.finishOp(types$1.bitShift, size);
        }
        if (next === 33 && code3 === 60 && !this.inModule && this.input.charCodeAt(this.pos + 2) === 45 && this.input.charCodeAt(this.pos + 3) === 45) {
          this.skipLineComment(4);
          this.skipSpace();
          return this.nextToken();
        }
        if (next === 61) {
          size = 2;
        }
        return this.finishOp(types$1.relational, size);
      };
      pp.readToken_eq_excl = function(code3) {
        var next = this.input.charCodeAt(this.pos + 1);
        if (next === 61) {
          return this.finishOp(types$1.equality, this.input.charCodeAt(this.pos + 2) === 61 ? 3 : 2);
        }
        if (code3 === 61 && next === 62 && this.options.ecmaVersion >= 6) {
          this.pos += 2;
          return this.finishToken(types$1.arrow);
        }
        return this.finishOp(code3 === 61 ? types$1.eq : types$1.prefix, 1);
      };
      pp.readToken_question = function() {
        var ecmaVersion = this.options.ecmaVersion;
        if (ecmaVersion >= 11) {
          var next = this.input.charCodeAt(this.pos + 1);
          if (next === 46) {
            var next2 = this.input.charCodeAt(this.pos + 2);
            if (next2 < 48 || next2 > 57) {
              return this.finishOp(types$1.questionDot, 2);
            }
          }
          if (next === 63) {
            if (ecmaVersion >= 12) {
              var next2$1 = this.input.charCodeAt(this.pos + 2);
              if (next2$1 === 61) {
                return this.finishOp(types$1.assign, 3);
              }
            }
            return this.finishOp(types$1.coalesce, 2);
          }
        }
        return this.finishOp(types$1.question, 1);
      };
      pp.readToken_numberSign = function() {
        var ecmaVersion = this.options.ecmaVersion;
        var code3 = 35;
        if (ecmaVersion >= 13) {
          ++this.pos;
          code3 = this.fullCharCodeAtPos();
          if (isIdentifierStart(code3, true) || code3 === 92) {
            return this.finishToken(types$1.privateId, this.readWord1());
          }
        }
        this.raise(this.pos, "Unexpected character '" + codePointToString(code3) + "'");
      };
      pp.getTokenFromCode = function(code3) {
        switch (code3) {
          // The interpretation of a dot depends on whether it is followed
          // by a digit or another two dots.
          case 46:
            return this.readToken_dot();
          // Punctuation tokens.
          case 40:
            ++this.pos;
            return this.finishToken(types$1.parenL);
          case 41:
            ++this.pos;
            return this.finishToken(types$1.parenR);
          case 59:
            ++this.pos;
            return this.finishToken(types$1.semi);
          case 44:
            ++this.pos;
            return this.finishToken(types$1.comma);
          case 91:
            ++this.pos;
            return this.finishToken(types$1.bracketL);
          case 93:
            ++this.pos;
            return this.finishToken(types$1.bracketR);
          case 123:
            ++this.pos;
            return this.finishToken(types$1.braceL);
          case 125:
            ++this.pos;
            return this.finishToken(types$1.braceR);
          case 58:
            ++this.pos;
            return this.finishToken(types$1.colon);
          case 96:
            if (this.options.ecmaVersion < 6) {
              break;
            }
            ++this.pos;
            return this.finishToken(types$1.backQuote);
          case 48:
            var next = this.input.charCodeAt(this.pos + 1);
            if (next === 120 || next === 88) {
              return this.readRadixNumber(16);
            }
            if (this.options.ecmaVersion >= 6) {
              if (next === 111 || next === 79) {
                return this.readRadixNumber(8);
              }
              if (next === 98 || next === 66) {
                return this.readRadixNumber(2);
              }
            }
          // Anything else beginning with a digit is an integer, octal
          // number, or float.
          case 49:
          case 50:
          case 51:
          case 52:
          case 53:
          case 54:
          case 55:
          case 56:
          case 57:
            return this.readNumber(false);
          // Quotes produce strings.
          case 34:
          case 39:
            return this.readString(code3);
          // Operators are parsed inline in tiny state machines. '=' (61) is
          // often referred to. `finishOp` simply skips the amount of
          // characters it is given as second argument, and returns a token
          // of the type given by its first argument.
          case 47:
            return this.readToken_slash();
          case 37:
          case 42:
            return this.readToken_mult_modulo_exp(code3);
          case 124:
          case 38:
            return this.readToken_pipe_amp(code3);
          case 94:
            return this.readToken_caret();
          case 43:
          case 45:
            return this.readToken_plus_min(code3);
          case 60:
          case 62:
            return this.readToken_lt_gt(code3);
          case 61:
          case 33:
            return this.readToken_eq_excl(code3);
          case 63:
            return this.readToken_question();
          case 126:
            return this.finishOp(types$1.prefix, 1);
          case 35:
            return this.readToken_numberSign();
        }
        this.raise(this.pos, "Unexpected character '" + codePointToString(code3) + "'");
      };
      pp.finishOp = function(type, size) {
        var str = this.input.slice(this.pos, this.pos + size);
        this.pos += size;
        return this.finishToken(type, str);
      };
      pp.readRegexp = function() {
        var escaped, inClass, start2 = this.pos;
        for (; ; ) {
          if (this.pos >= this.input.length) {
            this.raise(start2, "Unterminated regular expression");
          }
          var ch = this.input.charAt(this.pos);
          if (lineBreak.test(ch)) {
            this.raise(start2, "Unterminated regular expression");
          }
          if (!escaped) {
            if (ch === "[") {
              inClass = true;
            } else if (ch === "]" && inClass) {
              inClass = false;
            } else if (ch === "/" && !inClass) {
              break;
            }
            escaped = ch === "\\";
          } else {
            escaped = false;
          }
          ++this.pos;
        }
        var pattern = this.input.slice(start2, this.pos);
        ++this.pos;
        var flagsStart = this.pos;
        var flags = this.readWord1();
        if (this.containsEsc) {
          this.unexpected(flagsStart);
        }
        var state = this.regexpState || (this.regexpState = new RegExpValidationState(this));
        state.reset(start2, pattern, flags);
        this.validateRegExpFlags(state);
        this.validateRegExpPattern(state);
        var value = null;
        try {
          value = new RegExp(pattern, flags);
        } catch (e) {
        }
        return this.finishToken(types$1.regexp, { pattern, flags, value });
      };
      pp.readInt = function(radix, len, maybeLegacyOctalNumericLiteral) {
        var allowSeparators = this.options.ecmaVersion >= 12 && len === void 0;
        var isLegacyOctalNumericLiteral = maybeLegacyOctalNumericLiteral && this.input.charCodeAt(this.pos) === 48;
        var start2 = this.pos, total = 0, lastCode = 0;
        for (var i = 0, e = len == null ? Infinity : len; i < e; ++i, ++this.pos) {
          var code3 = this.input.charCodeAt(this.pos), val = void 0;
          if (allowSeparators && code3 === 95) {
            if (isLegacyOctalNumericLiteral) {
              this.raiseRecoverable(this.pos, "Numeric separator is not allowed in legacy octal numeric literals");
            }
            if (lastCode === 95) {
              this.raiseRecoverable(this.pos, "Numeric separator must be exactly one underscore");
            }
            if (i === 0) {
              this.raiseRecoverable(this.pos, "Numeric separator is not allowed at the first of digits");
            }
            lastCode = code3;
            continue;
          }
          if (code3 >= 97) {
            val = code3 - 97 + 10;
          } else if (code3 >= 65) {
            val = code3 - 65 + 10;
          } else if (code3 >= 48 && code3 <= 57) {
            val = code3 - 48;
          } else {
            val = Infinity;
          }
          if (val >= radix) {
            break;
          }
          lastCode = code3;
          total = total * radix + val;
        }
        if (allowSeparators && lastCode === 95) {
          this.raiseRecoverable(this.pos - 1, "Numeric separator is not allowed at the last of digits");
        }
        if (this.pos === start2 || len != null && this.pos - start2 !== len) {
          return null;
        }
        return total;
      };
      pp.readRadixNumber = function(radix) {
        var start2 = this.pos;
        this.pos += 2;
        var val = this.readInt(radix);
        if (val == null) {
          this.raise(this.start + 2, "Expected number in radix " + radix);
        }
        if (this.options.ecmaVersion >= 11 && this.input.charCodeAt(this.pos) === 110) {
          val = stringToBigInt(this.input.slice(start2, this.pos));
          ++this.pos;
        } else if (isIdentifierStart(this.fullCharCodeAtPos())) {
          this.raise(this.pos, "Identifier directly after number");
        }
        return this.finishToken(types$1.num, val);
      };
      pp.readNumber = function(startsWithDot) {
        var start2 = this.pos;
        if (!startsWithDot && this.readInt(10, void 0, true) === null) {
          this.raise(start2, "Invalid number");
        }
        var octal = this.pos - start2 >= 2 && this.input.charCodeAt(start2) === 48;
        if (octal && this.strict) {
          this.raise(start2, "Invalid number");
        }
        var next = this.input.charCodeAt(this.pos);
        if (!octal && !startsWithDot && this.options.ecmaVersion >= 11 && next === 110) {
          var val$1 = stringToBigInt(this.input.slice(start2, this.pos));
          ++this.pos;
          if (isIdentifierStart(this.fullCharCodeAtPos())) {
            this.raise(this.pos, "Identifier directly after number");
          }
          return this.finishToken(types$1.num, val$1);
        }
        if (octal && /[89]/.test(this.input.slice(start2, this.pos))) {
          octal = false;
        }
        if (next === 46 && !octal) {
          ++this.pos;
          this.readInt(10);
          next = this.input.charCodeAt(this.pos);
        }
        if ((next === 69 || next === 101) && !octal) {
          next = this.input.charCodeAt(++this.pos);
          if (next === 43 || next === 45) {
            ++this.pos;
          }
          if (this.readInt(10) === null) {
            this.raise(start2, "Invalid number");
          }
        }
        if (isIdentifierStart(this.fullCharCodeAtPos())) {
          this.raise(this.pos, "Identifier directly after number");
        }
        var val = stringToNumber(this.input.slice(start2, this.pos), octal);
        return this.finishToken(types$1.num, val);
      };
      pp.readCodePoint = function() {
        var ch = this.input.charCodeAt(this.pos), code3;
        if (ch === 123) {
          if (this.options.ecmaVersion < 6) {
            this.unexpected();
          }
          var codePos = ++this.pos;
          code3 = this.readHexChar(this.input.indexOf("}", this.pos) - this.pos);
          ++this.pos;
          if (code3 > 1114111) {
            this.invalidStringToken(codePos, "Code point out of bounds");
          }
        } else {
          code3 = this.readHexChar(4);
        }
        return code3;
      };
      pp.readString = function(quote) {
        var out = "", chunkStart = ++this.pos;
        for (; ; ) {
          if (this.pos >= this.input.length) {
            this.raise(this.start, "Unterminated string constant");
          }
          var ch = this.input.charCodeAt(this.pos);
          if (ch === quote) {
            break;
          }
          if (ch === 92) {
            out += this.input.slice(chunkStart, this.pos);
            out += this.readEscapedChar(false);
            chunkStart = this.pos;
          } else if (ch === 8232 || ch === 8233) {
            if (this.options.ecmaVersion < 10) {
              this.raise(this.start, "Unterminated string constant");
            }
            ++this.pos;
            if (this.options.locations) {
              this.curLine++;
              this.lineStart = this.pos;
            }
          } else {
            if (isNewLine(ch)) {
              this.raise(this.start, "Unterminated string constant");
            }
            ++this.pos;
          }
        }
        out += this.input.slice(chunkStart, this.pos++);
        return this.finishToken(types$1.string, out);
      };
      INVALID_TEMPLATE_ESCAPE_ERROR = {};
      pp.tryReadTemplateToken = function() {
        this.inTemplateElement = true;
        try {
          this.readTmplToken();
        } catch (err) {
          if (err === INVALID_TEMPLATE_ESCAPE_ERROR) {
            this.readInvalidTemplateToken();
          } else {
            throw err;
          }
        }
        this.inTemplateElement = false;
      };
      pp.invalidStringToken = function(position2, message) {
        if (this.inTemplateElement && this.options.ecmaVersion >= 9) {
          throw INVALID_TEMPLATE_ESCAPE_ERROR;
        } else {
          this.raise(position2, message);
        }
      };
      pp.readTmplToken = function() {
        var out = "", chunkStart = this.pos;
        for (; ; ) {
          if (this.pos >= this.input.length) {
            this.raise(this.start, "Unterminated template");
          }
          var ch = this.input.charCodeAt(this.pos);
          if (ch === 96 || ch === 36 && this.input.charCodeAt(this.pos + 1) === 123) {
            if (this.pos === this.start && (this.type === types$1.template || this.type === types$1.invalidTemplate)) {
              if (ch === 36) {
                this.pos += 2;
                return this.finishToken(types$1.dollarBraceL);
              } else {
                ++this.pos;
                return this.finishToken(types$1.backQuote);
              }
            }
            out += this.input.slice(chunkStart, this.pos);
            return this.finishToken(types$1.template, out);
          }
          if (ch === 92) {
            out += this.input.slice(chunkStart, this.pos);
            out += this.readEscapedChar(true);
            chunkStart = this.pos;
          } else if (isNewLine(ch)) {
            out += this.input.slice(chunkStart, this.pos);
            ++this.pos;
            switch (ch) {
              case 13:
                if (this.input.charCodeAt(this.pos) === 10) {
                  ++this.pos;
                }
              case 10:
                out += "\n";
                break;
              default:
                out += String.fromCharCode(ch);
                break;
            }
            if (this.options.locations) {
              ++this.curLine;
              this.lineStart = this.pos;
            }
            chunkStart = this.pos;
          } else {
            ++this.pos;
          }
        }
      };
      pp.readInvalidTemplateToken = function() {
        for (; this.pos < this.input.length; this.pos++) {
          switch (this.input[this.pos]) {
            case "\\":
              ++this.pos;
              break;
            case "$":
              if (this.input[this.pos + 1] !== "{") {
                break;
              }
            // fall through
            case "`":
              return this.finishToken(types$1.invalidTemplate, this.input.slice(this.start, this.pos));
            case "\r":
              if (this.input[this.pos + 1] === "\n") {
                ++this.pos;
              }
            // fall through
            case "\n":
            case "\u2028":
            case "\u2029":
              ++this.curLine;
              this.lineStart = this.pos + 1;
              break;
          }
        }
        this.raise(this.start, "Unterminated template");
      };
      pp.readEscapedChar = function(inTemplate) {
        var ch = this.input.charCodeAt(++this.pos);
        ++this.pos;
        switch (ch) {
          case 110:
            return "\n";
          // 'n' -> '\n'
          case 114:
            return "\r";
          // 'r' -> '\r'
          case 120:
            return String.fromCharCode(this.readHexChar(2));
          // 'x'
          case 117:
            return codePointToString(this.readCodePoint());
          // 'u'
          case 116:
            return "	";
          // 't' -> '\t'
          case 98:
            return "\b";
          // 'b' -> '\b'
          case 118:
            return "\v";
          // 'v' -> '\u000b'
          case 102:
            return "\f";
          // 'f' -> '\f'
          case 13:
            if (this.input.charCodeAt(this.pos) === 10) {
              ++this.pos;
            }
          // '\r\n'
          case 10:
            if (this.options.locations) {
              this.lineStart = this.pos;
              ++this.curLine;
            }
            return "";
          case 56:
          case 57:
            if (this.strict) {
              this.invalidStringToken(
                this.pos - 1,
                "Invalid escape sequence"
              );
            }
            if (inTemplate) {
              var codePos = this.pos - 1;
              this.invalidStringToken(
                codePos,
                "Invalid escape sequence in template string"
              );
            }
          default:
            if (ch >= 48 && ch <= 55) {
              var octalStr = this.input.substr(this.pos - 1, 3).match(/^[0-7]+/)[0];
              var octal = parseInt(octalStr, 8);
              if (octal > 255) {
                octalStr = octalStr.slice(0, -1);
                octal = parseInt(octalStr, 8);
              }
              this.pos += octalStr.length - 1;
              ch = this.input.charCodeAt(this.pos);
              if ((octalStr !== "0" || ch === 56 || ch === 57) && (this.strict || inTemplate)) {
                this.invalidStringToken(
                  this.pos - 1 - octalStr.length,
                  inTemplate ? "Octal literal in template string" : "Octal literal in strict mode"
                );
              }
              return String.fromCharCode(octal);
            }
            if (isNewLine(ch)) {
              if (this.options.locations) {
                this.lineStart = this.pos;
                ++this.curLine;
              }
              return "";
            }
            return String.fromCharCode(ch);
        }
      };
      pp.readHexChar = function(len) {
        var codePos = this.pos;
        var n = this.readInt(16, len);
        if (n === null) {
          this.invalidStringToken(codePos, "Bad character escape sequence");
        }
        return n;
      };
      pp.readWord1 = function() {
        this.containsEsc = false;
        var word = "", first = true, chunkStart = this.pos;
        var astral = this.options.ecmaVersion >= 6;
        while (this.pos < this.input.length) {
          var ch = this.fullCharCodeAtPos();
          if (isIdentifierChar(ch, astral)) {
            this.pos += ch <= 65535 ? 1 : 2;
          } else if (ch === 92) {
            this.containsEsc = true;
            word += this.input.slice(chunkStart, this.pos);
            var escStart = this.pos;
            if (this.input.charCodeAt(++this.pos) !== 117) {
              this.invalidStringToken(this.pos, "Expecting Unicode escape sequence \\uXXXX");
            }
            ++this.pos;
            var esc = this.readCodePoint();
            if (!(first ? isIdentifierStart : isIdentifierChar)(esc, astral)) {
              this.invalidStringToken(escStart, "Invalid Unicode escape");
            }
            word += codePointToString(esc);
            chunkStart = this.pos;
          } else {
            break;
          }
          first = false;
        }
        return word + this.input.slice(chunkStart, this.pos);
      };
      pp.readWord = function() {
        var word = this.readWord1();
        var type = types$1.name;
        if (this.keywords.test(word)) {
          type = keywords[word];
        }
        return this.finishToken(type, word);
      };
      version = "8.15.0";
      Parser.acorn = {
        Parser,
        version,
        defaultOptions,
        Position,
        SourceLocation,
        getLineInfo,
        Node,
        TokenType,
        tokTypes: types$1,
        keywordTypes: keywords,
        TokContext,
        tokContexts: types,
        isIdentifierChar,
        isIdentifierStart,
        Token,
        isNewLine,
        lineBreak,
        lineBreakG,
        nonASCIIwhitespace
      };
    }
  });

  // ../../../ui/node_modules/.pnpm/acorn-jsx@5.3.2_acorn@8.15.0/node_modules/acorn-jsx/xhtml.js
  var require_xhtml = __commonJS({
    "../../../ui/node_modules/.pnpm/acorn-jsx@5.3.2_acorn@8.15.0/node_modules/acorn-jsx/xhtml.js"(exports, module) {
      module.exports = {
        quot: '"',
        amp: "&",
        apos: "'",
        lt: "<",
        gt: ">",
        nbsp: " ",
        iexcl: "¡",
        cent: "¢",
        pound: "£",
        curren: "¤",
        yen: "¥",
        brvbar: "¦",
        sect: "§",
        uml: "¨",
        copy: "©",
        ordf: "ª",
        laquo: "«",
        not: "¬",
        shy: "­",
        reg: "®",
        macr: "¯",
        deg: "°",
        plusmn: "±",
        sup2: "²",
        sup3: "³",
        acute: "´",
        micro: "µ",
        para: "¶",
        middot: "·",
        cedil: "¸",
        sup1: "¹",
        ordm: "º",
        raquo: "»",
        frac14: "¼",
        frac12: "½",
        frac34: "¾",
        iquest: "¿",
        Agrave: "À",
        Aacute: "Á",
        Acirc: "Â",
        Atilde: "Ã",
        Auml: "Ä",
        Aring: "Å",
        AElig: "Æ",
        Ccedil: "Ç",
        Egrave: "È",
        Eacute: "É",
        Ecirc: "Ê",
        Euml: "Ë",
        Igrave: "Ì",
        Iacute: "Í",
        Icirc: "Î",
        Iuml: "Ï",
        ETH: "Ð",
        Ntilde: "Ñ",
        Ograve: "Ò",
        Oacute: "Ó",
        Ocirc: "Ô",
        Otilde: "Õ",
        Ouml: "Ö",
        times: "×",
        Oslash: "Ø",
        Ugrave: "Ù",
        Uacute: "Ú",
        Ucirc: "Û",
        Uuml: "Ü",
        Yacute: "Ý",
        THORN: "Þ",
        szlig: "ß",
        agrave: "à",
        aacute: "á",
        acirc: "â",
        atilde: "ã",
        auml: "ä",
        aring: "å",
        aelig: "æ",
        ccedil: "ç",
        egrave: "è",
        eacute: "é",
        ecirc: "ê",
        euml: "ë",
        igrave: "ì",
        iacute: "í",
        icirc: "î",
        iuml: "ï",
        eth: "ð",
        ntilde: "ñ",
        ograve: "ò",
        oacute: "ó",
        ocirc: "ô",
        otilde: "õ",
        ouml: "ö",
        divide: "÷",
        oslash: "ø",
        ugrave: "ù",
        uacute: "ú",
        ucirc: "û",
        uuml: "ü",
        yacute: "ý",
        thorn: "þ",
        yuml: "ÿ",
        OElig: "Œ",
        oelig: "œ",
        Scaron: "Š",
        scaron: "š",
        Yuml: "Ÿ",
        fnof: "ƒ",
        circ: "ˆ",
        tilde: "˜",
        Alpha: "Α",
        Beta: "Β",
        Gamma: "Γ",
        Delta: "Δ",
        Epsilon: "Ε",
        Zeta: "Ζ",
        Eta: "Η",
        Theta: "Θ",
        Iota: "Ι",
        Kappa: "Κ",
        Lambda: "Λ",
        Mu: "Μ",
        Nu: "Ν",
        Xi: "Ξ",
        Omicron: "Ο",
        Pi: "Π",
        Rho: "Ρ",
        Sigma: "Σ",
        Tau: "Τ",
        Upsilon: "Υ",
        Phi: "Φ",
        Chi: "Χ",
        Psi: "Ψ",
        Omega: "Ω",
        alpha: "α",
        beta: "β",
        gamma: "γ",
        delta: "δ",
        epsilon: "ε",
        zeta: "ζ",
        eta: "η",
        theta: "θ",
        iota: "ι",
        kappa: "κ",
        lambda: "λ",
        mu: "μ",
        nu: "ν",
        xi: "ξ",
        omicron: "ο",
        pi: "π",
        rho: "ρ",
        sigmaf: "ς",
        sigma: "σ",
        tau: "τ",
        upsilon: "υ",
        phi: "φ",
        chi: "χ",
        psi: "ψ",
        omega: "ω",
        thetasym: "ϑ",
        upsih: "ϒ",
        piv: "ϖ",
        ensp: " ",
        emsp: " ",
        thinsp: " ",
        zwnj: "‌",
        zwj: "‍",
        lrm: "‎",
        rlm: "‏",
        ndash: "–",
        mdash: "—",
        lsquo: "‘",
        rsquo: "’",
        sbquo: "‚",
        ldquo: "“",
        rdquo: "”",
        bdquo: "„",
        dagger: "†",
        Dagger: "‡",
        bull: "•",
        hellip: "…",
        permil: "‰",
        prime: "′",
        Prime: "″",
        lsaquo: "‹",
        rsaquo: "›",
        oline: "‾",
        frasl: "⁄",
        euro: "€",
        image: "ℑ",
        weierp: "℘",
        real: "ℜ",
        trade: "™",
        alefsym: "ℵ",
        larr: "←",
        uarr: "↑",
        rarr: "→",
        darr: "↓",
        harr: "↔",
        crarr: "↵",
        lArr: "⇐",
        uArr: "⇑",
        rArr: "⇒",
        dArr: "⇓",
        hArr: "⇔",
        forall: "∀",
        part: "∂",
        exist: "∃",
        empty: "∅",
        nabla: "∇",
        isin: "∈",
        notin: "∉",
        ni: "∋",
        prod: "∏",
        sum: "∑",
        minus: "−",
        lowast: "∗",
        radic: "√",
        prop: "∝",
        infin: "∞",
        ang: "∠",
        and: "∧",
        or: "∨",
        cap: "∩",
        cup: "∪",
        "int": "∫",
        there4: "∴",
        sim: "∼",
        cong: "≅",
        asymp: "≈",
        ne: "≠",
        equiv: "≡",
        le: "≤",
        ge: "≥",
        sub: "⊂",
        sup: "⊃",
        nsub: "⊄",
        sube: "⊆",
        supe: "⊇",
        oplus: "⊕",
        otimes: "⊗",
        perp: "⊥",
        sdot: "⋅",
        lceil: "⌈",
        rceil: "⌉",
        lfloor: "⌊",
        rfloor: "⌋",
        lang: "〈",
        rang: "〉",
        loz: "◊",
        spades: "♠",
        clubs: "♣",
        hearts: "♥",
        diams: "♦"
      };
    }
  });

  // ../../../ui/node_modules/.pnpm/acorn-jsx@5.3.2_acorn@8.15.0/node_modules/acorn-jsx/index.js
  var require_acorn_jsx = __commonJS({
    "../../../ui/node_modules/.pnpm/acorn-jsx@5.3.2_acorn@8.15.0/node_modules/acorn-jsx/index.js"(exports, module) {
      "use strict";
      var XHTMLEntities = require_xhtml();
      var hexNumber = /^[\da-fA-F]+$/;
      var decimalNumber = /^\d+$/;
      var acornJsxMap = /* @__PURE__ */ new WeakMap();
      function getJsxTokens(acorn) {
        acorn = acorn.Parser.acorn || acorn;
        let acornJsx2 = acornJsxMap.get(acorn);
        if (!acornJsx2) {
          const tt = acorn.tokTypes;
          const TokContext3 = acorn.TokContext;
          const TokenType3 = acorn.TokenType;
          const tc_oTag = new TokContext3("<tag", false);
          const tc_cTag = new TokContext3("</tag", false);
          const tc_expr = new TokContext3("<tag>...</tag>", true, true);
          const tokContexts = {
            tc_oTag,
            tc_cTag,
            tc_expr
          };
          const tokTypes = {
            jsxName: new TokenType3("jsxName"),
            jsxText: new TokenType3("jsxText", { beforeExpr: true }),
            jsxTagStart: new TokenType3("jsxTagStart", { startsExpr: true }),
            jsxTagEnd: new TokenType3("jsxTagEnd")
          };
          tokTypes.jsxTagStart.updateContext = function() {
            this.context.push(tc_expr);
            this.context.push(tc_oTag);
            this.exprAllowed = false;
          };
          tokTypes.jsxTagEnd.updateContext = function(prevType) {
            let out = this.context.pop();
            if (out === tc_oTag && prevType === tt.slash || out === tc_cTag) {
              this.context.pop();
              this.exprAllowed = this.curContext() === tc_expr;
            } else {
              this.exprAllowed = true;
            }
          };
          acornJsx2 = { tokContexts, tokTypes };
          acornJsxMap.set(acorn, acornJsx2);
        }
        return acornJsx2;
      }
      function getQualifiedJSXName(object) {
        if (!object)
          return object;
        if (object.type === "JSXIdentifier")
          return object.name;
        if (object.type === "JSXNamespacedName")
          return object.namespace.name + ":" + object.name.name;
        if (object.type === "JSXMemberExpression")
          return getQualifiedJSXName(object.object) + "." + getQualifiedJSXName(object.property);
      }
      module.exports = function(options) {
        options = options || {};
        return function(Parser3) {
          return plugin({
            allowNamespaces: options.allowNamespaces !== false,
            allowNamespacedObjects: !!options.allowNamespacedObjects
          }, Parser3);
        };
      };
      Object.defineProperty(module.exports, "tokTypes", {
        get: function get_tokTypes() {
          return getJsxTokens((init_acorn(), __toCommonJS(acorn_exports))).tokTypes;
        },
        configurable: true,
        enumerable: true
      });
      function plugin(options, Parser3) {
        const acorn = Parser3.acorn || (init_acorn(), __toCommonJS(acorn_exports));
        const acornJsx2 = getJsxTokens(acorn);
        const tt = acorn.tokTypes;
        const tok = acornJsx2.tokTypes;
        const tokContexts = acorn.tokContexts;
        const tc_oTag = acornJsx2.tokContexts.tc_oTag;
        const tc_cTag = acornJsx2.tokContexts.tc_cTag;
        const tc_expr = acornJsx2.tokContexts.tc_expr;
        const isNewLine2 = acorn.isNewLine;
        const isIdentifierStart2 = acorn.isIdentifierStart;
        const isIdentifierChar2 = acorn.isIdentifierChar;
        return class extends Parser3 {
          // Expose actual `tokTypes` and `tokContexts` to other plugins.
          static get acornJsx() {
            return acornJsx2;
          }
          // Reads inline JSX contents token.
          jsx_readToken() {
            let out = "", chunkStart = this.pos;
            for (; ; ) {
              if (this.pos >= this.input.length)
                this.raise(this.start, "Unterminated JSX contents");
              let ch = this.input.charCodeAt(this.pos);
              switch (ch) {
                case 60:
                // '<'
                case 123:
                  if (this.pos === this.start) {
                    if (ch === 60 && this.exprAllowed) {
                      ++this.pos;
                      return this.finishToken(tok.jsxTagStart);
                    }
                    return this.getTokenFromCode(ch);
                  }
                  out += this.input.slice(chunkStart, this.pos);
                  return this.finishToken(tok.jsxText, out);
                case 38:
                  out += this.input.slice(chunkStart, this.pos);
                  out += this.jsx_readEntity();
                  chunkStart = this.pos;
                  break;
                case 62:
                // '>'
                case 125:
                  this.raise(
                    this.pos,
                    "Unexpected token `" + this.input[this.pos] + "`. Did you mean `" + (ch === 62 ? "&gt;" : "&rbrace;") + '` or `{"' + this.input[this.pos] + '"}`?'
                  );
                default:
                  if (isNewLine2(ch)) {
                    out += this.input.slice(chunkStart, this.pos);
                    out += this.jsx_readNewLine(true);
                    chunkStart = this.pos;
                  } else {
                    ++this.pos;
                  }
              }
            }
          }
          jsx_readNewLine(normalizeCRLF) {
            let ch = this.input.charCodeAt(this.pos);
            let out;
            ++this.pos;
            if (ch === 13 && this.input.charCodeAt(this.pos) === 10) {
              ++this.pos;
              out = normalizeCRLF ? "\n" : "\r\n";
            } else {
              out = String.fromCharCode(ch);
            }
            if (this.options.locations) {
              ++this.curLine;
              this.lineStart = this.pos;
            }
            return out;
          }
          jsx_readString(quote) {
            let out = "", chunkStart = ++this.pos;
            for (; ; ) {
              if (this.pos >= this.input.length)
                this.raise(this.start, "Unterminated string constant");
              let ch = this.input.charCodeAt(this.pos);
              if (ch === quote) break;
              if (ch === 38) {
                out += this.input.slice(chunkStart, this.pos);
                out += this.jsx_readEntity();
                chunkStart = this.pos;
              } else if (isNewLine2(ch)) {
                out += this.input.slice(chunkStart, this.pos);
                out += this.jsx_readNewLine(false);
                chunkStart = this.pos;
              } else {
                ++this.pos;
              }
            }
            out += this.input.slice(chunkStart, this.pos++);
            return this.finishToken(tt.string, out);
          }
          jsx_readEntity() {
            let str = "", count3 = 0, entity;
            let ch = this.input[this.pos];
            if (ch !== "&")
              this.raise(this.pos, "Entity must start with an ampersand");
            let startPos = ++this.pos;
            while (this.pos < this.input.length && count3++ < 10) {
              ch = this.input[this.pos++];
              if (ch === ";") {
                if (str[0] === "#") {
                  if (str[1] === "x") {
                    str = str.substr(2);
                    if (hexNumber.test(str))
                      entity = String.fromCharCode(parseInt(str, 16));
                  } else {
                    str = str.substr(1);
                    if (decimalNumber.test(str))
                      entity = String.fromCharCode(parseInt(str, 10));
                  }
                } else {
                  entity = XHTMLEntities[str];
                }
                break;
              }
              str += ch;
            }
            if (!entity) {
              this.pos = startPos;
              return "&";
            }
            return entity;
          }
          // Read a JSX identifier (valid tag or attribute name).
          //
          // Optimized version since JSX identifiers can't contain
          // escape characters and so can be read as single slice.
          // Also assumes that first character was already checked
          // by isIdentifierStart in readToken.
          jsx_readWord() {
            let ch, start2 = this.pos;
            do {
              ch = this.input.charCodeAt(++this.pos);
            } while (isIdentifierChar2(ch) || ch === 45);
            return this.finishToken(tok.jsxName, this.input.slice(start2, this.pos));
          }
          // Parse next token as JSX identifier
          jsx_parseIdentifier() {
            let node2 = this.startNode();
            if (this.type === tok.jsxName)
              node2.name = this.value;
            else if (this.type.keyword)
              node2.name = this.type.keyword;
            else
              this.unexpected();
            this.next();
            return this.finishNode(node2, "JSXIdentifier");
          }
          // Parse namespaced identifier.
          jsx_parseNamespacedName() {
            let startPos = this.start, startLoc = this.startLoc;
            let name2 = this.jsx_parseIdentifier();
            if (!options.allowNamespaces || !this.eat(tt.colon)) return name2;
            var node2 = this.startNodeAt(startPos, startLoc);
            node2.namespace = name2;
            node2.name = this.jsx_parseIdentifier();
            return this.finishNode(node2, "JSXNamespacedName");
          }
          // Parses element name in any form - namespaced, member
          // or single identifier.
          jsx_parseElementName() {
            if (this.type === tok.jsxTagEnd) return "";
            let startPos = this.start, startLoc = this.startLoc;
            let node2 = this.jsx_parseNamespacedName();
            if (this.type === tt.dot && node2.type === "JSXNamespacedName" && !options.allowNamespacedObjects) {
              this.unexpected();
            }
            while (this.eat(tt.dot)) {
              let newNode = this.startNodeAt(startPos, startLoc);
              newNode.object = node2;
              newNode.property = this.jsx_parseIdentifier();
              node2 = this.finishNode(newNode, "JSXMemberExpression");
            }
            return node2;
          }
          // Parses any type of JSX attribute value.
          jsx_parseAttributeValue() {
            switch (this.type) {
              case tt.braceL:
                let node2 = this.jsx_parseExpressionContainer();
                if (node2.expression.type === "JSXEmptyExpression")
                  this.raise(node2.start, "JSX attributes must only be assigned a non-empty expression");
                return node2;
              case tok.jsxTagStart:
              case tt.string:
                return this.parseExprAtom();
              default:
                this.raise(this.start, "JSX value should be either an expression or a quoted JSX text");
            }
          }
          // JSXEmptyExpression is unique type since it doesn't actually parse anything,
          // and so it should start at the end of last read token (left brace) and finish
          // at the beginning of the next one (right brace).
          jsx_parseEmptyExpression() {
            let node2 = this.startNodeAt(this.lastTokEnd, this.lastTokEndLoc);
            return this.finishNodeAt(node2, "JSXEmptyExpression", this.start, this.startLoc);
          }
          // Parses JSX expression enclosed into curly brackets.
          jsx_parseExpressionContainer() {
            let node2 = this.startNode();
            this.next();
            node2.expression = this.type === tt.braceR ? this.jsx_parseEmptyExpression() : this.parseExpression();
            this.expect(tt.braceR);
            return this.finishNode(node2, "JSXExpressionContainer");
          }
          // Parses following JSX attribute name-value pair.
          jsx_parseAttribute() {
            let node2 = this.startNode();
            if (this.eat(tt.braceL)) {
              this.expect(tt.ellipsis);
              node2.argument = this.parseMaybeAssign();
              this.expect(tt.braceR);
              return this.finishNode(node2, "JSXSpreadAttribute");
            }
            node2.name = this.jsx_parseNamespacedName();
            node2.value = this.eat(tt.eq) ? this.jsx_parseAttributeValue() : null;
            return this.finishNode(node2, "JSXAttribute");
          }
          // Parses JSX opening tag starting after '<'.
          jsx_parseOpeningElementAt(startPos, startLoc) {
            let node2 = this.startNodeAt(startPos, startLoc);
            node2.attributes = [];
            let nodeName = this.jsx_parseElementName();
            if (nodeName) node2.name = nodeName;
            while (this.type !== tt.slash && this.type !== tok.jsxTagEnd)
              node2.attributes.push(this.jsx_parseAttribute());
            node2.selfClosing = this.eat(tt.slash);
            this.expect(tok.jsxTagEnd);
            return this.finishNode(node2, nodeName ? "JSXOpeningElement" : "JSXOpeningFragment");
          }
          // Parses JSX closing tag starting after '</'.
          jsx_parseClosingElementAt(startPos, startLoc) {
            let node2 = this.startNodeAt(startPos, startLoc);
            let nodeName = this.jsx_parseElementName();
            if (nodeName) node2.name = nodeName;
            this.expect(tok.jsxTagEnd);
            return this.finishNode(node2, nodeName ? "JSXClosingElement" : "JSXClosingFragment");
          }
          // Parses entire JSX element, including it's opening tag
          // (starting after '<'), attributes, contents and closing tag.
          jsx_parseElementAt(startPos, startLoc) {
            let node2 = this.startNodeAt(startPos, startLoc);
            let children = [];
            let openingElement = this.jsx_parseOpeningElementAt(startPos, startLoc);
            let closingElement = null;
            if (!openingElement.selfClosing) {
              contents: for (; ; ) {
                switch (this.type) {
                  case tok.jsxTagStart:
                    startPos = this.start;
                    startLoc = this.startLoc;
                    this.next();
                    if (this.eat(tt.slash)) {
                      closingElement = this.jsx_parseClosingElementAt(startPos, startLoc);
                      break contents;
                    }
                    children.push(this.jsx_parseElementAt(startPos, startLoc));
                    break;
                  case tok.jsxText:
                    children.push(this.parseExprAtom());
                    break;
                  case tt.braceL:
                    children.push(this.jsx_parseExpressionContainer());
                    break;
                  default:
                    this.unexpected();
                }
              }
              if (getQualifiedJSXName(closingElement.name) !== getQualifiedJSXName(openingElement.name)) {
                this.raise(
                  closingElement.start,
                  "Expected corresponding JSX closing tag for <" + getQualifiedJSXName(openingElement.name) + ">"
                );
              }
            }
            let fragmentOrElement = openingElement.name ? "Element" : "Fragment";
            node2["opening" + fragmentOrElement] = openingElement;
            node2["closing" + fragmentOrElement] = closingElement;
            node2.children = children;
            if (this.type === tt.relational && this.value === "<") {
              this.raise(this.start, "Adjacent JSX elements must be wrapped in an enclosing tag");
            }
            return this.finishNode(node2, "JSX" + fragmentOrElement);
          }
          // Parse JSX text
          jsx_parseText() {
            let node2 = this.parseLiteral(this.value);
            node2.type = "JSXText";
            return node2;
          }
          // Parses entire JSX element from current position.
          jsx_parseElement() {
            let startPos = this.start, startLoc = this.startLoc;
            this.next();
            return this.jsx_parseElementAt(startPos, startLoc);
          }
          parseExprAtom(refShortHandDefaultPos) {
            if (this.type === tok.jsxText)
              return this.jsx_parseText();
            else if (this.type === tok.jsxTagStart)
              return this.jsx_parseElement();
            else
              return super.parseExprAtom(refShortHandDefaultPos);
          }
          readToken(code3) {
            let context = this.curContext();
            if (context === tc_expr) return this.jsx_readToken();
            if (context === tc_oTag || context === tc_cTag) {
              if (isIdentifierStart2(code3)) return this.jsx_readWord();
              if (code3 == 62) {
                ++this.pos;
                return this.finishToken(tok.jsxTagEnd);
              }
              if ((code3 === 34 || code3 === 39) && context == tc_oTag)
                return this.jsx_readString(code3);
            }
            if (code3 === 60 && this.exprAllowed && this.input.charCodeAt(this.pos + 1) !== 33) {
              ++this.pos;
              return this.finishToken(tok.jsxTagStart);
            }
            return super.readToken(code3);
          }
          updateContext(prevType) {
            if (this.type == tt.braceL) {
              var curContext = this.curContext();
              if (curContext == tc_oTag) this.context.push(tokContexts.b_expr);
              else if (curContext == tc_expr) this.context.push(tokContexts.b_tmpl);
              else super.updateContext(prevType);
              this.exprAllowed = true;
            } else if (this.type === tt.slash && prevType === tok.jsxTagStart) {
              this.context.length -= 2;
              this.context.push(tc_cTag);
              this.exprAllowed = false;
            } else {
              return super.updateContext(prevType);
            }
          }
        };
      }
    }
  });

  // ../../../ui/node_modules/.pnpm/bail@2.0.2/node_modules/bail/index.js
  function bail(error) {
    if (error) {
      throw error;
    }
  }

  // ../../../ui/node_modules/.pnpm/unified@11.0.5/node_modules/unified/lib/index.js
  var import_extend = __toESM(require_extend(), 1);

  // ../../../ui/node_modules/.pnpm/devlop@1.1.0/node_modules/devlop/lib/default.js
  function ok() {
  }

  // ../../../ui/node_modules/.pnpm/is-plain-obj@4.1.0/node_modules/is-plain-obj/index.js
  function isPlainObject(value) {
    if (typeof value !== "object" || value === null) {
      return false;
    }
    const prototype = Object.getPrototypeOf(value);
    return (prototype === null || prototype === Object.prototype || Object.getPrototypeOf(prototype) === null) && !(Symbol.toStringTag in value) && !(Symbol.iterator in value);
  }

  // ../../../ui/node_modules/.pnpm/trough@2.2.0/node_modules/trough/lib/index.js
  function trough() {
    const fns = [];
    const pipeline = { run, use };
    return pipeline;
    function run(...values) {
      let middlewareIndex = -1;
      const callback = values.pop();
      if (typeof callback !== "function") {
        throw new TypeError("Expected function as last argument, not " + callback);
      }
      next(null, ...values);
      function next(error, ...output) {
        const fn = fns[++middlewareIndex];
        let index2 = -1;
        if (error) {
          callback(error);
          return;
        }
        while (++index2 < values.length) {
          if (output[index2] === null || output[index2] === void 0) {
            output[index2] = values[index2];
          }
        }
        values = output;
        if (fn) {
          wrap(fn, next)(...output);
        } else {
          callback(null, ...output);
        }
      }
    }
    function use(middelware) {
      if (typeof middelware !== "function") {
        throw new TypeError(
          "Expected `middelware` to be a function, not " + middelware
        );
      }
      fns.push(middelware);
      return pipeline;
    }
  }
  function wrap(middleware, callback) {
    let called;
    return wrapped;
    function wrapped(...parameters) {
      const fnExpectsCallback = middleware.length > parameters.length;
      let result;
      if (fnExpectsCallback) {
        parameters.push(done);
      }
      try {
        result = middleware.apply(this, parameters);
      } catch (error) {
        const exception = (
          /** @type {Error} */
          error
        );
        if (fnExpectsCallback && called) {
          throw exception;
        }
        return done(exception);
      }
      if (!fnExpectsCallback) {
        if (result && result.then && typeof result.then === "function") {
          result.then(then, done);
        } else if (result instanceof Error) {
          done(result);
        } else {
          then(result);
        }
      }
    }
    function done(error, ...output) {
      if (!called) {
        called = true;
        callback(error, ...output);
      }
    }
    function then(value) {
      done(null, value);
    }
  }

  // ../../../ui/node_modules/.pnpm/unist-util-stringify-position@4.0.0/node_modules/unist-util-stringify-position/lib/index.js
  function stringifyPosition(value) {
    if (!value || typeof value !== "object") {
      return "";
    }
    if ("position" in value || "type" in value) {
      return position(value.position);
    }
    if ("start" in value || "end" in value) {
      return position(value);
    }
    if ("line" in value || "column" in value) {
      return point(value);
    }
    return "";
  }
  function point(point3) {
    return index(point3 && point3.line) + ":" + index(point3 && point3.column);
  }
  function position(pos) {
    return point(pos && pos.start) + "-" + point(pos && pos.end);
  }
  function index(value) {
    return value && typeof value === "number" ? value : 1;
  }

  // ../../../ui/node_modules/.pnpm/vfile-message@4.0.3/node_modules/vfile-message/lib/index.js
  var VFileMessage = class extends Error {
    /**
     * Create a message for `reason`.
     *
     * > 🪦 **Note**: also has obsolete signatures.
     *
     * @overload
     * @param {string} reason
     * @param {Options | null | undefined} [options]
     * @returns
     *
     * @overload
     * @param {string} reason
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @overload
     * @param {string} reason
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @overload
     * @param {string} reason
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {string | null | undefined} [origin]
     * @returns
     *
     * @param {Error | VFileMessage | string} causeOrReason
     *   Reason for message, should use markdown.
     * @param {Node | NodeLike | Options | Point | Position | string | null | undefined} [optionsOrParentOrPlace]
     *   Configuration (optional).
     * @param {string | null | undefined} [origin]
     *   Place in code where the message originates (example:
     *   `'my-package:my-rule'` or `'my-rule'`).
     * @returns
     *   Instance of `VFileMessage`.
     */
    // eslint-disable-next-line complexity
    constructor(causeOrReason, optionsOrParentOrPlace, origin) {
      super();
      if (typeof optionsOrParentOrPlace === "string") {
        origin = optionsOrParentOrPlace;
        optionsOrParentOrPlace = void 0;
      }
      let reason = "";
      let options = {};
      let legacyCause = false;
      if (optionsOrParentOrPlace) {
        if ("line" in optionsOrParentOrPlace && "column" in optionsOrParentOrPlace) {
          options = { place: optionsOrParentOrPlace };
        } else if ("start" in optionsOrParentOrPlace && "end" in optionsOrParentOrPlace) {
          options = { place: optionsOrParentOrPlace };
        } else if ("type" in optionsOrParentOrPlace) {
          options = {
            ancestors: [optionsOrParentOrPlace],
            place: optionsOrParentOrPlace.position
          };
        } else {
          options = { ...optionsOrParentOrPlace };
        }
      }
      if (typeof causeOrReason === "string") {
        reason = causeOrReason;
      } else if (!options.cause && causeOrReason) {
        legacyCause = true;
        reason = causeOrReason.message;
        options.cause = causeOrReason;
      }
      if (!options.ruleId && !options.source && typeof origin === "string") {
        const index2 = origin.indexOf(":");
        if (index2 === -1) {
          options.ruleId = origin;
        } else {
          options.source = origin.slice(0, index2);
          options.ruleId = origin.slice(index2 + 1);
        }
      }
      if (!options.place && options.ancestors && options.ancestors) {
        const parent = options.ancestors[options.ancestors.length - 1];
        if (parent) {
          options.place = parent.position;
        }
      }
      const start2 = options.place && "start" in options.place ? options.place.start : options.place;
      this.ancestors = options.ancestors || void 0;
      this.cause = options.cause || void 0;
      this.column = start2 ? start2.column : void 0;
      this.fatal = void 0;
      this.file = "";
      this.message = reason;
      this.line = start2 ? start2.line : void 0;
      this.name = stringifyPosition(options.place) || "1:1";
      this.place = options.place || void 0;
      this.reason = this.message;
      this.ruleId = options.ruleId || void 0;
      this.source = options.source || void 0;
      this.stack = legacyCause && options.cause && typeof options.cause.stack === "string" ? options.cause.stack : "";
      this.actual = void 0;
      this.expected = void 0;
      this.note = void 0;
      this.url = void 0;
    }
  };
  VFileMessage.prototype.file = "";
  VFileMessage.prototype.name = "";
  VFileMessage.prototype.reason = "";
  VFileMessage.prototype.message = "";
  VFileMessage.prototype.stack = "";
  VFileMessage.prototype.column = void 0;
  VFileMessage.prototype.line = void 0;
  VFileMessage.prototype.ancestors = void 0;
  VFileMessage.prototype.cause = void 0;
  VFileMessage.prototype.fatal = void 0;
  VFileMessage.prototype.place = void 0;
  VFileMessage.prototype.ruleId = void 0;
  VFileMessage.prototype.source = void 0;

  // ../../../ui/node_modules/.pnpm/vfile@6.0.3/node_modules/vfile/lib/minpath.browser.js
  var minpath = { basename, dirname, extname, join, sep: "/" };
  function basename(path2, extname3) {
    if (extname3 !== void 0 && typeof extname3 !== "string") {
      throw new TypeError('"ext" argument must be a string');
    }
    assertPath(path2);
    let start2 = 0;
    let end = -1;
    let index2 = path2.length;
    let seenNonSlash;
    if (extname3 === void 0 || extname3.length === 0 || extname3.length > path2.length) {
      while (index2--) {
        if (path2.codePointAt(index2) === 47) {
          if (seenNonSlash) {
            start2 = index2 + 1;
            break;
          }
        } else if (end < 0) {
          seenNonSlash = true;
          end = index2 + 1;
        }
      }
      return end < 0 ? "" : path2.slice(start2, end);
    }
    if (extname3 === path2) {
      return "";
    }
    let firstNonSlashEnd = -1;
    let extnameIndex = extname3.length - 1;
    while (index2--) {
      if (path2.codePointAt(index2) === 47) {
        if (seenNonSlash) {
          start2 = index2 + 1;
          break;
        }
      } else {
        if (firstNonSlashEnd < 0) {
          seenNonSlash = true;
          firstNonSlashEnd = index2 + 1;
        }
        if (extnameIndex > -1) {
          if (path2.codePointAt(index2) === extname3.codePointAt(extnameIndex--)) {
            if (extnameIndex < 0) {
              end = index2;
            }
          } else {
            extnameIndex = -1;
            end = firstNonSlashEnd;
          }
        }
      }
    }
    if (start2 === end) {
      end = firstNonSlashEnd;
    } else if (end < 0) {
      end = path2.length;
    }
    return path2.slice(start2, end);
  }
  function dirname(path2) {
    assertPath(path2);
    if (path2.length === 0) {
      return ".";
    }
    let end = -1;
    let index2 = path2.length;
    let unmatchedSlash;
    while (--index2) {
      if (path2.codePointAt(index2) === 47) {
        if (unmatchedSlash) {
          end = index2;
          break;
        }
      } else if (!unmatchedSlash) {
        unmatchedSlash = true;
      }
    }
    return end < 0 ? path2.codePointAt(0) === 47 ? "/" : "." : end === 1 && path2.codePointAt(0) === 47 ? "//" : path2.slice(0, end);
  }
  function extname(path2) {
    assertPath(path2);
    let index2 = path2.length;
    let end = -1;
    let startPart = 0;
    let startDot = -1;
    let preDotState = 0;
    let unmatchedSlash;
    while (index2--) {
      const code3 = path2.codePointAt(index2);
      if (code3 === 47) {
        if (unmatchedSlash) {
          startPart = index2 + 1;
          break;
        }
        continue;
      }
      if (end < 0) {
        unmatchedSlash = true;
        end = index2 + 1;
      }
      if (code3 === 46) {
        if (startDot < 0) {
          startDot = index2;
        } else if (preDotState !== 1) {
          preDotState = 1;
        }
      } else if (startDot > -1) {
        preDotState = -1;
      }
    }
    if (startDot < 0 || end < 0 || // We saw a non-dot character immediately before the dot.
    preDotState === 0 || // The (right-most) trimmed path component is exactly `..`.
    preDotState === 1 && startDot === end - 1 && startDot === startPart + 1) {
      return "";
    }
    return path2.slice(startDot, end);
  }
  function join(...segments) {
    let index2 = -1;
    let joined;
    while (++index2 < segments.length) {
      assertPath(segments[index2]);
      if (segments[index2]) {
        joined = joined === void 0 ? segments[index2] : joined + "/" + segments[index2];
      }
    }
    return joined === void 0 ? "." : normalize(joined);
  }
  function normalize(path2) {
    assertPath(path2);
    const absolute = path2.codePointAt(0) === 47;
    let value = normalizeString(path2, !absolute);
    if (value.length === 0 && !absolute) {
      value = ".";
    }
    if (value.length > 0 && path2.codePointAt(path2.length - 1) === 47) {
      value += "/";
    }
    return absolute ? "/" + value : value;
  }
  function normalizeString(path2, allowAboveRoot) {
    let result = "";
    let lastSegmentLength = 0;
    let lastSlash = -1;
    let dots = 0;
    let index2 = -1;
    let code3;
    let lastSlashIndex;
    while (++index2 <= path2.length) {
      if (index2 < path2.length) {
        code3 = path2.codePointAt(index2);
      } else if (code3 === 47) {
        break;
      } else {
        code3 = 47;
      }
      if (code3 === 47) {
        if (lastSlash === index2 - 1 || dots === 1) {
        } else if (lastSlash !== index2 - 1 && dots === 2) {
          if (result.length < 2 || lastSegmentLength !== 2 || result.codePointAt(result.length - 1) !== 46 || result.codePointAt(result.length - 2) !== 46) {
            if (result.length > 2) {
              lastSlashIndex = result.lastIndexOf("/");
              if (lastSlashIndex !== result.length - 1) {
                if (lastSlashIndex < 0) {
                  result = "";
                  lastSegmentLength = 0;
                } else {
                  result = result.slice(0, lastSlashIndex);
                  lastSegmentLength = result.length - 1 - result.lastIndexOf("/");
                }
                lastSlash = index2;
                dots = 0;
                continue;
              }
            } else if (result.length > 0) {
              result = "";
              lastSegmentLength = 0;
              lastSlash = index2;
              dots = 0;
              continue;
            }
          }
          if (allowAboveRoot) {
            result = result.length > 0 ? result + "/.." : "..";
            lastSegmentLength = 2;
          }
        } else {
          if (result.length > 0) {
            result += "/" + path2.slice(lastSlash + 1, index2);
          } else {
            result = path2.slice(lastSlash + 1, index2);
          }
          lastSegmentLength = index2 - lastSlash - 1;
        }
        lastSlash = index2;
        dots = 0;
      } else if (code3 === 46 && dots > -1) {
        dots++;
      } else {
        dots = -1;
      }
    }
    return result;
  }
  function assertPath(path2) {
    if (typeof path2 !== "string") {
      throw new TypeError(
        "Path must be a string. Received " + JSON.stringify(path2)
      );
    }
  }

  // ../../../ui/node_modules/.pnpm/vfile@6.0.3/node_modules/vfile/lib/minproc.browser.js
  var minproc = { cwd };
  function cwd() {
    return "/";
  }

  // ../../../ui/node_modules/.pnpm/vfile@6.0.3/node_modules/vfile/lib/minurl.shared.js
  function isUrl(fileUrlOrPath) {
    return Boolean(
      fileUrlOrPath !== null && typeof fileUrlOrPath === "object" && "href" in fileUrlOrPath && fileUrlOrPath.href && "protocol" in fileUrlOrPath && fileUrlOrPath.protocol && // @ts-expect-error: indexing is fine.
      fileUrlOrPath.auth === void 0
    );
  }

  // ../../../ui/node_modules/.pnpm/vfile@6.0.3/node_modules/vfile/lib/minurl.browser.js
  function urlToPath(path2) {
    if (typeof path2 === "string") {
      path2 = new URL(path2);
    } else if (!isUrl(path2)) {
      const error = new TypeError(
        'The "path" argument must be of type string or an instance of URL. Received `' + path2 + "`"
      );
      error.code = "ERR_INVALID_ARG_TYPE";
      throw error;
    }
    if (path2.protocol !== "file:") {
      const error = new TypeError("The URL must be of scheme file");
      error.code = "ERR_INVALID_URL_SCHEME";
      throw error;
    }
    return getPathFromURLPosix(path2);
  }
  function getPathFromURLPosix(url) {
    if (url.hostname !== "") {
      const error = new TypeError(
        'File URL host must be "localhost" or empty on darwin'
      );
      error.code = "ERR_INVALID_FILE_URL_HOST";
      throw error;
    }
    const pathname = url.pathname;
    let index2 = -1;
    while (++index2 < pathname.length) {
      if (pathname.codePointAt(index2) === 37 && pathname.codePointAt(index2 + 1) === 50) {
        const third = pathname.codePointAt(index2 + 2);
        if (third === 70 || third === 102) {
          const error = new TypeError(
            "File URL path must not include encoded / characters"
          );
          error.code = "ERR_INVALID_FILE_URL_PATH";
          throw error;
        }
      }
    }
    return decodeURIComponent(pathname);
  }

  // ../../../ui/node_modules/.pnpm/vfile@6.0.3/node_modules/vfile/lib/index.js
  var order = (
    /** @type {const} */
    [
      "history",
      "path",
      "basename",
      "stem",
      "extname",
      "dirname"
    ]
  );
  var VFile = class {
    /**
     * Create a new virtual file.
     *
     * `options` is treated as:
     *
     * *   `string` or `Uint8Array` — `{value: options}`
     * *   `URL` — `{path: options}`
     * *   `VFile` — shallow copies its data over to the new file
     * *   `object` — all fields are shallow copied over to the new file
     *
     * Path related fields are set in the following order (least specific to
     * most specific): `history`, `path`, `basename`, `stem`, `extname`,
     * `dirname`.
     *
     * You cannot set `dirname` or `extname` without setting either `history`,
     * `path`, `basename`, or `stem` too.
     *
     * @param {Compatible | null | undefined} [value]
     *   File value.
     * @returns
     *   New instance.
     */
    constructor(value) {
      let options;
      if (!value) {
        options = {};
      } else if (isUrl(value)) {
        options = { path: value };
      } else if (typeof value === "string" || isUint8Array(value)) {
        options = { value };
      } else {
        options = value;
      }
      this.cwd = "cwd" in options ? "" : minproc.cwd();
      this.data = {};
      this.history = [];
      this.messages = [];
      this.value;
      this.map;
      this.result;
      this.stored;
      let index2 = -1;
      while (++index2 < order.length) {
        const field2 = order[index2];
        if (field2 in options && options[field2] !== void 0 && options[field2] !== null) {
          this[field2] = field2 === "history" ? [...options[field2]] : options[field2];
        }
      }
      let field;
      for (field in options) {
        if (!order.includes(field)) {
          this[field] = options[field];
        }
      }
    }
    /**
     * Get the basename (including extname) (example: `'index.min.js'`).
     *
     * @returns {string | undefined}
     *   Basename.
     */
    get basename() {
      return typeof this.path === "string" ? minpath.basename(this.path) : void 0;
    }
    /**
     * Set basename (including extname) (`'index.min.js'`).
     *
     * Cannot contain path separators (`'/'` on unix, macOS, and browsers, `'\'`
     * on windows).
     * Cannot be nullified (use `file.path = file.dirname` instead).
     *
     * @param {string} basename
     *   Basename.
     * @returns {undefined}
     *   Nothing.
     */
    set basename(basename3) {
      assertNonEmpty(basename3, "basename");
      assertPart(basename3, "basename");
      this.path = minpath.join(this.dirname || "", basename3);
    }
    /**
     * Get the parent path (example: `'~'`).
     *
     * @returns {string | undefined}
     *   Dirname.
     */
    get dirname() {
      return typeof this.path === "string" ? minpath.dirname(this.path) : void 0;
    }
    /**
     * Set the parent path (example: `'~'`).
     *
     * Cannot be set if there’s no `path` yet.
     *
     * @param {string | undefined} dirname
     *   Dirname.
     * @returns {undefined}
     *   Nothing.
     */
    set dirname(dirname2) {
      assertPath2(this.basename, "dirname");
      this.path = minpath.join(dirname2 || "", this.basename);
    }
    /**
     * Get the extname (including dot) (example: `'.js'`).
     *
     * @returns {string | undefined}
     *   Extname.
     */
    get extname() {
      return typeof this.path === "string" ? minpath.extname(this.path) : void 0;
    }
    /**
     * Set the extname (including dot) (example: `'.js'`).
     *
     * Cannot contain path separators (`'/'` on unix, macOS, and browsers, `'\'`
     * on windows).
     * Cannot be set if there’s no `path` yet.
     *
     * @param {string | undefined} extname
     *   Extname.
     * @returns {undefined}
     *   Nothing.
     */
    set extname(extname3) {
      assertPart(extname3, "extname");
      assertPath2(this.dirname, "extname");
      if (extname3) {
        if (extname3.codePointAt(0) !== 46) {
          throw new Error("`extname` must start with `.`");
        }
        if (extname3.includes(".", 1)) {
          throw new Error("`extname` cannot contain multiple dots");
        }
      }
      this.path = minpath.join(this.dirname, this.stem + (extname3 || ""));
    }
    /**
     * Get the full path (example: `'~/index.min.js'`).
     *
     * @returns {string}
     *   Path.
     */
    get path() {
      return this.history[this.history.length - 1];
    }
    /**
     * Set the full path (example: `'~/index.min.js'`).
     *
     * Cannot be nullified.
     * You can set a file URL (a `URL` object with a `file:` protocol) which will
     * be turned into a path with `url.fileURLToPath`.
     *
     * @param {URL | string} path
     *   Path.
     * @returns {undefined}
     *   Nothing.
     */
    set path(path2) {
      if (isUrl(path2)) {
        path2 = urlToPath(path2);
      }
      assertNonEmpty(path2, "path");
      if (this.path !== path2) {
        this.history.push(path2);
      }
    }
    /**
     * Get the stem (basename w/o extname) (example: `'index.min'`).
     *
     * @returns {string | undefined}
     *   Stem.
     */
    get stem() {
      return typeof this.path === "string" ? minpath.basename(this.path, this.extname) : void 0;
    }
    /**
     * Set the stem (basename w/o extname) (example: `'index.min'`).
     *
     * Cannot contain path separators (`'/'` on unix, macOS, and browsers, `'\'`
     * on windows).
     * Cannot be nullified (use `file.path = file.dirname` instead).
     *
     * @param {string} stem
     *   Stem.
     * @returns {undefined}
     *   Nothing.
     */
    set stem(stem) {
      assertNonEmpty(stem, "stem");
      assertPart(stem, "stem");
      this.path = minpath.join(this.dirname || "", stem + (this.extname || ""));
    }
    // Normal prototypal methods.
    /**
     * Create a fatal message for `reason` associated with the file.
     *
     * The `fatal` field of the message is set to `true` (error; file not usable)
     * and the `file` field is set to the current file path.
     * The message is added to the `messages` field on `file`.
     *
     * > 🪦 **Note**: also has obsolete signatures.
     *
     * @overload
     * @param {string} reason
     * @param {MessageOptions | null | undefined} [options]
     * @returns {never}
     *
     * @overload
     * @param {string} reason
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @overload
     * @param {string} reason
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @overload
     * @param {string} reason
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {string | null | undefined} [origin]
     * @returns {never}
     *
     * @param {Error | VFileMessage | string} causeOrReason
     *   Reason for message, should use markdown.
     * @param {Node | NodeLike | MessageOptions | Point | Position | string | null | undefined} [optionsOrParentOrPlace]
     *   Configuration (optional).
     * @param {string | null | undefined} [origin]
     *   Place in code where the message originates (example:
     *   `'my-package:my-rule'` or `'my-rule'`).
     * @returns {never}
     *   Never.
     * @throws {VFileMessage}
     *   Message.
     */
    fail(causeOrReason, optionsOrParentOrPlace, origin) {
      const message = this.message(causeOrReason, optionsOrParentOrPlace, origin);
      message.fatal = true;
      throw message;
    }
    /**
     * Create an info message for `reason` associated with the file.
     *
     * The `fatal` field of the message is set to `undefined` (info; change
     * likely not needed) and the `file` field is set to the current file path.
     * The message is added to the `messages` field on `file`.
     *
     * > 🪦 **Note**: also has obsolete signatures.
     *
     * @overload
     * @param {string} reason
     * @param {MessageOptions | null | undefined} [options]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @param {Error | VFileMessage | string} causeOrReason
     *   Reason for message, should use markdown.
     * @param {Node | NodeLike | MessageOptions | Point | Position | string | null | undefined} [optionsOrParentOrPlace]
     *   Configuration (optional).
     * @param {string | null | undefined} [origin]
     *   Place in code where the message originates (example:
     *   `'my-package:my-rule'` or `'my-rule'`).
     * @returns {VFileMessage}
     *   Message.
     */
    info(causeOrReason, optionsOrParentOrPlace, origin) {
      const message = this.message(causeOrReason, optionsOrParentOrPlace, origin);
      message.fatal = void 0;
      return message;
    }
    /**
     * Create a message for `reason` associated with the file.
     *
     * The `fatal` field of the message is set to `false` (warning; change may be
     * needed) and the `file` field is set to the current file path.
     * The message is added to the `messages` field on `file`.
     *
     * > 🪦 **Note**: also has obsolete signatures.
     *
     * @overload
     * @param {string} reason
     * @param {MessageOptions | null | undefined} [options]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {string} reason
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Node | NodeLike | null | undefined} parent
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {Point | Position | null | undefined} place
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @overload
     * @param {Error | VFileMessage} cause
     * @param {string | null | undefined} [origin]
     * @returns {VFileMessage}
     *
     * @param {Error | VFileMessage | string} causeOrReason
     *   Reason for message, should use markdown.
     * @param {Node | NodeLike | MessageOptions | Point | Position | string | null | undefined} [optionsOrParentOrPlace]
     *   Configuration (optional).
     * @param {string | null | undefined} [origin]
     *   Place in code where the message originates (example:
     *   `'my-package:my-rule'` or `'my-rule'`).
     * @returns {VFileMessage}
     *   Message.
     */
    message(causeOrReason, optionsOrParentOrPlace, origin) {
      const message = new VFileMessage(
        // @ts-expect-error: the overloads are fine.
        causeOrReason,
        optionsOrParentOrPlace,
        origin
      );
      if (this.path) {
        message.name = this.path + ":" + message.name;
        message.file = this.path;
      }
      message.fatal = false;
      this.messages.push(message);
      return message;
    }
    /**
     * Serialize the file.
     *
     * > **Note**: which encodings are supported depends on the engine.
     * > For info on Node.js, see:
     * > <https://nodejs.org/api/util.html#whatwg-supported-encodings>.
     *
     * @param {string | null | undefined} [encoding='utf8']
     *   Character encoding to understand `value` as when it’s a `Uint8Array`
     *   (default: `'utf-8'`).
     * @returns {string}
     *   Serialized file.
     */
    toString(encoding) {
      if (this.value === void 0) {
        return "";
      }
      if (typeof this.value === "string") {
        return this.value;
      }
      const decoder = new TextDecoder(encoding || void 0);
      return decoder.decode(this.value);
    }
  };
  function assertPart(part, name2) {
    if (part && part.includes(minpath.sep)) {
      throw new Error(
        "`" + name2 + "` cannot be a path: did not expect `" + minpath.sep + "`"
      );
    }
  }
  function assertNonEmpty(part, name2) {
    if (!part) {
      throw new Error("`" + name2 + "` cannot be empty");
    }
  }
  function assertPath2(path2, name2) {
    if (!path2) {
      throw new Error("Setting `" + name2 + "` requires `path` to be set too");
    }
  }
  function isUint8Array(value) {
    return Boolean(
      value && typeof value === "object" && "byteLength" in value && "byteOffset" in value
    );
  }

  // ../../../ui/node_modules/.pnpm/unified@11.0.5/node_modules/unified/lib/callable-instance.js
  var CallableInstance = (
    /**
     * @type {new <Parameters extends Array<unknown>, Result>(property: string | symbol) => (...parameters: Parameters) => Result}
     */
    /** @type {unknown} */
    /**
     * @this {Function}
     * @param {string | symbol} property
     * @returns {(...parameters: Array<unknown>) => unknown}
     */
    (function(property2) {
      const self2 = this;
      const constr = self2.constructor;
      const proto = (
        /** @type {Record<string | symbol, Function>} */
        // Prototypes do exist.
        // type-coverage:ignore-next-line
        constr.prototype
      );
      const value = proto[property2];
      const apply = function() {
        return value.apply(apply, arguments);
      };
      Object.setPrototypeOf(apply, proto);
      return apply;
    })
  );

  // ../../../ui/node_modules/.pnpm/unified@11.0.5/node_modules/unified/lib/index.js
  var own = {}.hasOwnProperty;
  var Processor = class _Processor extends CallableInstance {
    /**
     * Create a processor.
     */
    constructor() {
      super("copy");
      this.Compiler = void 0;
      this.Parser = void 0;
      this.attachers = [];
      this.compiler = void 0;
      this.freezeIndex = -1;
      this.frozen = void 0;
      this.namespace = {};
      this.parser = void 0;
      this.transformers = trough();
    }
    /**
     * Copy a processor.
     *
     * @deprecated
     *   This is a private internal method and should not be used.
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *   New *unfrozen* processor ({@linkcode Processor}) that is
     *   configured to work the same as its ancestor.
     *   When the descendant processor is configured in the future it does not
     *   affect the ancestral processor.
     */
    copy() {
      const destination = (
        /** @type {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>} */
        new _Processor()
      );
      let index2 = -1;
      while (++index2 < this.attachers.length) {
        const attacher = this.attachers[index2];
        destination.use(...attacher);
      }
      destination.data((0, import_extend.default)(true, {}, this.namespace));
      return destination;
    }
    /**
     * Configure the processor with info available to all plugins.
     * Information is stored in an object.
     *
     * Typically, options can be given to a specific plugin, but sometimes it
     * makes sense to have information shared with several plugins.
     * For example, a list of HTML elements that are self-closing, which is
     * needed during all phases.
     *
     * > **Note**: setting information cannot occur on *frozen* processors.
     * > Call the processor first to create a new unfrozen processor.
     *
     * > **Note**: to register custom data in TypeScript, augment the
     * > {@linkcode Data} interface.
     *
     * @example
     *   This example show how to get and set info:
     *
     *   ```js
     *   import {unified} from 'unified'
     *
     *   const processor = unified().data('alpha', 'bravo')
     *
     *   processor.data('alpha') // => 'bravo'
     *
     *   processor.data() // => {alpha: 'bravo'}
     *
     *   processor.data({charlie: 'delta'})
     *
     *   processor.data() // => {charlie: 'delta'}
     *   ```
     *
     * @template {keyof Data} Key
     *
     * @overload
     * @returns {Data}
     *
     * @overload
     * @param {Data} dataset
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *
     * @overload
     * @param {Key} key
     * @returns {Data[Key]}
     *
     * @overload
     * @param {Key} key
     * @param {Data[Key]} value
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *
     * @param {Data | Key} [key]
     *   Key to get or set, or entire dataset to set, or nothing to get the
     *   entire dataset (optional).
     * @param {Data[Key]} [value]
     *   Value to set (optional).
     * @returns {unknown}
     *   The current processor when setting, the value at `key` when getting, or
     *   the entire dataset when getting without key.
     */
    data(key, value) {
      if (typeof key === "string") {
        if (arguments.length === 2) {
          assertUnfrozen("data", this.frozen);
          this.namespace[key] = value;
          return this;
        }
        return own.call(this.namespace, key) && this.namespace[key] || void 0;
      }
      if (key) {
        assertUnfrozen("data", this.frozen);
        this.namespace = key;
        return this;
      }
      return this.namespace;
    }
    /**
     * Freeze a processor.
     *
     * Frozen processors are meant to be extended and not to be configured
     * directly.
     *
     * When a processor is frozen it cannot be unfrozen.
     * New processors working the same way can be created by calling the
     * processor.
     *
     * It’s possible to freeze processors explicitly by calling `.freeze()`.
     * Processors freeze automatically when `.parse()`, `.run()`, `.runSync()`,
     * `.stringify()`, `.process()`, or `.processSync()` are called.
     *
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *   The current processor.
     */
    freeze() {
      if (this.frozen) {
        return this;
      }
      const self2 = (
        /** @type {Processor} */
        /** @type {unknown} */
        this
      );
      while (++this.freezeIndex < this.attachers.length) {
        const [attacher, ...options] = this.attachers[this.freezeIndex];
        if (options[0] === false) {
          continue;
        }
        if (options[0] === true) {
          options[0] = void 0;
        }
        const transformer = attacher.call(self2, ...options);
        if (typeof transformer === "function") {
          this.transformers.use(transformer);
        }
      }
      this.frozen = true;
      this.freezeIndex = Number.POSITIVE_INFINITY;
      return this;
    }
    /**
     * Parse text to a syntax tree.
     *
     * > **Note**: `parse` freezes the processor if not already *frozen*.
     *
     * > **Note**: `parse` performs the parse phase, not the run phase or other
     * > phases.
     *
     * @param {Compatible | undefined} [file]
     *   file to parse (optional); typically `string` or `VFile`; any value
     *   accepted as `x` in `new VFile(x)`.
     * @returns {ParseTree extends undefined ? Node : ParseTree}
     *   Syntax tree representing `file`.
     */
    parse(file) {
      this.freeze();
      const realFile = vfile(file);
      const parser = this.parser || this.Parser;
      assertParser("parse", parser);
      return parser(String(realFile), realFile);
    }
    /**
     * Process the given file as configured on the processor.
     *
     * > **Note**: `process` freezes the processor if not already *frozen*.
     *
     * > **Note**: `process` performs the parse, run, and stringify phases.
     *
     * @overload
     * @param {Compatible | undefined} file
     * @param {ProcessCallback<VFileWithOutput<CompileResult>>} done
     * @returns {undefined}
     *
     * @overload
     * @param {Compatible | undefined} [file]
     * @returns {Promise<VFileWithOutput<CompileResult>>}
     *
     * @param {Compatible | undefined} [file]
     *   File (optional); typically `string` or `VFile`]; any value accepted as
     *   `x` in `new VFile(x)`.
     * @param {ProcessCallback<VFileWithOutput<CompileResult>> | undefined} [done]
     *   Callback (optional).
     * @returns {Promise<VFile> | undefined}
     *   Nothing if `done` is given.
     *   Otherwise a promise, rejected with a fatal error or resolved with the
     *   processed file.
     *
     *   The parsed, transformed, and compiled value is available at
     *   `file.value` (see note).
     *
     *   > **Note**: unified typically compiles by serializing: most
     *   > compilers return `string` (or `Uint8Array`).
     *   > Some compilers, such as the one configured with
     *   > [`rehype-react`][rehype-react], return other values (in this case, a
     *   > React tree).
     *   > If you’re using a compiler that doesn’t serialize, expect different
     *   > result values.
     *   >
     *   > To register custom results in TypeScript, add them to
     *   > {@linkcode CompileResultMap}.
     *
     *   [rehype-react]: https://github.com/rehypejs/rehype-react
     */
    process(file, done) {
      const self2 = this;
      this.freeze();
      assertParser("process", this.parser || this.Parser);
      assertCompiler("process", this.compiler || this.Compiler);
      return done ? executor(void 0, done) : new Promise(executor);
      function executor(resolve, reject) {
        const realFile = vfile(file);
        const parseTree = (
          /** @type {HeadTree extends undefined ? Node : HeadTree} */
          /** @type {unknown} */
          self2.parse(realFile)
        );
        self2.run(parseTree, realFile, function(error, tree, file2) {
          if (error || !tree || !file2) {
            return realDone(error);
          }
          const compileTree = (
            /** @type {CompileTree extends undefined ? Node : CompileTree} */
            /** @type {unknown} */
            tree
          );
          const compileResult = self2.stringify(compileTree, file2);
          if (looksLikeAValue(compileResult)) {
            file2.value = compileResult;
          } else {
            file2.result = compileResult;
          }
          realDone(
            error,
            /** @type {VFileWithOutput<CompileResult>} */
            file2
          );
        });
        function realDone(error, file2) {
          if (error || !file2) {
            reject(error);
          } else if (resolve) {
            resolve(file2);
          } else {
            ok(done, "`done` is defined if `resolve` is not");
            done(void 0, file2);
          }
        }
      }
    }
    /**
     * Process the given file as configured on the processor.
     *
     * An error is thrown if asynchronous transforms are configured.
     *
     * > **Note**: `processSync` freezes the processor if not already *frozen*.
     *
     * > **Note**: `processSync` performs the parse, run, and stringify phases.
     *
     * @param {Compatible | undefined} [file]
     *   File (optional); typically `string` or `VFile`; any value accepted as
     *   `x` in `new VFile(x)`.
     * @returns {VFileWithOutput<CompileResult>}
     *   The processed file.
     *
     *   The parsed, transformed, and compiled value is available at
     *   `file.value` (see note).
     *
     *   > **Note**: unified typically compiles by serializing: most
     *   > compilers return `string` (or `Uint8Array`).
     *   > Some compilers, such as the one configured with
     *   > [`rehype-react`][rehype-react], return other values (in this case, a
     *   > React tree).
     *   > If you’re using a compiler that doesn’t serialize, expect different
     *   > result values.
     *   >
     *   > To register custom results in TypeScript, add them to
     *   > {@linkcode CompileResultMap}.
     *
     *   [rehype-react]: https://github.com/rehypejs/rehype-react
     */
    processSync(file) {
      let complete = false;
      let result;
      this.freeze();
      assertParser("processSync", this.parser || this.Parser);
      assertCompiler("processSync", this.compiler || this.Compiler);
      this.process(file, realDone);
      assertDone("processSync", "process", complete);
      ok(result, "we either bailed on an error or have a tree");
      return result;
      function realDone(error, file2) {
        complete = true;
        bail(error);
        result = file2;
      }
    }
    /**
     * Run *transformers* on a syntax tree.
     *
     * > **Note**: `run` freezes the processor if not already *frozen*.
     *
     * > **Note**: `run` performs the run phase, not other phases.
     *
     * @overload
     * @param {HeadTree extends undefined ? Node : HeadTree} tree
     * @param {RunCallback<TailTree extends undefined ? Node : TailTree>} done
     * @returns {undefined}
     *
     * @overload
     * @param {HeadTree extends undefined ? Node : HeadTree} tree
     * @param {Compatible | undefined} file
     * @param {RunCallback<TailTree extends undefined ? Node : TailTree>} done
     * @returns {undefined}
     *
     * @overload
     * @param {HeadTree extends undefined ? Node : HeadTree} tree
     * @param {Compatible | undefined} [file]
     * @returns {Promise<TailTree extends undefined ? Node : TailTree>}
     *
     * @param {HeadTree extends undefined ? Node : HeadTree} tree
     *   Tree to transform and inspect.
     * @param {(
     *   RunCallback<TailTree extends undefined ? Node : TailTree> |
     *   Compatible
     * )} [file]
     *   File associated with `node` (optional); any value accepted as `x` in
     *   `new VFile(x)`.
     * @param {RunCallback<TailTree extends undefined ? Node : TailTree>} [done]
     *   Callback (optional).
     * @returns {Promise<TailTree extends undefined ? Node : TailTree> | undefined}
     *   Nothing if `done` is given.
     *   Otherwise, a promise rejected with a fatal error or resolved with the
     *   transformed tree.
     */
    run(tree, file, done) {
      assertNode(tree);
      this.freeze();
      const transformers = this.transformers;
      if (!done && typeof file === "function") {
        done = file;
        file = void 0;
      }
      return done ? executor(void 0, done) : new Promise(executor);
      function executor(resolve, reject) {
        ok(
          typeof file !== "function",
          "`file` can’t be a `done` anymore, we checked"
        );
        const realFile = vfile(file);
        transformers.run(tree, realFile, realDone);
        function realDone(error, outputTree, file2) {
          const resultingTree = (
            /** @type {TailTree extends undefined ? Node : TailTree} */
            outputTree || tree
          );
          if (error) {
            reject(error);
          } else if (resolve) {
            resolve(resultingTree);
          } else {
            ok(done, "`done` is defined if `resolve` is not");
            done(void 0, resultingTree, file2);
          }
        }
      }
    }
    /**
     * Run *transformers* on a syntax tree.
     *
     * An error is thrown if asynchronous transforms are configured.
     *
     * > **Note**: `runSync` freezes the processor if not already *frozen*.
     *
     * > **Note**: `runSync` performs the run phase, not other phases.
     *
     * @param {HeadTree extends undefined ? Node : HeadTree} tree
     *   Tree to transform and inspect.
     * @param {Compatible | undefined} [file]
     *   File associated with `node` (optional); any value accepted as `x` in
     *   `new VFile(x)`.
     * @returns {TailTree extends undefined ? Node : TailTree}
     *   Transformed tree.
     */
    runSync(tree, file) {
      let complete = false;
      let result;
      this.run(tree, file, realDone);
      assertDone("runSync", "run", complete);
      ok(result, "we either bailed on an error or have a tree");
      return result;
      function realDone(error, tree2) {
        bail(error);
        result = tree2;
        complete = true;
      }
    }
    /**
     * Compile a syntax tree.
     *
     * > **Note**: `stringify` freezes the processor if not already *frozen*.
     *
     * > **Note**: `stringify` performs the stringify phase, not the run phase
     * > or other phases.
     *
     * @param {CompileTree extends undefined ? Node : CompileTree} tree
     *   Tree to compile.
     * @param {Compatible | undefined} [file]
     *   File associated with `node` (optional); any value accepted as `x` in
     *   `new VFile(x)`.
     * @returns {CompileResult extends undefined ? Value : CompileResult}
     *   Textual representation of the tree (see note).
     *
     *   > **Note**: unified typically compiles by serializing: most compilers
     *   > return `string` (or `Uint8Array`).
     *   > Some compilers, such as the one configured with
     *   > [`rehype-react`][rehype-react], return other values (in this case, a
     *   > React tree).
     *   > If you’re using a compiler that doesn’t serialize, expect different
     *   > result values.
     *   >
     *   > To register custom results in TypeScript, add them to
     *   > {@linkcode CompileResultMap}.
     *
     *   [rehype-react]: https://github.com/rehypejs/rehype-react
     */
    stringify(tree, file) {
      this.freeze();
      const realFile = vfile(file);
      const compiler2 = this.compiler || this.Compiler;
      assertCompiler("stringify", compiler2);
      assertNode(tree);
      return compiler2(tree, realFile);
    }
    /**
     * Configure the processor to use a plugin, a list of usable values, or a
     * preset.
     *
     * If the processor is already using a plugin, the previous plugin
     * configuration is changed based on the options that are passed in.
     * In other words, the plugin is not added a second time.
     *
     * > **Note**: `use` cannot be called on *frozen* processors.
     * > Call the processor first to create a new unfrozen processor.
     *
     * @example
     *   There are many ways to pass plugins to `.use()`.
     *   This example gives an overview:
     *
     *   ```js
     *   import {unified} from 'unified'
     *
     *   unified()
     *     // Plugin with options:
     *     .use(pluginA, {x: true, y: true})
     *     // Passing the same plugin again merges configuration (to `{x: true, y: false, z: true}`):
     *     .use(pluginA, {y: false, z: true})
     *     // Plugins:
     *     .use([pluginB, pluginC])
     *     // Two plugins, the second with options:
     *     .use([pluginD, [pluginE, {}]])
     *     // Preset with plugins and settings:
     *     .use({plugins: [pluginF, [pluginG, {}]], settings: {position: false}})
     *     // Settings only:
     *     .use({settings: {position: false}})
     *   ```
     *
     * @template {Array<unknown>} [Parameters=[]]
     * @template {Node | string | undefined} [Input=undefined]
     * @template [Output=Input]
     *
     * @overload
     * @param {Preset | null | undefined} [preset]
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *
     * @overload
     * @param {PluggableList} list
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *
     * @overload
     * @param {Plugin<Parameters, Input, Output>} plugin
     * @param {...(Parameters | [boolean])} parameters
     * @returns {UsePlugin<ParseTree, HeadTree, TailTree, CompileTree, CompileResult, Input, Output>}
     *
     * @param {PluggableList | Plugin | Preset | null | undefined} value
     *   Usable value.
     * @param {...unknown} parameters
     *   Parameters, when a plugin is given as a usable value.
     * @returns {Processor<ParseTree, HeadTree, TailTree, CompileTree, CompileResult>}
     *   Current processor.
     */
    use(value, ...parameters) {
      const attachers = this.attachers;
      const namespace = this.namespace;
      assertUnfrozen("use", this.frozen);
      if (value === null || value === void 0) {
      } else if (typeof value === "function") {
        addPlugin(value, parameters);
      } else if (typeof value === "object") {
        if (Array.isArray(value)) {
          addList(value);
        } else {
          addPreset(value);
        }
      } else {
        throw new TypeError("Expected usable value, not `" + value + "`");
      }
      return this;
      function add(value2) {
        if (typeof value2 === "function") {
          addPlugin(value2, []);
        } else if (typeof value2 === "object") {
          if (Array.isArray(value2)) {
            const [plugin, ...parameters2] = (
              /** @type {PluginTuple<Array<unknown>>} */
              value2
            );
            addPlugin(plugin, parameters2);
          } else {
            addPreset(value2);
          }
        } else {
          throw new TypeError("Expected usable value, not `" + value2 + "`");
        }
      }
      function addPreset(result) {
        if (!("plugins" in result) && !("settings" in result)) {
          throw new Error(
            "Expected usable value but received an empty preset, which is probably a mistake: presets typically come with `plugins` and sometimes with `settings`, but this has neither"
          );
        }
        addList(result.plugins);
        if (result.settings) {
          namespace.settings = (0, import_extend.default)(true, namespace.settings, result.settings);
        }
      }
      function addList(plugins) {
        let index2 = -1;
        if (plugins === null || plugins === void 0) {
        } else if (Array.isArray(plugins)) {
          while (++index2 < plugins.length) {
            const thing = plugins[index2];
            add(thing);
          }
        } else {
          throw new TypeError("Expected a list of plugins, not `" + plugins + "`");
        }
      }
      function addPlugin(plugin, parameters2) {
        let index2 = -1;
        let entryIndex = -1;
        while (++index2 < attachers.length) {
          if (attachers[index2][0] === plugin) {
            entryIndex = index2;
            break;
          }
        }
        if (entryIndex === -1) {
          attachers.push([plugin, ...parameters2]);
        } else if (parameters2.length > 0) {
          let [primary, ...rest] = parameters2;
          const currentPrimary = attachers[entryIndex][1];
          if (isPlainObject(currentPrimary) && isPlainObject(primary)) {
            primary = (0, import_extend.default)(true, currentPrimary, primary);
          }
          attachers[entryIndex] = [plugin, primary, ...rest];
        }
      }
    }
  };
  var unified = new Processor().freeze();
  function assertParser(name2, value) {
    if (typeof value !== "function") {
      throw new TypeError("Cannot `" + name2 + "` without `parser`");
    }
  }
  function assertCompiler(name2, value) {
    if (typeof value !== "function") {
      throw new TypeError("Cannot `" + name2 + "` without `compiler`");
    }
  }
  function assertUnfrozen(name2, frozen) {
    if (frozen) {
      throw new Error(
        "Cannot call `" + name2 + "` on a frozen processor.\nCreate a new processor first, by calling it: use `processor()` instead of `processor`."
      );
    }
  }
  function assertNode(node2) {
    if (!isPlainObject(node2) || typeof node2.type !== "string") {
      throw new TypeError("Expected node, got `" + node2 + "`");
    }
  }
  function assertDone(name2, asyncName, complete) {
    if (!complete) {
      throw new Error(
        "`" + name2 + "` finished async. Use `" + asyncName + "` instead"
      );
    }
  }
  function vfile(value) {
    return looksLikeAVFile(value) ? value : new VFile(value);
  }
  function looksLikeAVFile(value) {
    return Boolean(
      value && typeof value === "object" && "message" in value && "messages" in value
    );
  }
  function looksLikeAValue(value) {
    return typeof value === "string" || isUint8Array2(value);
  }
  function isUint8Array2(value) {
    return Boolean(
      value && typeof value === "object" && "byteLength" in value && "byteOffset" in value
    );
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-mdx-expression@2.0.1/node_modules/mdast-util-mdx-expression/lib/index.js
  function mdxExpressionFromMarkdown() {
    return {
      enter: {
        mdxFlowExpression: enterMdxFlowExpression,
        mdxTextExpression: enterMdxTextExpression
      },
      exit: {
        mdxFlowExpression: exitMdxExpression,
        mdxFlowExpressionChunk: exitMdxExpressionData,
        mdxTextExpression: exitMdxExpression,
        mdxTextExpressionChunk: exitMdxExpressionData
      }
    };
  }
  function mdxExpressionToMarkdown() {
    return {
      handlers: {
        mdxFlowExpression: handleMdxExpression,
        mdxTextExpression: handleMdxExpression
      },
      unsafe: [
        { character: "{", inConstruct: ["phrasing"] },
        { atBreak: true, character: "{" }
      ]
    };
  }
  function enterMdxFlowExpression(token) {
    this.enter({ type: "mdxFlowExpression", value: "" }, token);
    this.buffer();
  }
  function enterMdxTextExpression(token) {
    this.enter({ type: "mdxTextExpression", value: "" }, token);
    this.buffer();
  }
  function exitMdxExpression(token) {
    const value = this.resume();
    const estree = token.estree;
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "mdxFlowExpression" || node2.type === "mdxTextExpression");
    this.exit(token);
    node2.value = value;
    if (estree) {
      node2.data = { estree };
    }
  }
  function exitMdxExpressionData(token) {
    this.config.enter.data.call(this, token);
    this.config.exit.data.call(this, token);
  }
  function handleMdxExpression(node2, parent, state) {
    const value = node2.value || "";
    const result = state.indentLines(value, function(line, index2, blank) {
      return (index2 === 0 || blank ? "" : "  ") + line;
    });
    return "{" + result + "}";
  }

  // ../../../ui/node_modules/.pnpm/ccount@2.0.1/node_modules/ccount/index.js
  function ccount(value, character) {
    const source = String(value);
    if (typeof character !== "string") {
      throw new TypeError("Expected character");
    }
    let count3 = 0;
    let index2 = source.indexOf(character);
    while (index2 !== -1) {
      count3++;
      index2 = source.indexOf(character, index2 + character.length);
    }
    return count3;
  }

  // ../../../ui/node_modules/.pnpm/character-entities-legacy@3.0.0/node_modules/character-entities-legacy/index.js
  var characterEntitiesLegacy = [
    "AElig",
    "AMP",
    "Aacute",
    "Acirc",
    "Agrave",
    "Aring",
    "Atilde",
    "Auml",
    "COPY",
    "Ccedil",
    "ETH",
    "Eacute",
    "Ecirc",
    "Egrave",
    "Euml",
    "GT",
    "Iacute",
    "Icirc",
    "Igrave",
    "Iuml",
    "LT",
    "Ntilde",
    "Oacute",
    "Ocirc",
    "Ograve",
    "Oslash",
    "Otilde",
    "Ouml",
    "QUOT",
    "REG",
    "THORN",
    "Uacute",
    "Ucirc",
    "Ugrave",
    "Uuml",
    "Yacute",
    "aacute",
    "acirc",
    "acute",
    "aelig",
    "agrave",
    "amp",
    "aring",
    "atilde",
    "auml",
    "brvbar",
    "ccedil",
    "cedil",
    "cent",
    "copy",
    "curren",
    "deg",
    "divide",
    "eacute",
    "ecirc",
    "egrave",
    "eth",
    "euml",
    "frac12",
    "frac14",
    "frac34",
    "gt",
    "iacute",
    "icirc",
    "iexcl",
    "igrave",
    "iquest",
    "iuml",
    "laquo",
    "lt",
    "macr",
    "micro",
    "middot",
    "nbsp",
    "not",
    "ntilde",
    "oacute",
    "ocirc",
    "ograve",
    "ordf",
    "ordm",
    "oslash",
    "otilde",
    "ouml",
    "para",
    "plusmn",
    "pound",
    "quot",
    "raquo",
    "reg",
    "sect",
    "shy",
    "sup1",
    "sup2",
    "sup3",
    "szlig",
    "thorn",
    "times",
    "uacute",
    "ucirc",
    "ugrave",
    "uml",
    "uuml",
    "yacute",
    "yen",
    "yuml"
  ];

  // ../../../ui/node_modules/.pnpm/character-reference-invalid@2.0.1/node_modules/character-reference-invalid/index.js
  var characterReferenceInvalid = {
    0: "�",
    128: "€",
    130: "‚",
    131: "ƒ",
    132: "„",
    133: "…",
    134: "†",
    135: "‡",
    136: "ˆ",
    137: "‰",
    138: "Š",
    139: "‹",
    140: "Œ",
    142: "Ž",
    145: "‘",
    146: "’",
    147: "“",
    148: "”",
    149: "•",
    150: "–",
    151: "—",
    152: "˜",
    153: "™",
    154: "š",
    155: "›",
    156: "œ",
    158: "ž",
    159: "Ÿ"
  };

  // ../../../ui/node_modules/.pnpm/is-decimal@2.0.1/node_modules/is-decimal/index.js
  function isDecimal(character) {
    const code3 = typeof character === "string" ? character.charCodeAt(0) : character;
    return code3 >= 48 && code3 <= 57;
  }

  // ../../../ui/node_modules/.pnpm/is-hexadecimal@2.0.1/node_modules/is-hexadecimal/index.js
  function isHexadecimal(character) {
    const code3 = typeof character === "string" ? character.charCodeAt(0) : character;
    return code3 >= 97 && code3 <= 102 || code3 >= 65 && code3 <= 70 || code3 >= 48 && code3 <= 57;
  }

  // ../../../ui/node_modules/.pnpm/is-alphabetical@2.0.1/node_modules/is-alphabetical/index.js
  function isAlphabetical(character) {
    const code3 = typeof character === "string" ? character.charCodeAt(0) : character;
    return code3 >= 97 && code3 <= 122 || code3 >= 65 && code3 <= 90;
  }

  // ../../../ui/node_modules/.pnpm/is-alphanumerical@2.0.1/node_modules/is-alphanumerical/index.js
  function isAlphanumerical(character) {
    return isAlphabetical(character) || isDecimal(character);
  }

  // ../../../ui/node_modules/.pnpm/character-entities@2.0.2/node_modules/character-entities/index.js
  var characterEntities = {
    AElig: "Æ",
    AMP: "&",
    Aacute: "Á",
    Abreve: "Ă",
    Acirc: "Â",
    Acy: "А",
    Afr: "𝔄",
    Agrave: "À",
    Alpha: "Α",
    Amacr: "Ā",
    And: "⩓",
    Aogon: "Ą",
    Aopf: "𝔸",
    ApplyFunction: "⁡",
    Aring: "Å",
    Ascr: "𝒜",
    Assign: "≔",
    Atilde: "Ã",
    Auml: "Ä",
    Backslash: "∖",
    Barv: "⫧",
    Barwed: "⌆",
    Bcy: "Б",
    Because: "∵",
    Bernoullis: "ℬ",
    Beta: "Β",
    Bfr: "𝔅",
    Bopf: "𝔹",
    Breve: "˘",
    Bscr: "ℬ",
    Bumpeq: "≎",
    CHcy: "Ч",
    COPY: "©",
    Cacute: "Ć",
    Cap: "⋒",
    CapitalDifferentialD: "ⅅ",
    Cayleys: "ℭ",
    Ccaron: "Č",
    Ccedil: "Ç",
    Ccirc: "Ĉ",
    Cconint: "∰",
    Cdot: "Ċ",
    Cedilla: "¸",
    CenterDot: "·",
    Cfr: "ℭ",
    Chi: "Χ",
    CircleDot: "⊙",
    CircleMinus: "⊖",
    CirclePlus: "⊕",
    CircleTimes: "⊗",
    ClockwiseContourIntegral: "∲",
    CloseCurlyDoubleQuote: "”",
    CloseCurlyQuote: "’",
    Colon: "∷",
    Colone: "⩴",
    Congruent: "≡",
    Conint: "∯",
    ContourIntegral: "∮",
    Copf: "ℂ",
    Coproduct: "∐",
    CounterClockwiseContourIntegral: "∳",
    Cross: "⨯",
    Cscr: "𝒞",
    Cup: "⋓",
    CupCap: "≍",
    DD: "ⅅ",
    DDotrahd: "⤑",
    DJcy: "Ђ",
    DScy: "Ѕ",
    DZcy: "Џ",
    Dagger: "‡",
    Darr: "↡",
    Dashv: "⫤",
    Dcaron: "Ď",
    Dcy: "Д",
    Del: "∇",
    Delta: "Δ",
    Dfr: "𝔇",
    DiacriticalAcute: "´",
    DiacriticalDot: "˙",
    DiacriticalDoubleAcute: "˝",
    DiacriticalGrave: "`",
    DiacriticalTilde: "˜",
    Diamond: "⋄",
    DifferentialD: "ⅆ",
    Dopf: "𝔻",
    Dot: "¨",
    DotDot: "⃜",
    DotEqual: "≐",
    DoubleContourIntegral: "∯",
    DoubleDot: "¨",
    DoubleDownArrow: "⇓",
    DoubleLeftArrow: "⇐",
    DoubleLeftRightArrow: "⇔",
    DoubleLeftTee: "⫤",
    DoubleLongLeftArrow: "⟸",
    DoubleLongLeftRightArrow: "⟺",
    DoubleLongRightArrow: "⟹",
    DoubleRightArrow: "⇒",
    DoubleRightTee: "⊨",
    DoubleUpArrow: "⇑",
    DoubleUpDownArrow: "⇕",
    DoubleVerticalBar: "∥",
    DownArrow: "↓",
    DownArrowBar: "⤓",
    DownArrowUpArrow: "⇵",
    DownBreve: "̑",
    DownLeftRightVector: "⥐",
    DownLeftTeeVector: "⥞",
    DownLeftVector: "↽",
    DownLeftVectorBar: "⥖",
    DownRightTeeVector: "⥟",
    DownRightVector: "⇁",
    DownRightVectorBar: "⥗",
    DownTee: "⊤",
    DownTeeArrow: "↧",
    Downarrow: "⇓",
    Dscr: "𝒟",
    Dstrok: "Đ",
    ENG: "Ŋ",
    ETH: "Ð",
    Eacute: "É",
    Ecaron: "Ě",
    Ecirc: "Ê",
    Ecy: "Э",
    Edot: "Ė",
    Efr: "𝔈",
    Egrave: "È",
    Element: "∈",
    Emacr: "Ē",
    EmptySmallSquare: "◻",
    EmptyVerySmallSquare: "▫",
    Eogon: "Ę",
    Eopf: "𝔼",
    Epsilon: "Ε",
    Equal: "⩵",
    EqualTilde: "≂",
    Equilibrium: "⇌",
    Escr: "ℰ",
    Esim: "⩳",
    Eta: "Η",
    Euml: "Ë",
    Exists: "∃",
    ExponentialE: "ⅇ",
    Fcy: "Ф",
    Ffr: "𝔉",
    FilledSmallSquare: "◼",
    FilledVerySmallSquare: "▪",
    Fopf: "𝔽",
    ForAll: "∀",
    Fouriertrf: "ℱ",
    Fscr: "ℱ",
    GJcy: "Ѓ",
    GT: ">",
    Gamma: "Γ",
    Gammad: "Ϝ",
    Gbreve: "Ğ",
    Gcedil: "Ģ",
    Gcirc: "Ĝ",
    Gcy: "Г",
    Gdot: "Ġ",
    Gfr: "𝔊",
    Gg: "⋙",
    Gopf: "𝔾",
    GreaterEqual: "≥",
    GreaterEqualLess: "⋛",
    GreaterFullEqual: "≧",
    GreaterGreater: "⪢",
    GreaterLess: "≷",
    GreaterSlantEqual: "⩾",
    GreaterTilde: "≳",
    Gscr: "𝒢",
    Gt: "≫",
    HARDcy: "Ъ",
    Hacek: "ˇ",
    Hat: "^",
    Hcirc: "Ĥ",
    Hfr: "ℌ",
    HilbertSpace: "ℋ",
    Hopf: "ℍ",
    HorizontalLine: "─",
    Hscr: "ℋ",
    Hstrok: "Ħ",
    HumpDownHump: "≎",
    HumpEqual: "≏",
    IEcy: "Е",
    IJlig: "Ĳ",
    IOcy: "Ё",
    Iacute: "Í",
    Icirc: "Î",
    Icy: "И",
    Idot: "İ",
    Ifr: "ℑ",
    Igrave: "Ì",
    Im: "ℑ",
    Imacr: "Ī",
    ImaginaryI: "ⅈ",
    Implies: "⇒",
    Int: "∬",
    Integral: "∫",
    Intersection: "⋂",
    InvisibleComma: "⁣",
    InvisibleTimes: "⁢",
    Iogon: "Į",
    Iopf: "𝕀",
    Iota: "Ι",
    Iscr: "ℐ",
    Itilde: "Ĩ",
    Iukcy: "І",
    Iuml: "Ï",
    Jcirc: "Ĵ",
    Jcy: "Й",
    Jfr: "𝔍",
    Jopf: "𝕁",
    Jscr: "𝒥",
    Jsercy: "Ј",
    Jukcy: "Є",
    KHcy: "Х",
    KJcy: "Ќ",
    Kappa: "Κ",
    Kcedil: "Ķ",
    Kcy: "К",
    Kfr: "𝔎",
    Kopf: "𝕂",
    Kscr: "𝒦",
    LJcy: "Љ",
    LT: "<",
    Lacute: "Ĺ",
    Lambda: "Λ",
    Lang: "⟪",
    Laplacetrf: "ℒ",
    Larr: "↞",
    Lcaron: "Ľ",
    Lcedil: "Ļ",
    Lcy: "Л",
    LeftAngleBracket: "⟨",
    LeftArrow: "←",
    LeftArrowBar: "⇤",
    LeftArrowRightArrow: "⇆",
    LeftCeiling: "⌈",
    LeftDoubleBracket: "⟦",
    LeftDownTeeVector: "⥡",
    LeftDownVector: "⇃",
    LeftDownVectorBar: "⥙",
    LeftFloor: "⌊",
    LeftRightArrow: "↔",
    LeftRightVector: "⥎",
    LeftTee: "⊣",
    LeftTeeArrow: "↤",
    LeftTeeVector: "⥚",
    LeftTriangle: "⊲",
    LeftTriangleBar: "⧏",
    LeftTriangleEqual: "⊴",
    LeftUpDownVector: "⥑",
    LeftUpTeeVector: "⥠",
    LeftUpVector: "↿",
    LeftUpVectorBar: "⥘",
    LeftVector: "↼",
    LeftVectorBar: "⥒",
    Leftarrow: "⇐",
    Leftrightarrow: "⇔",
    LessEqualGreater: "⋚",
    LessFullEqual: "≦",
    LessGreater: "≶",
    LessLess: "⪡",
    LessSlantEqual: "⩽",
    LessTilde: "≲",
    Lfr: "𝔏",
    Ll: "⋘",
    Lleftarrow: "⇚",
    Lmidot: "Ŀ",
    LongLeftArrow: "⟵",
    LongLeftRightArrow: "⟷",
    LongRightArrow: "⟶",
    Longleftarrow: "⟸",
    Longleftrightarrow: "⟺",
    Longrightarrow: "⟹",
    Lopf: "𝕃",
    LowerLeftArrow: "↙",
    LowerRightArrow: "↘",
    Lscr: "ℒ",
    Lsh: "↰",
    Lstrok: "Ł",
    Lt: "≪",
    Map: "⤅",
    Mcy: "М",
    MediumSpace: " ",
    Mellintrf: "ℳ",
    Mfr: "𝔐",
    MinusPlus: "∓",
    Mopf: "𝕄",
    Mscr: "ℳ",
    Mu: "Μ",
    NJcy: "Њ",
    Nacute: "Ń",
    Ncaron: "Ň",
    Ncedil: "Ņ",
    Ncy: "Н",
    NegativeMediumSpace: "​",
    NegativeThickSpace: "​",
    NegativeThinSpace: "​",
    NegativeVeryThinSpace: "​",
    NestedGreaterGreater: "≫",
    NestedLessLess: "≪",
    NewLine: "\n",
    Nfr: "𝔑",
    NoBreak: "⁠",
    NonBreakingSpace: " ",
    Nopf: "ℕ",
    Not: "⫬",
    NotCongruent: "≢",
    NotCupCap: "≭",
    NotDoubleVerticalBar: "∦",
    NotElement: "∉",
    NotEqual: "≠",
    NotEqualTilde: "≂̸",
    NotExists: "∄",
    NotGreater: "≯",
    NotGreaterEqual: "≱",
    NotGreaterFullEqual: "≧̸",
    NotGreaterGreater: "≫̸",
    NotGreaterLess: "≹",
    NotGreaterSlantEqual: "⩾̸",
    NotGreaterTilde: "≵",
    NotHumpDownHump: "≎̸",
    NotHumpEqual: "≏̸",
    NotLeftTriangle: "⋪",
    NotLeftTriangleBar: "⧏̸",
    NotLeftTriangleEqual: "⋬",
    NotLess: "≮",
    NotLessEqual: "≰",
    NotLessGreater: "≸",
    NotLessLess: "≪̸",
    NotLessSlantEqual: "⩽̸",
    NotLessTilde: "≴",
    NotNestedGreaterGreater: "⪢̸",
    NotNestedLessLess: "⪡̸",
    NotPrecedes: "⊀",
    NotPrecedesEqual: "⪯̸",
    NotPrecedesSlantEqual: "⋠",
    NotReverseElement: "∌",
    NotRightTriangle: "⋫",
    NotRightTriangleBar: "⧐̸",
    NotRightTriangleEqual: "⋭",
    NotSquareSubset: "⊏̸",
    NotSquareSubsetEqual: "⋢",
    NotSquareSuperset: "⊐̸",
    NotSquareSupersetEqual: "⋣",
    NotSubset: "⊂⃒",
    NotSubsetEqual: "⊈",
    NotSucceeds: "⊁",
    NotSucceedsEqual: "⪰̸",
    NotSucceedsSlantEqual: "⋡",
    NotSucceedsTilde: "≿̸",
    NotSuperset: "⊃⃒",
    NotSupersetEqual: "⊉",
    NotTilde: "≁",
    NotTildeEqual: "≄",
    NotTildeFullEqual: "≇",
    NotTildeTilde: "≉",
    NotVerticalBar: "∤",
    Nscr: "𝒩",
    Ntilde: "Ñ",
    Nu: "Ν",
    OElig: "Œ",
    Oacute: "Ó",
    Ocirc: "Ô",
    Ocy: "О",
    Odblac: "Ő",
    Ofr: "𝔒",
    Ograve: "Ò",
    Omacr: "Ō",
    Omega: "Ω",
    Omicron: "Ο",
    Oopf: "𝕆",
    OpenCurlyDoubleQuote: "“",
    OpenCurlyQuote: "‘",
    Or: "⩔",
    Oscr: "𝒪",
    Oslash: "Ø",
    Otilde: "Õ",
    Otimes: "⨷",
    Ouml: "Ö",
    OverBar: "‾",
    OverBrace: "⏞",
    OverBracket: "⎴",
    OverParenthesis: "⏜",
    PartialD: "∂",
    Pcy: "П",
    Pfr: "𝔓",
    Phi: "Φ",
    Pi: "Π",
    PlusMinus: "±",
    Poincareplane: "ℌ",
    Popf: "ℙ",
    Pr: "⪻",
    Precedes: "≺",
    PrecedesEqual: "⪯",
    PrecedesSlantEqual: "≼",
    PrecedesTilde: "≾",
    Prime: "″",
    Product: "∏",
    Proportion: "∷",
    Proportional: "∝",
    Pscr: "𝒫",
    Psi: "Ψ",
    QUOT: '"',
    Qfr: "𝔔",
    Qopf: "ℚ",
    Qscr: "𝒬",
    RBarr: "⤐",
    REG: "®",
    Racute: "Ŕ",
    Rang: "⟫",
    Rarr: "↠",
    Rarrtl: "⤖",
    Rcaron: "Ř",
    Rcedil: "Ŗ",
    Rcy: "Р",
    Re: "ℜ",
    ReverseElement: "∋",
    ReverseEquilibrium: "⇋",
    ReverseUpEquilibrium: "⥯",
    Rfr: "ℜ",
    Rho: "Ρ",
    RightAngleBracket: "⟩",
    RightArrow: "→",
    RightArrowBar: "⇥",
    RightArrowLeftArrow: "⇄",
    RightCeiling: "⌉",
    RightDoubleBracket: "⟧",
    RightDownTeeVector: "⥝",
    RightDownVector: "⇂",
    RightDownVectorBar: "⥕",
    RightFloor: "⌋",
    RightTee: "⊢",
    RightTeeArrow: "↦",
    RightTeeVector: "⥛",
    RightTriangle: "⊳",
    RightTriangleBar: "⧐",
    RightTriangleEqual: "⊵",
    RightUpDownVector: "⥏",
    RightUpTeeVector: "⥜",
    RightUpVector: "↾",
    RightUpVectorBar: "⥔",
    RightVector: "⇀",
    RightVectorBar: "⥓",
    Rightarrow: "⇒",
    Ropf: "ℝ",
    RoundImplies: "⥰",
    Rrightarrow: "⇛",
    Rscr: "ℛ",
    Rsh: "↱",
    RuleDelayed: "⧴",
    SHCHcy: "Щ",
    SHcy: "Ш",
    SOFTcy: "Ь",
    Sacute: "Ś",
    Sc: "⪼",
    Scaron: "Š",
    Scedil: "Ş",
    Scirc: "Ŝ",
    Scy: "С",
    Sfr: "𝔖",
    ShortDownArrow: "↓",
    ShortLeftArrow: "←",
    ShortRightArrow: "→",
    ShortUpArrow: "↑",
    Sigma: "Σ",
    SmallCircle: "∘",
    Sopf: "𝕊",
    Sqrt: "√",
    Square: "□",
    SquareIntersection: "⊓",
    SquareSubset: "⊏",
    SquareSubsetEqual: "⊑",
    SquareSuperset: "⊐",
    SquareSupersetEqual: "⊒",
    SquareUnion: "⊔",
    Sscr: "𝒮",
    Star: "⋆",
    Sub: "⋐",
    Subset: "⋐",
    SubsetEqual: "⊆",
    Succeeds: "≻",
    SucceedsEqual: "⪰",
    SucceedsSlantEqual: "≽",
    SucceedsTilde: "≿",
    SuchThat: "∋",
    Sum: "∑",
    Sup: "⋑",
    Superset: "⊃",
    SupersetEqual: "⊇",
    Supset: "⋑",
    THORN: "Þ",
    TRADE: "™",
    TSHcy: "Ћ",
    TScy: "Ц",
    Tab: "	",
    Tau: "Τ",
    Tcaron: "Ť",
    Tcedil: "Ţ",
    Tcy: "Т",
    Tfr: "𝔗",
    Therefore: "∴",
    Theta: "Θ",
    ThickSpace: "  ",
    ThinSpace: " ",
    Tilde: "∼",
    TildeEqual: "≃",
    TildeFullEqual: "≅",
    TildeTilde: "≈",
    Topf: "𝕋",
    TripleDot: "⃛",
    Tscr: "𝒯",
    Tstrok: "Ŧ",
    Uacute: "Ú",
    Uarr: "↟",
    Uarrocir: "⥉",
    Ubrcy: "Ў",
    Ubreve: "Ŭ",
    Ucirc: "Û",
    Ucy: "У",
    Udblac: "Ű",
    Ufr: "𝔘",
    Ugrave: "Ù",
    Umacr: "Ū",
    UnderBar: "_",
    UnderBrace: "⏟",
    UnderBracket: "⎵",
    UnderParenthesis: "⏝",
    Union: "⋃",
    UnionPlus: "⊎",
    Uogon: "Ų",
    Uopf: "𝕌",
    UpArrow: "↑",
    UpArrowBar: "⤒",
    UpArrowDownArrow: "⇅",
    UpDownArrow: "↕",
    UpEquilibrium: "⥮",
    UpTee: "⊥",
    UpTeeArrow: "↥",
    Uparrow: "⇑",
    Updownarrow: "⇕",
    UpperLeftArrow: "↖",
    UpperRightArrow: "↗",
    Upsi: "ϒ",
    Upsilon: "Υ",
    Uring: "Ů",
    Uscr: "𝒰",
    Utilde: "Ũ",
    Uuml: "Ü",
    VDash: "⊫",
    Vbar: "⫫",
    Vcy: "В",
    Vdash: "⊩",
    Vdashl: "⫦",
    Vee: "⋁",
    Verbar: "‖",
    Vert: "‖",
    VerticalBar: "∣",
    VerticalLine: "|",
    VerticalSeparator: "❘",
    VerticalTilde: "≀",
    VeryThinSpace: " ",
    Vfr: "𝔙",
    Vopf: "𝕍",
    Vscr: "𝒱",
    Vvdash: "⊪",
    Wcirc: "Ŵ",
    Wedge: "⋀",
    Wfr: "𝔚",
    Wopf: "𝕎",
    Wscr: "𝒲",
    Xfr: "𝔛",
    Xi: "Ξ",
    Xopf: "𝕏",
    Xscr: "𝒳",
    YAcy: "Я",
    YIcy: "Ї",
    YUcy: "Ю",
    Yacute: "Ý",
    Ycirc: "Ŷ",
    Ycy: "Ы",
    Yfr: "𝔜",
    Yopf: "𝕐",
    Yscr: "𝒴",
    Yuml: "Ÿ",
    ZHcy: "Ж",
    Zacute: "Ź",
    Zcaron: "Ž",
    Zcy: "З",
    Zdot: "Ż",
    ZeroWidthSpace: "​",
    Zeta: "Ζ",
    Zfr: "ℨ",
    Zopf: "ℤ",
    Zscr: "𝒵",
    aacute: "á",
    abreve: "ă",
    ac: "∾",
    acE: "∾̳",
    acd: "∿",
    acirc: "â",
    acute: "´",
    acy: "а",
    aelig: "æ",
    af: "⁡",
    afr: "𝔞",
    agrave: "à",
    alefsym: "ℵ",
    aleph: "ℵ",
    alpha: "α",
    amacr: "ā",
    amalg: "⨿",
    amp: "&",
    and: "∧",
    andand: "⩕",
    andd: "⩜",
    andslope: "⩘",
    andv: "⩚",
    ang: "∠",
    ange: "⦤",
    angle: "∠",
    angmsd: "∡",
    angmsdaa: "⦨",
    angmsdab: "⦩",
    angmsdac: "⦪",
    angmsdad: "⦫",
    angmsdae: "⦬",
    angmsdaf: "⦭",
    angmsdag: "⦮",
    angmsdah: "⦯",
    angrt: "∟",
    angrtvb: "⊾",
    angrtvbd: "⦝",
    angsph: "∢",
    angst: "Å",
    angzarr: "⍼",
    aogon: "ą",
    aopf: "𝕒",
    ap: "≈",
    apE: "⩰",
    apacir: "⩯",
    ape: "≊",
    apid: "≋",
    apos: "'",
    approx: "≈",
    approxeq: "≊",
    aring: "å",
    ascr: "𝒶",
    ast: "*",
    asymp: "≈",
    asympeq: "≍",
    atilde: "ã",
    auml: "ä",
    awconint: "∳",
    awint: "⨑",
    bNot: "⫭",
    backcong: "≌",
    backepsilon: "϶",
    backprime: "‵",
    backsim: "∽",
    backsimeq: "⋍",
    barvee: "⊽",
    barwed: "⌅",
    barwedge: "⌅",
    bbrk: "⎵",
    bbrktbrk: "⎶",
    bcong: "≌",
    bcy: "б",
    bdquo: "„",
    becaus: "∵",
    because: "∵",
    bemptyv: "⦰",
    bepsi: "϶",
    bernou: "ℬ",
    beta: "β",
    beth: "ℶ",
    between: "≬",
    bfr: "𝔟",
    bigcap: "⋂",
    bigcirc: "◯",
    bigcup: "⋃",
    bigodot: "⨀",
    bigoplus: "⨁",
    bigotimes: "⨂",
    bigsqcup: "⨆",
    bigstar: "★",
    bigtriangledown: "▽",
    bigtriangleup: "△",
    biguplus: "⨄",
    bigvee: "⋁",
    bigwedge: "⋀",
    bkarow: "⤍",
    blacklozenge: "⧫",
    blacksquare: "▪",
    blacktriangle: "▴",
    blacktriangledown: "▾",
    blacktriangleleft: "◂",
    blacktriangleright: "▸",
    blank: "␣",
    blk12: "▒",
    blk14: "░",
    blk34: "▓",
    block: "█",
    bne: "=⃥",
    bnequiv: "≡⃥",
    bnot: "⌐",
    bopf: "𝕓",
    bot: "⊥",
    bottom: "⊥",
    bowtie: "⋈",
    boxDL: "╗",
    boxDR: "╔",
    boxDl: "╖",
    boxDr: "╓",
    boxH: "═",
    boxHD: "╦",
    boxHU: "╩",
    boxHd: "╤",
    boxHu: "╧",
    boxUL: "╝",
    boxUR: "╚",
    boxUl: "╜",
    boxUr: "╙",
    boxV: "║",
    boxVH: "╬",
    boxVL: "╣",
    boxVR: "╠",
    boxVh: "╫",
    boxVl: "╢",
    boxVr: "╟",
    boxbox: "⧉",
    boxdL: "╕",
    boxdR: "╒",
    boxdl: "┐",
    boxdr: "┌",
    boxh: "─",
    boxhD: "╥",
    boxhU: "╨",
    boxhd: "┬",
    boxhu: "┴",
    boxminus: "⊟",
    boxplus: "⊞",
    boxtimes: "⊠",
    boxuL: "╛",
    boxuR: "╘",
    boxul: "┘",
    boxur: "└",
    boxv: "│",
    boxvH: "╪",
    boxvL: "╡",
    boxvR: "╞",
    boxvh: "┼",
    boxvl: "┤",
    boxvr: "├",
    bprime: "‵",
    breve: "˘",
    brvbar: "¦",
    bscr: "𝒷",
    bsemi: "⁏",
    bsim: "∽",
    bsime: "⋍",
    bsol: "\\",
    bsolb: "⧅",
    bsolhsub: "⟈",
    bull: "•",
    bullet: "•",
    bump: "≎",
    bumpE: "⪮",
    bumpe: "≏",
    bumpeq: "≏",
    cacute: "ć",
    cap: "∩",
    capand: "⩄",
    capbrcup: "⩉",
    capcap: "⩋",
    capcup: "⩇",
    capdot: "⩀",
    caps: "∩︀",
    caret: "⁁",
    caron: "ˇ",
    ccaps: "⩍",
    ccaron: "č",
    ccedil: "ç",
    ccirc: "ĉ",
    ccups: "⩌",
    ccupssm: "⩐",
    cdot: "ċ",
    cedil: "¸",
    cemptyv: "⦲",
    cent: "¢",
    centerdot: "·",
    cfr: "𝔠",
    chcy: "ч",
    check: "✓",
    checkmark: "✓",
    chi: "χ",
    cir: "○",
    cirE: "⧃",
    circ: "ˆ",
    circeq: "≗",
    circlearrowleft: "↺",
    circlearrowright: "↻",
    circledR: "®",
    circledS: "Ⓢ",
    circledast: "⊛",
    circledcirc: "⊚",
    circleddash: "⊝",
    cire: "≗",
    cirfnint: "⨐",
    cirmid: "⫯",
    cirscir: "⧂",
    clubs: "♣",
    clubsuit: "♣",
    colon: ":",
    colone: "≔",
    coloneq: "≔",
    comma: ",",
    commat: "@",
    comp: "∁",
    compfn: "∘",
    complement: "∁",
    complexes: "ℂ",
    cong: "≅",
    congdot: "⩭",
    conint: "∮",
    copf: "𝕔",
    coprod: "∐",
    copy: "©",
    copysr: "℗",
    crarr: "↵",
    cross: "✗",
    cscr: "𝒸",
    csub: "⫏",
    csube: "⫑",
    csup: "⫐",
    csupe: "⫒",
    ctdot: "⋯",
    cudarrl: "⤸",
    cudarrr: "⤵",
    cuepr: "⋞",
    cuesc: "⋟",
    cularr: "↶",
    cularrp: "⤽",
    cup: "∪",
    cupbrcap: "⩈",
    cupcap: "⩆",
    cupcup: "⩊",
    cupdot: "⊍",
    cupor: "⩅",
    cups: "∪︀",
    curarr: "↷",
    curarrm: "⤼",
    curlyeqprec: "⋞",
    curlyeqsucc: "⋟",
    curlyvee: "⋎",
    curlywedge: "⋏",
    curren: "¤",
    curvearrowleft: "↶",
    curvearrowright: "↷",
    cuvee: "⋎",
    cuwed: "⋏",
    cwconint: "∲",
    cwint: "∱",
    cylcty: "⌭",
    dArr: "⇓",
    dHar: "⥥",
    dagger: "†",
    daleth: "ℸ",
    darr: "↓",
    dash: "‐",
    dashv: "⊣",
    dbkarow: "⤏",
    dblac: "˝",
    dcaron: "ď",
    dcy: "д",
    dd: "ⅆ",
    ddagger: "‡",
    ddarr: "⇊",
    ddotseq: "⩷",
    deg: "°",
    delta: "δ",
    demptyv: "⦱",
    dfisht: "⥿",
    dfr: "𝔡",
    dharl: "⇃",
    dharr: "⇂",
    diam: "⋄",
    diamond: "⋄",
    diamondsuit: "♦",
    diams: "♦",
    die: "¨",
    digamma: "ϝ",
    disin: "⋲",
    div: "÷",
    divide: "÷",
    divideontimes: "⋇",
    divonx: "⋇",
    djcy: "ђ",
    dlcorn: "⌞",
    dlcrop: "⌍",
    dollar: "$",
    dopf: "𝕕",
    dot: "˙",
    doteq: "≐",
    doteqdot: "≑",
    dotminus: "∸",
    dotplus: "∔",
    dotsquare: "⊡",
    doublebarwedge: "⌆",
    downarrow: "↓",
    downdownarrows: "⇊",
    downharpoonleft: "⇃",
    downharpoonright: "⇂",
    drbkarow: "⤐",
    drcorn: "⌟",
    drcrop: "⌌",
    dscr: "𝒹",
    dscy: "ѕ",
    dsol: "⧶",
    dstrok: "đ",
    dtdot: "⋱",
    dtri: "▿",
    dtrif: "▾",
    duarr: "⇵",
    duhar: "⥯",
    dwangle: "⦦",
    dzcy: "џ",
    dzigrarr: "⟿",
    eDDot: "⩷",
    eDot: "≑",
    eacute: "é",
    easter: "⩮",
    ecaron: "ě",
    ecir: "≖",
    ecirc: "ê",
    ecolon: "≕",
    ecy: "э",
    edot: "ė",
    ee: "ⅇ",
    efDot: "≒",
    efr: "𝔢",
    eg: "⪚",
    egrave: "è",
    egs: "⪖",
    egsdot: "⪘",
    el: "⪙",
    elinters: "⏧",
    ell: "ℓ",
    els: "⪕",
    elsdot: "⪗",
    emacr: "ē",
    empty: "∅",
    emptyset: "∅",
    emptyv: "∅",
    emsp13: " ",
    emsp14: " ",
    emsp: " ",
    eng: "ŋ",
    ensp: " ",
    eogon: "ę",
    eopf: "𝕖",
    epar: "⋕",
    eparsl: "⧣",
    eplus: "⩱",
    epsi: "ε",
    epsilon: "ε",
    epsiv: "ϵ",
    eqcirc: "≖",
    eqcolon: "≕",
    eqsim: "≂",
    eqslantgtr: "⪖",
    eqslantless: "⪕",
    equals: "=",
    equest: "≟",
    equiv: "≡",
    equivDD: "⩸",
    eqvparsl: "⧥",
    erDot: "≓",
    erarr: "⥱",
    escr: "ℯ",
    esdot: "≐",
    esim: "≂",
    eta: "η",
    eth: "ð",
    euml: "ë",
    euro: "€",
    excl: "!",
    exist: "∃",
    expectation: "ℰ",
    exponentiale: "ⅇ",
    fallingdotseq: "≒",
    fcy: "ф",
    female: "♀",
    ffilig: "ﬃ",
    fflig: "ﬀ",
    ffllig: "ﬄ",
    ffr: "𝔣",
    filig: "ﬁ",
    fjlig: "fj",
    flat: "♭",
    fllig: "ﬂ",
    fltns: "▱",
    fnof: "ƒ",
    fopf: "𝕗",
    forall: "∀",
    fork: "⋔",
    forkv: "⫙",
    fpartint: "⨍",
    frac12: "½",
    frac13: "⅓",
    frac14: "¼",
    frac15: "⅕",
    frac16: "⅙",
    frac18: "⅛",
    frac23: "⅔",
    frac25: "⅖",
    frac34: "¾",
    frac35: "⅗",
    frac38: "⅜",
    frac45: "⅘",
    frac56: "⅚",
    frac58: "⅝",
    frac78: "⅞",
    frasl: "⁄",
    frown: "⌢",
    fscr: "𝒻",
    gE: "≧",
    gEl: "⪌",
    gacute: "ǵ",
    gamma: "γ",
    gammad: "ϝ",
    gap: "⪆",
    gbreve: "ğ",
    gcirc: "ĝ",
    gcy: "г",
    gdot: "ġ",
    ge: "≥",
    gel: "⋛",
    geq: "≥",
    geqq: "≧",
    geqslant: "⩾",
    ges: "⩾",
    gescc: "⪩",
    gesdot: "⪀",
    gesdoto: "⪂",
    gesdotol: "⪄",
    gesl: "⋛︀",
    gesles: "⪔",
    gfr: "𝔤",
    gg: "≫",
    ggg: "⋙",
    gimel: "ℷ",
    gjcy: "ѓ",
    gl: "≷",
    glE: "⪒",
    gla: "⪥",
    glj: "⪤",
    gnE: "≩",
    gnap: "⪊",
    gnapprox: "⪊",
    gne: "⪈",
    gneq: "⪈",
    gneqq: "≩",
    gnsim: "⋧",
    gopf: "𝕘",
    grave: "`",
    gscr: "ℊ",
    gsim: "≳",
    gsime: "⪎",
    gsiml: "⪐",
    gt: ">",
    gtcc: "⪧",
    gtcir: "⩺",
    gtdot: "⋗",
    gtlPar: "⦕",
    gtquest: "⩼",
    gtrapprox: "⪆",
    gtrarr: "⥸",
    gtrdot: "⋗",
    gtreqless: "⋛",
    gtreqqless: "⪌",
    gtrless: "≷",
    gtrsim: "≳",
    gvertneqq: "≩︀",
    gvnE: "≩︀",
    hArr: "⇔",
    hairsp: " ",
    half: "½",
    hamilt: "ℋ",
    hardcy: "ъ",
    harr: "↔",
    harrcir: "⥈",
    harrw: "↭",
    hbar: "ℏ",
    hcirc: "ĥ",
    hearts: "♥",
    heartsuit: "♥",
    hellip: "…",
    hercon: "⊹",
    hfr: "𝔥",
    hksearow: "⤥",
    hkswarow: "⤦",
    hoarr: "⇿",
    homtht: "∻",
    hookleftarrow: "↩",
    hookrightarrow: "↪",
    hopf: "𝕙",
    horbar: "―",
    hscr: "𝒽",
    hslash: "ℏ",
    hstrok: "ħ",
    hybull: "⁃",
    hyphen: "‐",
    iacute: "í",
    ic: "⁣",
    icirc: "î",
    icy: "и",
    iecy: "е",
    iexcl: "¡",
    iff: "⇔",
    ifr: "𝔦",
    igrave: "ì",
    ii: "ⅈ",
    iiiint: "⨌",
    iiint: "∭",
    iinfin: "⧜",
    iiota: "℩",
    ijlig: "ĳ",
    imacr: "ī",
    image: "ℑ",
    imagline: "ℐ",
    imagpart: "ℑ",
    imath: "ı",
    imof: "⊷",
    imped: "Ƶ",
    in: "∈",
    incare: "℅",
    infin: "∞",
    infintie: "⧝",
    inodot: "ı",
    int: "∫",
    intcal: "⊺",
    integers: "ℤ",
    intercal: "⊺",
    intlarhk: "⨗",
    intprod: "⨼",
    iocy: "ё",
    iogon: "į",
    iopf: "𝕚",
    iota: "ι",
    iprod: "⨼",
    iquest: "¿",
    iscr: "𝒾",
    isin: "∈",
    isinE: "⋹",
    isindot: "⋵",
    isins: "⋴",
    isinsv: "⋳",
    isinv: "∈",
    it: "⁢",
    itilde: "ĩ",
    iukcy: "і",
    iuml: "ï",
    jcirc: "ĵ",
    jcy: "й",
    jfr: "𝔧",
    jmath: "ȷ",
    jopf: "𝕛",
    jscr: "𝒿",
    jsercy: "ј",
    jukcy: "є",
    kappa: "κ",
    kappav: "ϰ",
    kcedil: "ķ",
    kcy: "к",
    kfr: "𝔨",
    kgreen: "ĸ",
    khcy: "х",
    kjcy: "ќ",
    kopf: "𝕜",
    kscr: "𝓀",
    lAarr: "⇚",
    lArr: "⇐",
    lAtail: "⤛",
    lBarr: "⤎",
    lE: "≦",
    lEg: "⪋",
    lHar: "⥢",
    lacute: "ĺ",
    laemptyv: "⦴",
    lagran: "ℒ",
    lambda: "λ",
    lang: "⟨",
    langd: "⦑",
    langle: "⟨",
    lap: "⪅",
    laquo: "«",
    larr: "←",
    larrb: "⇤",
    larrbfs: "⤟",
    larrfs: "⤝",
    larrhk: "↩",
    larrlp: "↫",
    larrpl: "⤹",
    larrsim: "⥳",
    larrtl: "↢",
    lat: "⪫",
    latail: "⤙",
    late: "⪭",
    lates: "⪭︀",
    lbarr: "⤌",
    lbbrk: "❲",
    lbrace: "{",
    lbrack: "[",
    lbrke: "⦋",
    lbrksld: "⦏",
    lbrkslu: "⦍",
    lcaron: "ľ",
    lcedil: "ļ",
    lceil: "⌈",
    lcub: "{",
    lcy: "л",
    ldca: "⤶",
    ldquo: "“",
    ldquor: "„",
    ldrdhar: "⥧",
    ldrushar: "⥋",
    ldsh: "↲",
    le: "≤",
    leftarrow: "←",
    leftarrowtail: "↢",
    leftharpoondown: "↽",
    leftharpoonup: "↼",
    leftleftarrows: "⇇",
    leftrightarrow: "↔",
    leftrightarrows: "⇆",
    leftrightharpoons: "⇋",
    leftrightsquigarrow: "↭",
    leftthreetimes: "⋋",
    leg: "⋚",
    leq: "≤",
    leqq: "≦",
    leqslant: "⩽",
    les: "⩽",
    lescc: "⪨",
    lesdot: "⩿",
    lesdoto: "⪁",
    lesdotor: "⪃",
    lesg: "⋚︀",
    lesges: "⪓",
    lessapprox: "⪅",
    lessdot: "⋖",
    lesseqgtr: "⋚",
    lesseqqgtr: "⪋",
    lessgtr: "≶",
    lesssim: "≲",
    lfisht: "⥼",
    lfloor: "⌊",
    lfr: "𝔩",
    lg: "≶",
    lgE: "⪑",
    lhard: "↽",
    lharu: "↼",
    lharul: "⥪",
    lhblk: "▄",
    ljcy: "љ",
    ll: "≪",
    llarr: "⇇",
    llcorner: "⌞",
    llhard: "⥫",
    lltri: "◺",
    lmidot: "ŀ",
    lmoust: "⎰",
    lmoustache: "⎰",
    lnE: "≨",
    lnap: "⪉",
    lnapprox: "⪉",
    lne: "⪇",
    lneq: "⪇",
    lneqq: "≨",
    lnsim: "⋦",
    loang: "⟬",
    loarr: "⇽",
    lobrk: "⟦",
    longleftarrow: "⟵",
    longleftrightarrow: "⟷",
    longmapsto: "⟼",
    longrightarrow: "⟶",
    looparrowleft: "↫",
    looparrowright: "↬",
    lopar: "⦅",
    lopf: "𝕝",
    loplus: "⨭",
    lotimes: "⨴",
    lowast: "∗",
    lowbar: "_",
    loz: "◊",
    lozenge: "◊",
    lozf: "⧫",
    lpar: "(",
    lparlt: "⦓",
    lrarr: "⇆",
    lrcorner: "⌟",
    lrhar: "⇋",
    lrhard: "⥭",
    lrm: "‎",
    lrtri: "⊿",
    lsaquo: "‹",
    lscr: "𝓁",
    lsh: "↰",
    lsim: "≲",
    lsime: "⪍",
    lsimg: "⪏",
    lsqb: "[",
    lsquo: "‘",
    lsquor: "‚",
    lstrok: "ł",
    lt: "<",
    ltcc: "⪦",
    ltcir: "⩹",
    ltdot: "⋖",
    lthree: "⋋",
    ltimes: "⋉",
    ltlarr: "⥶",
    ltquest: "⩻",
    ltrPar: "⦖",
    ltri: "◃",
    ltrie: "⊴",
    ltrif: "◂",
    lurdshar: "⥊",
    luruhar: "⥦",
    lvertneqq: "≨︀",
    lvnE: "≨︀",
    mDDot: "∺",
    macr: "¯",
    male: "♂",
    malt: "✠",
    maltese: "✠",
    map: "↦",
    mapsto: "↦",
    mapstodown: "↧",
    mapstoleft: "↤",
    mapstoup: "↥",
    marker: "▮",
    mcomma: "⨩",
    mcy: "м",
    mdash: "—",
    measuredangle: "∡",
    mfr: "𝔪",
    mho: "℧",
    micro: "µ",
    mid: "∣",
    midast: "*",
    midcir: "⫰",
    middot: "·",
    minus: "−",
    minusb: "⊟",
    minusd: "∸",
    minusdu: "⨪",
    mlcp: "⫛",
    mldr: "…",
    mnplus: "∓",
    models: "⊧",
    mopf: "𝕞",
    mp: "∓",
    mscr: "𝓂",
    mstpos: "∾",
    mu: "μ",
    multimap: "⊸",
    mumap: "⊸",
    nGg: "⋙̸",
    nGt: "≫⃒",
    nGtv: "≫̸",
    nLeftarrow: "⇍",
    nLeftrightarrow: "⇎",
    nLl: "⋘̸",
    nLt: "≪⃒",
    nLtv: "≪̸",
    nRightarrow: "⇏",
    nVDash: "⊯",
    nVdash: "⊮",
    nabla: "∇",
    nacute: "ń",
    nang: "∠⃒",
    nap: "≉",
    napE: "⩰̸",
    napid: "≋̸",
    napos: "ŉ",
    napprox: "≉",
    natur: "♮",
    natural: "♮",
    naturals: "ℕ",
    nbsp: " ",
    nbump: "≎̸",
    nbumpe: "≏̸",
    ncap: "⩃",
    ncaron: "ň",
    ncedil: "ņ",
    ncong: "≇",
    ncongdot: "⩭̸",
    ncup: "⩂",
    ncy: "н",
    ndash: "–",
    ne: "≠",
    neArr: "⇗",
    nearhk: "⤤",
    nearr: "↗",
    nearrow: "↗",
    nedot: "≐̸",
    nequiv: "≢",
    nesear: "⤨",
    nesim: "≂̸",
    nexist: "∄",
    nexists: "∄",
    nfr: "𝔫",
    ngE: "≧̸",
    nge: "≱",
    ngeq: "≱",
    ngeqq: "≧̸",
    ngeqslant: "⩾̸",
    nges: "⩾̸",
    ngsim: "≵",
    ngt: "≯",
    ngtr: "≯",
    nhArr: "⇎",
    nharr: "↮",
    nhpar: "⫲",
    ni: "∋",
    nis: "⋼",
    nisd: "⋺",
    niv: "∋",
    njcy: "њ",
    nlArr: "⇍",
    nlE: "≦̸",
    nlarr: "↚",
    nldr: "‥",
    nle: "≰",
    nleftarrow: "↚",
    nleftrightarrow: "↮",
    nleq: "≰",
    nleqq: "≦̸",
    nleqslant: "⩽̸",
    nles: "⩽̸",
    nless: "≮",
    nlsim: "≴",
    nlt: "≮",
    nltri: "⋪",
    nltrie: "⋬",
    nmid: "∤",
    nopf: "𝕟",
    not: "¬",
    notin: "∉",
    notinE: "⋹̸",
    notindot: "⋵̸",
    notinva: "∉",
    notinvb: "⋷",
    notinvc: "⋶",
    notni: "∌",
    notniva: "∌",
    notnivb: "⋾",
    notnivc: "⋽",
    npar: "∦",
    nparallel: "∦",
    nparsl: "⫽⃥",
    npart: "∂̸",
    npolint: "⨔",
    npr: "⊀",
    nprcue: "⋠",
    npre: "⪯̸",
    nprec: "⊀",
    npreceq: "⪯̸",
    nrArr: "⇏",
    nrarr: "↛",
    nrarrc: "⤳̸",
    nrarrw: "↝̸",
    nrightarrow: "↛",
    nrtri: "⋫",
    nrtrie: "⋭",
    nsc: "⊁",
    nsccue: "⋡",
    nsce: "⪰̸",
    nscr: "𝓃",
    nshortmid: "∤",
    nshortparallel: "∦",
    nsim: "≁",
    nsime: "≄",
    nsimeq: "≄",
    nsmid: "∤",
    nspar: "∦",
    nsqsube: "⋢",
    nsqsupe: "⋣",
    nsub: "⊄",
    nsubE: "⫅̸",
    nsube: "⊈",
    nsubset: "⊂⃒",
    nsubseteq: "⊈",
    nsubseteqq: "⫅̸",
    nsucc: "⊁",
    nsucceq: "⪰̸",
    nsup: "⊅",
    nsupE: "⫆̸",
    nsupe: "⊉",
    nsupset: "⊃⃒",
    nsupseteq: "⊉",
    nsupseteqq: "⫆̸",
    ntgl: "≹",
    ntilde: "ñ",
    ntlg: "≸",
    ntriangleleft: "⋪",
    ntrianglelefteq: "⋬",
    ntriangleright: "⋫",
    ntrianglerighteq: "⋭",
    nu: "ν",
    num: "#",
    numero: "№",
    numsp: " ",
    nvDash: "⊭",
    nvHarr: "⤄",
    nvap: "≍⃒",
    nvdash: "⊬",
    nvge: "≥⃒",
    nvgt: ">⃒",
    nvinfin: "⧞",
    nvlArr: "⤂",
    nvle: "≤⃒",
    nvlt: "<⃒",
    nvltrie: "⊴⃒",
    nvrArr: "⤃",
    nvrtrie: "⊵⃒",
    nvsim: "∼⃒",
    nwArr: "⇖",
    nwarhk: "⤣",
    nwarr: "↖",
    nwarrow: "↖",
    nwnear: "⤧",
    oS: "Ⓢ",
    oacute: "ó",
    oast: "⊛",
    ocir: "⊚",
    ocirc: "ô",
    ocy: "о",
    odash: "⊝",
    odblac: "ő",
    odiv: "⨸",
    odot: "⊙",
    odsold: "⦼",
    oelig: "œ",
    ofcir: "⦿",
    ofr: "𝔬",
    ogon: "˛",
    ograve: "ò",
    ogt: "⧁",
    ohbar: "⦵",
    ohm: "Ω",
    oint: "∮",
    olarr: "↺",
    olcir: "⦾",
    olcross: "⦻",
    oline: "‾",
    olt: "⧀",
    omacr: "ō",
    omega: "ω",
    omicron: "ο",
    omid: "⦶",
    ominus: "⊖",
    oopf: "𝕠",
    opar: "⦷",
    operp: "⦹",
    oplus: "⊕",
    or: "∨",
    orarr: "↻",
    ord: "⩝",
    order: "ℴ",
    orderof: "ℴ",
    ordf: "ª",
    ordm: "º",
    origof: "⊶",
    oror: "⩖",
    orslope: "⩗",
    orv: "⩛",
    oscr: "ℴ",
    oslash: "ø",
    osol: "⊘",
    otilde: "õ",
    otimes: "⊗",
    otimesas: "⨶",
    ouml: "ö",
    ovbar: "⌽",
    par: "∥",
    para: "¶",
    parallel: "∥",
    parsim: "⫳",
    parsl: "⫽",
    part: "∂",
    pcy: "п",
    percnt: "%",
    period: ".",
    permil: "‰",
    perp: "⊥",
    pertenk: "‱",
    pfr: "𝔭",
    phi: "φ",
    phiv: "ϕ",
    phmmat: "ℳ",
    phone: "☎",
    pi: "π",
    pitchfork: "⋔",
    piv: "ϖ",
    planck: "ℏ",
    planckh: "ℎ",
    plankv: "ℏ",
    plus: "+",
    plusacir: "⨣",
    plusb: "⊞",
    pluscir: "⨢",
    plusdo: "∔",
    plusdu: "⨥",
    pluse: "⩲",
    plusmn: "±",
    plussim: "⨦",
    plustwo: "⨧",
    pm: "±",
    pointint: "⨕",
    popf: "𝕡",
    pound: "£",
    pr: "≺",
    prE: "⪳",
    prap: "⪷",
    prcue: "≼",
    pre: "⪯",
    prec: "≺",
    precapprox: "⪷",
    preccurlyeq: "≼",
    preceq: "⪯",
    precnapprox: "⪹",
    precneqq: "⪵",
    precnsim: "⋨",
    precsim: "≾",
    prime: "′",
    primes: "ℙ",
    prnE: "⪵",
    prnap: "⪹",
    prnsim: "⋨",
    prod: "∏",
    profalar: "⌮",
    profline: "⌒",
    profsurf: "⌓",
    prop: "∝",
    propto: "∝",
    prsim: "≾",
    prurel: "⊰",
    pscr: "𝓅",
    psi: "ψ",
    puncsp: " ",
    qfr: "𝔮",
    qint: "⨌",
    qopf: "𝕢",
    qprime: "⁗",
    qscr: "𝓆",
    quaternions: "ℍ",
    quatint: "⨖",
    quest: "?",
    questeq: "≟",
    quot: '"',
    rAarr: "⇛",
    rArr: "⇒",
    rAtail: "⤜",
    rBarr: "⤏",
    rHar: "⥤",
    race: "∽̱",
    racute: "ŕ",
    radic: "√",
    raemptyv: "⦳",
    rang: "⟩",
    rangd: "⦒",
    range: "⦥",
    rangle: "⟩",
    raquo: "»",
    rarr: "→",
    rarrap: "⥵",
    rarrb: "⇥",
    rarrbfs: "⤠",
    rarrc: "⤳",
    rarrfs: "⤞",
    rarrhk: "↪",
    rarrlp: "↬",
    rarrpl: "⥅",
    rarrsim: "⥴",
    rarrtl: "↣",
    rarrw: "↝",
    ratail: "⤚",
    ratio: "∶",
    rationals: "ℚ",
    rbarr: "⤍",
    rbbrk: "❳",
    rbrace: "}",
    rbrack: "]",
    rbrke: "⦌",
    rbrksld: "⦎",
    rbrkslu: "⦐",
    rcaron: "ř",
    rcedil: "ŗ",
    rceil: "⌉",
    rcub: "}",
    rcy: "р",
    rdca: "⤷",
    rdldhar: "⥩",
    rdquo: "”",
    rdquor: "”",
    rdsh: "↳",
    real: "ℜ",
    realine: "ℛ",
    realpart: "ℜ",
    reals: "ℝ",
    rect: "▭",
    reg: "®",
    rfisht: "⥽",
    rfloor: "⌋",
    rfr: "𝔯",
    rhard: "⇁",
    rharu: "⇀",
    rharul: "⥬",
    rho: "ρ",
    rhov: "ϱ",
    rightarrow: "→",
    rightarrowtail: "↣",
    rightharpoondown: "⇁",
    rightharpoonup: "⇀",
    rightleftarrows: "⇄",
    rightleftharpoons: "⇌",
    rightrightarrows: "⇉",
    rightsquigarrow: "↝",
    rightthreetimes: "⋌",
    ring: "˚",
    risingdotseq: "≓",
    rlarr: "⇄",
    rlhar: "⇌",
    rlm: "‏",
    rmoust: "⎱",
    rmoustache: "⎱",
    rnmid: "⫮",
    roang: "⟭",
    roarr: "⇾",
    robrk: "⟧",
    ropar: "⦆",
    ropf: "𝕣",
    roplus: "⨮",
    rotimes: "⨵",
    rpar: ")",
    rpargt: "⦔",
    rppolint: "⨒",
    rrarr: "⇉",
    rsaquo: "›",
    rscr: "𝓇",
    rsh: "↱",
    rsqb: "]",
    rsquo: "’",
    rsquor: "’",
    rthree: "⋌",
    rtimes: "⋊",
    rtri: "▹",
    rtrie: "⊵",
    rtrif: "▸",
    rtriltri: "⧎",
    ruluhar: "⥨",
    rx: "℞",
    sacute: "ś",
    sbquo: "‚",
    sc: "≻",
    scE: "⪴",
    scap: "⪸",
    scaron: "š",
    sccue: "≽",
    sce: "⪰",
    scedil: "ş",
    scirc: "ŝ",
    scnE: "⪶",
    scnap: "⪺",
    scnsim: "⋩",
    scpolint: "⨓",
    scsim: "≿",
    scy: "с",
    sdot: "⋅",
    sdotb: "⊡",
    sdote: "⩦",
    seArr: "⇘",
    searhk: "⤥",
    searr: "↘",
    searrow: "↘",
    sect: "§",
    semi: ";",
    seswar: "⤩",
    setminus: "∖",
    setmn: "∖",
    sext: "✶",
    sfr: "𝔰",
    sfrown: "⌢",
    sharp: "♯",
    shchcy: "щ",
    shcy: "ш",
    shortmid: "∣",
    shortparallel: "∥",
    shy: "­",
    sigma: "σ",
    sigmaf: "ς",
    sigmav: "ς",
    sim: "∼",
    simdot: "⩪",
    sime: "≃",
    simeq: "≃",
    simg: "⪞",
    simgE: "⪠",
    siml: "⪝",
    simlE: "⪟",
    simne: "≆",
    simplus: "⨤",
    simrarr: "⥲",
    slarr: "←",
    smallsetminus: "∖",
    smashp: "⨳",
    smeparsl: "⧤",
    smid: "∣",
    smile: "⌣",
    smt: "⪪",
    smte: "⪬",
    smtes: "⪬︀",
    softcy: "ь",
    sol: "/",
    solb: "⧄",
    solbar: "⌿",
    sopf: "𝕤",
    spades: "♠",
    spadesuit: "♠",
    spar: "∥",
    sqcap: "⊓",
    sqcaps: "⊓︀",
    sqcup: "⊔",
    sqcups: "⊔︀",
    sqsub: "⊏",
    sqsube: "⊑",
    sqsubset: "⊏",
    sqsubseteq: "⊑",
    sqsup: "⊐",
    sqsupe: "⊒",
    sqsupset: "⊐",
    sqsupseteq: "⊒",
    squ: "□",
    square: "□",
    squarf: "▪",
    squf: "▪",
    srarr: "→",
    sscr: "𝓈",
    ssetmn: "∖",
    ssmile: "⌣",
    sstarf: "⋆",
    star: "☆",
    starf: "★",
    straightepsilon: "ϵ",
    straightphi: "ϕ",
    strns: "¯",
    sub: "⊂",
    subE: "⫅",
    subdot: "⪽",
    sube: "⊆",
    subedot: "⫃",
    submult: "⫁",
    subnE: "⫋",
    subne: "⊊",
    subplus: "⪿",
    subrarr: "⥹",
    subset: "⊂",
    subseteq: "⊆",
    subseteqq: "⫅",
    subsetneq: "⊊",
    subsetneqq: "⫋",
    subsim: "⫇",
    subsub: "⫕",
    subsup: "⫓",
    succ: "≻",
    succapprox: "⪸",
    succcurlyeq: "≽",
    succeq: "⪰",
    succnapprox: "⪺",
    succneqq: "⪶",
    succnsim: "⋩",
    succsim: "≿",
    sum: "∑",
    sung: "♪",
    sup1: "¹",
    sup2: "²",
    sup3: "³",
    sup: "⊃",
    supE: "⫆",
    supdot: "⪾",
    supdsub: "⫘",
    supe: "⊇",
    supedot: "⫄",
    suphsol: "⟉",
    suphsub: "⫗",
    suplarr: "⥻",
    supmult: "⫂",
    supnE: "⫌",
    supne: "⊋",
    supplus: "⫀",
    supset: "⊃",
    supseteq: "⊇",
    supseteqq: "⫆",
    supsetneq: "⊋",
    supsetneqq: "⫌",
    supsim: "⫈",
    supsub: "⫔",
    supsup: "⫖",
    swArr: "⇙",
    swarhk: "⤦",
    swarr: "↙",
    swarrow: "↙",
    swnwar: "⤪",
    szlig: "ß",
    target: "⌖",
    tau: "τ",
    tbrk: "⎴",
    tcaron: "ť",
    tcedil: "ţ",
    tcy: "т",
    tdot: "⃛",
    telrec: "⌕",
    tfr: "𝔱",
    there4: "∴",
    therefore: "∴",
    theta: "θ",
    thetasym: "ϑ",
    thetav: "ϑ",
    thickapprox: "≈",
    thicksim: "∼",
    thinsp: " ",
    thkap: "≈",
    thksim: "∼",
    thorn: "þ",
    tilde: "˜",
    times: "×",
    timesb: "⊠",
    timesbar: "⨱",
    timesd: "⨰",
    tint: "∭",
    toea: "⤨",
    top: "⊤",
    topbot: "⌶",
    topcir: "⫱",
    topf: "𝕥",
    topfork: "⫚",
    tosa: "⤩",
    tprime: "‴",
    trade: "™",
    triangle: "▵",
    triangledown: "▿",
    triangleleft: "◃",
    trianglelefteq: "⊴",
    triangleq: "≜",
    triangleright: "▹",
    trianglerighteq: "⊵",
    tridot: "◬",
    trie: "≜",
    triminus: "⨺",
    triplus: "⨹",
    trisb: "⧍",
    tritime: "⨻",
    trpezium: "⏢",
    tscr: "𝓉",
    tscy: "ц",
    tshcy: "ћ",
    tstrok: "ŧ",
    twixt: "≬",
    twoheadleftarrow: "↞",
    twoheadrightarrow: "↠",
    uArr: "⇑",
    uHar: "⥣",
    uacute: "ú",
    uarr: "↑",
    ubrcy: "ў",
    ubreve: "ŭ",
    ucirc: "û",
    ucy: "у",
    udarr: "⇅",
    udblac: "ű",
    udhar: "⥮",
    ufisht: "⥾",
    ufr: "𝔲",
    ugrave: "ù",
    uharl: "↿",
    uharr: "↾",
    uhblk: "▀",
    ulcorn: "⌜",
    ulcorner: "⌜",
    ulcrop: "⌏",
    ultri: "◸",
    umacr: "ū",
    uml: "¨",
    uogon: "ų",
    uopf: "𝕦",
    uparrow: "↑",
    updownarrow: "↕",
    upharpoonleft: "↿",
    upharpoonright: "↾",
    uplus: "⊎",
    upsi: "υ",
    upsih: "ϒ",
    upsilon: "υ",
    upuparrows: "⇈",
    urcorn: "⌝",
    urcorner: "⌝",
    urcrop: "⌎",
    uring: "ů",
    urtri: "◹",
    uscr: "𝓊",
    utdot: "⋰",
    utilde: "ũ",
    utri: "▵",
    utrif: "▴",
    uuarr: "⇈",
    uuml: "ü",
    uwangle: "⦧",
    vArr: "⇕",
    vBar: "⫨",
    vBarv: "⫩",
    vDash: "⊨",
    vangrt: "⦜",
    varepsilon: "ϵ",
    varkappa: "ϰ",
    varnothing: "∅",
    varphi: "ϕ",
    varpi: "ϖ",
    varpropto: "∝",
    varr: "↕",
    varrho: "ϱ",
    varsigma: "ς",
    varsubsetneq: "⊊︀",
    varsubsetneqq: "⫋︀",
    varsupsetneq: "⊋︀",
    varsupsetneqq: "⫌︀",
    vartheta: "ϑ",
    vartriangleleft: "⊲",
    vartriangleright: "⊳",
    vcy: "в",
    vdash: "⊢",
    vee: "∨",
    veebar: "⊻",
    veeeq: "≚",
    vellip: "⋮",
    verbar: "|",
    vert: "|",
    vfr: "𝔳",
    vltri: "⊲",
    vnsub: "⊂⃒",
    vnsup: "⊃⃒",
    vopf: "𝕧",
    vprop: "∝",
    vrtri: "⊳",
    vscr: "𝓋",
    vsubnE: "⫋︀",
    vsubne: "⊊︀",
    vsupnE: "⫌︀",
    vsupne: "⊋︀",
    vzigzag: "⦚",
    wcirc: "ŵ",
    wedbar: "⩟",
    wedge: "∧",
    wedgeq: "≙",
    weierp: "℘",
    wfr: "𝔴",
    wopf: "𝕨",
    wp: "℘",
    wr: "≀",
    wreath: "≀",
    wscr: "𝓌",
    xcap: "⋂",
    xcirc: "◯",
    xcup: "⋃",
    xdtri: "▽",
    xfr: "𝔵",
    xhArr: "⟺",
    xharr: "⟷",
    xi: "ξ",
    xlArr: "⟸",
    xlarr: "⟵",
    xmap: "⟼",
    xnis: "⋻",
    xodot: "⨀",
    xopf: "𝕩",
    xoplus: "⨁",
    xotime: "⨂",
    xrArr: "⟹",
    xrarr: "⟶",
    xscr: "𝓍",
    xsqcup: "⨆",
    xuplus: "⨄",
    xutri: "△",
    xvee: "⋁",
    xwedge: "⋀",
    yacute: "ý",
    yacy: "я",
    ycirc: "ŷ",
    ycy: "ы",
    yen: "¥",
    yfr: "𝔶",
    yicy: "ї",
    yopf: "𝕪",
    yscr: "𝓎",
    yucy: "ю",
    yuml: "ÿ",
    zacute: "ź",
    zcaron: "ž",
    zcy: "з",
    zdot: "ż",
    zeetrf: "ℨ",
    zeta: "ζ",
    zfr: "𝔷",
    zhcy: "ж",
    zigrarr: "⇝",
    zopf: "𝕫",
    zscr: "𝓏",
    zwj: "‍",
    zwnj: "‌"
  };

  // ../../../ui/node_modules/.pnpm/decode-named-character-reference@1.2.0/node_modules/decode-named-character-reference/index.js
  var own2 = {}.hasOwnProperty;
  function decodeNamedCharacterReference(value) {
    return own2.call(characterEntities, value) ? characterEntities[value] : false;
  }

  // ../../../ui/node_modules/.pnpm/parse-entities@4.0.2/node_modules/parse-entities/lib/index.js
  var messages = [
    "",
    /* 1: Non terminated (named) */
    "Named character references must be terminated by a semicolon",
    /* 2: Non terminated (numeric) */
    "Numeric character references must be terminated by a semicolon",
    /* 3: Empty (named) */
    "Named character references cannot be empty",
    /* 4: Empty (numeric) */
    "Numeric character references cannot be empty",
    /* 5: Unknown (named) */
    "Named character references must be known",
    /* 6: Disallowed (numeric) */
    "Numeric character references cannot be disallowed",
    /* 7: Prohibited (numeric) */
    "Numeric character references cannot be outside the permissible Unicode range"
  ];
  function parseEntities(value, options) {
    const settings = options || {};
    const additional = typeof settings.additional === "string" ? settings.additional.charCodeAt(0) : settings.additional;
    const result = [];
    let index2 = 0;
    let lines = -1;
    let queue = "";
    let point3;
    let indent3;
    if (settings.position) {
      if ("start" in settings.position || "indent" in settings.position) {
        indent3 = settings.position.indent;
        point3 = settings.position.start;
      } else {
        point3 = settings.position;
      }
    }
    let line = (point3 ? point3.line : 0) || 1;
    let column = (point3 ? point3.column : 0) || 1;
    let previous3 = now();
    let character;
    index2--;
    while (++index2 <= value.length) {
      if (character === 10) {
        column = (indent3 ? indent3[lines] : 0) || 1;
      }
      character = value.charCodeAt(index2);
      if (character === 38) {
        const following = value.charCodeAt(index2 + 1);
        if (following === 9 || following === 10 || following === 12 || following === 32 || following === 38 || following === 60 || Number.isNaN(following) || additional && following === additional) {
          queue += String.fromCharCode(character);
          column++;
          continue;
        }
        const start2 = index2 + 1;
        let begin = start2;
        let end = start2;
        let type;
        if (following === 35) {
          end = ++begin;
          const following2 = value.charCodeAt(end);
          if (following2 === 88 || following2 === 120) {
            type = "hexadecimal";
            end = ++begin;
          } else {
            type = "decimal";
          }
        } else {
          type = "named";
        }
        let characterReferenceCharacters = "";
        let characterReference2 = "";
        let characters = "";
        const test = type === "named" ? isAlphanumerical : type === "decimal" ? isDecimal : isHexadecimal;
        end--;
        while (++end <= value.length) {
          const following2 = value.charCodeAt(end);
          if (!test(following2)) {
            break;
          }
          characters += String.fromCharCode(following2);
          if (type === "named" && characterEntitiesLegacy.includes(characters)) {
            characterReferenceCharacters = characters;
            characterReference2 = decodeNamedCharacterReference(characters);
          }
        }
        let terminated = value.charCodeAt(end) === 59;
        if (terminated) {
          end++;
          const namedReference = type === "named" ? decodeNamedCharacterReference(characters) : false;
          if (namedReference) {
            characterReferenceCharacters = characters;
            characterReference2 = namedReference;
          }
        }
        let diff = 1 + end - start2;
        let reference = "";
        if (!terminated && settings.nonTerminated === false) {
        } else if (!characters) {
          if (type !== "named") {
            warning(4, diff);
          }
        } else if (type === "named") {
          if (terminated && !characterReference2) {
            warning(5, 1);
          } else {
            if (characterReferenceCharacters !== characters) {
              end = begin + characterReferenceCharacters.length;
              diff = 1 + end - begin;
              terminated = false;
            }
            if (!terminated) {
              const reason = characterReferenceCharacters ? 1 : 3;
              if (settings.attribute) {
                const following2 = value.charCodeAt(end);
                if (following2 === 61) {
                  warning(reason, diff);
                  characterReference2 = "";
                } else if (isAlphanumerical(following2)) {
                  characterReference2 = "";
                } else {
                  warning(reason, diff);
                }
              } else {
                warning(reason, diff);
              }
            }
          }
          reference = characterReference2;
        } else {
          if (!terminated) {
            warning(2, diff);
          }
          let referenceCode = Number.parseInt(
            characters,
            type === "hexadecimal" ? 16 : 10
          );
          if (prohibited(referenceCode)) {
            warning(7, diff);
            reference = String.fromCharCode(
              65533
              /* `�` */
            );
          } else if (referenceCode in characterReferenceInvalid) {
            warning(6, diff);
            reference = characterReferenceInvalid[referenceCode];
          } else {
            let output = "";
            if (disallowed(referenceCode)) {
              warning(6, diff);
            }
            if (referenceCode > 65535) {
              referenceCode -= 65536;
              output += String.fromCharCode(
                referenceCode >>> (10 & 1023) | 55296
              );
              referenceCode = 56320 | referenceCode & 1023;
            }
            reference = output + String.fromCharCode(referenceCode);
          }
        }
        if (reference) {
          flush();
          previous3 = now();
          index2 = end - 1;
          column += end - start2 + 1;
          result.push(reference);
          const next = now();
          next.offset++;
          if (settings.reference) {
            settings.reference.call(
              settings.referenceContext || void 0,
              reference,
              { start: previous3, end: next },
              value.slice(start2 - 1, end)
            );
          }
          previous3 = next;
        } else {
          characters = value.slice(start2 - 1, end);
          queue += characters;
          column += characters.length;
          index2 = end - 1;
        }
      } else {
        if (character === 10) {
          line++;
          lines++;
          column = 0;
        }
        if (Number.isNaN(character)) {
          flush();
        } else {
          queue += String.fromCharCode(character);
          column++;
        }
      }
    }
    return result.join("");
    function now() {
      return {
        line,
        column,
        offset: index2 + ((point3 ? point3.offset : 0) || 0)
      };
    }
    function warning(code3, offset2) {
      let position2;
      if (settings.warning) {
        position2 = now();
        position2.column += offset2;
        position2.offset += offset2;
        settings.warning.call(
          settings.warningContext || void 0,
          messages[code3],
          position2,
          code3
        );
      }
    }
    function flush() {
      if (queue) {
        result.push(queue);
        if (settings.text) {
          settings.text.call(settings.textContext || void 0, queue, {
            start: previous3,
            end: now()
          });
        }
        queue = "";
      }
    }
  }
  function prohibited(code3) {
    return code3 >= 55296 && code3 <= 57343 || code3 > 1114111;
  }
  function disallowed(code3) {
    return code3 >= 1 && code3 <= 8 || code3 === 11 || code3 >= 13 && code3 <= 31 || code3 >= 127 && code3 <= 159 || code3 >= 64976 && code3 <= 65007 || (code3 & 65535) === 65535 || (code3 & 65535) === 65534;
  }

  // ../../../ui/node_modules/.pnpm/stringify-entities@4.0.4/node_modules/stringify-entities/lib/core.js
  var defaultSubsetRegex = /["&'<>`]/g;
  var surrogatePairsRegex = /[\uD800-\uDBFF][\uDC00-\uDFFF]/g;
  var controlCharactersRegex = (
    // eslint-disable-next-line no-control-regex, unicorn/no-hex-escape
    /[\x01-\t\v\f\x0E-\x1F\x7F\x81\x8D\x8F\x90\x9D\xA0-\uFFFF]/g
  );
  var regexEscapeRegex = /[|\\{}()[\]^$+*?.]/g;
  var subsetToRegexCache = /* @__PURE__ */ new WeakMap();
  function core(value, options) {
    value = value.replace(
      options.subset ? charactersToExpressionCached(options.subset) : defaultSubsetRegex,
      basic
    );
    if (options.subset || options.escapeOnly) {
      return value;
    }
    return value.replace(surrogatePairsRegex, surrogate).replace(controlCharactersRegex, basic);
    function surrogate(pair, index2, all2) {
      return options.format(
        (pair.charCodeAt(0) - 55296) * 1024 + pair.charCodeAt(1) - 56320 + 65536,
        all2.charCodeAt(index2 + 2),
        options
      );
    }
    function basic(character, index2, all2) {
      return options.format(
        character.charCodeAt(0),
        all2.charCodeAt(index2 + 1),
        options
      );
    }
  }
  function charactersToExpressionCached(subset) {
    let cached = subsetToRegexCache.get(subset);
    if (!cached) {
      cached = charactersToExpression(subset);
      subsetToRegexCache.set(subset, cached);
    }
    return cached;
  }
  function charactersToExpression(subset) {
    const groups = [];
    let index2 = -1;
    while (++index2 < subset.length) {
      groups.push(subset[index2].replace(regexEscapeRegex, "\\$&"));
    }
    return new RegExp("(?:" + groups.join("|") + ")", "g");
  }

  // ../../../ui/node_modules/.pnpm/stringify-entities@4.0.4/node_modules/stringify-entities/lib/util/format-basic.js
  function formatBasic(code3) {
    return "&#x" + code3.toString(16).toUpperCase() + ";";
  }

  // ../../../ui/node_modules/.pnpm/stringify-entities@4.0.4/node_modules/stringify-entities/lib/index.js
  function stringifyEntitiesLight(value, options) {
    return core(value, Object.assign({ format: formatBasic }, options));
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-mdx-jsx@3.2.0/node_modules/mdast-util-mdx-jsx/lib/index.js
  var indent = "  ";
  function mdxJsxFromMarkdown() {
    return {
      canContainEols: ["mdxJsxTextElement"],
      enter: {
        mdxJsxFlowTag: enterMdxJsxTag,
        mdxJsxFlowTagClosingMarker: enterMdxJsxTagClosingMarker,
        mdxJsxFlowTagAttribute: enterMdxJsxTagAttribute,
        mdxJsxFlowTagExpressionAttribute: enterMdxJsxTagExpressionAttribute,
        mdxJsxFlowTagAttributeValueLiteral: buffer,
        mdxJsxFlowTagAttributeValueExpression: buffer,
        mdxJsxFlowTagSelfClosingMarker: enterMdxJsxTagSelfClosingMarker,
        mdxJsxTextTag: enterMdxJsxTag,
        mdxJsxTextTagClosingMarker: enterMdxJsxTagClosingMarker,
        mdxJsxTextTagAttribute: enterMdxJsxTagAttribute,
        mdxJsxTextTagExpressionAttribute: enterMdxJsxTagExpressionAttribute,
        mdxJsxTextTagAttributeValueLiteral: buffer,
        mdxJsxTextTagAttributeValueExpression: buffer,
        mdxJsxTextTagSelfClosingMarker: enterMdxJsxTagSelfClosingMarker
      },
      exit: {
        mdxJsxFlowTagClosingMarker: exitMdxJsxTagClosingMarker,
        mdxJsxFlowTagNamePrimary: exitMdxJsxTagNamePrimary,
        mdxJsxFlowTagNameMember: exitMdxJsxTagNameMember,
        mdxJsxFlowTagNameLocal: exitMdxJsxTagNameLocal,
        mdxJsxFlowTagExpressionAttribute: exitMdxJsxTagExpressionAttribute,
        mdxJsxFlowTagExpressionAttributeValue: data2,
        mdxJsxFlowTagAttributeNamePrimary: exitMdxJsxTagAttributeNamePrimary,
        mdxJsxFlowTagAttributeNameLocal: exitMdxJsxTagAttributeNameLocal,
        mdxJsxFlowTagAttributeValueLiteral: exitMdxJsxTagAttributeValueLiteral,
        mdxJsxFlowTagAttributeValueLiteralValue: data2,
        mdxJsxFlowTagAttributeValueExpression: exitMdxJsxTagAttributeValueExpression,
        mdxJsxFlowTagAttributeValueExpressionValue: data2,
        mdxJsxFlowTagSelfClosingMarker: exitMdxJsxTagSelfClosingMarker,
        mdxJsxFlowTag: exitMdxJsxTag,
        mdxJsxTextTagClosingMarker: exitMdxJsxTagClosingMarker,
        mdxJsxTextTagNamePrimary: exitMdxJsxTagNamePrimary,
        mdxJsxTextTagNameMember: exitMdxJsxTagNameMember,
        mdxJsxTextTagNameLocal: exitMdxJsxTagNameLocal,
        mdxJsxTextTagExpressionAttribute: exitMdxJsxTagExpressionAttribute,
        mdxJsxTextTagExpressionAttributeValue: data2,
        mdxJsxTextTagAttributeNamePrimary: exitMdxJsxTagAttributeNamePrimary,
        mdxJsxTextTagAttributeNameLocal: exitMdxJsxTagAttributeNameLocal,
        mdxJsxTextTagAttributeValueLiteral: exitMdxJsxTagAttributeValueLiteral,
        mdxJsxTextTagAttributeValueLiteralValue: data2,
        mdxJsxTextTagAttributeValueExpression: exitMdxJsxTagAttributeValueExpression,
        mdxJsxTextTagAttributeValueExpressionValue: data2,
        mdxJsxTextTagSelfClosingMarker: exitMdxJsxTagSelfClosingMarker,
        mdxJsxTextTag: exitMdxJsxTag
      }
    };
    function buffer() {
      this.buffer();
    }
    function point3(d) {
      return { line: d.line, column: d.column, offset: d.offset };
    }
    function data2(token) {
      this.config.enter.data.call(this, token);
      this.config.exit.data.call(this, token);
    }
    function enterMdxJsxTag(token) {
      const tag = {
        name: void 0,
        attributes: [],
        close: false,
        selfClosing: false,
        start: token.start,
        end: token.end
      };
      if (!this.data.mdxJsxTagStack) this.data.mdxJsxTagStack = [];
      this.data.mdxJsxTag = tag;
      this.buffer();
    }
    function enterMdxJsxTagClosingMarker(token) {
      const stack = this.data.mdxJsxTagStack;
      ok(stack, "expected `mdxJsxTagStack`");
      if (stack.length === 0) {
        throw new VFileMessage(
          "Unexpected closing slash `/` in tag, expected an open tag first",
          { start: token.start, end: token.end },
          "mdast-util-mdx-jsx:unexpected-closing-slash"
        );
      }
    }
    function enterMdxJsxTagAnyAttribute(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      if (tag.close) {
        throw new VFileMessage(
          "Unexpected attribute in closing tag, expected the end of the tag",
          { start: token.start, end: token.end },
          "mdast-util-mdx-jsx:unexpected-attribute"
        );
      }
    }
    function enterMdxJsxTagSelfClosingMarker(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      if (tag.close) {
        throw new VFileMessage(
          "Unexpected self-closing slash `/` in closing tag, expected the end of the tag",
          { start: token.start, end: token.end },
          "mdast-util-mdx-jsx:unexpected-self-closing-slash"
        );
      }
    }
    function exitMdxJsxTagClosingMarker() {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      tag.close = true;
    }
    function exitMdxJsxTagNamePrimary(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      tag.name = this.sliceSerialize(token);
    }
    function exitMdxJsxTagNameMember(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      tag.name += "." + this.sliceSerialize(token);
    }
    function exitMdxJsxTagNameLocal(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      tag.name += ":" + this.sliceSerialize(token);
    }
    function enterMdxJsxTagAttribute(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      enterMdxJsxTagAnyAttribute.call(this, token);
      tag.attributes.push({
        type: "mdxJsxAttribute",
        name: "",
        value: null,
        position: {
          start: point3(token.start),
          // @ts-expect-error: `end` will be patched later.
          end: void 0
        }
      });
    }
    function enterMdxJsxTagExpressionAttribute(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      enterMdxJsxTagAnyAttribute.call(this, token);
      tag.attributes.push({
        type: "mdxJsxExpressionAttribute",
        value: "",
        position: {
          start: point3(token.start),
          // @ts-expect-error: `end` will be patched later.
          end: void 0
        }
      });
      this.buffer();
    }
    function exitMdxJsxTagExpressionAttribute(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const tail = tag.attributes[tag.attributes.length - 1];
      ok(tail.type === "mdxJsxExpressionAttribute");
      const estree = token.estree;
      tail.value = this.resume();
      ok(tail.position !== void 0);
      tail.position.end = point3(token.end);
      if (estree) {
        tail.data = { estree };
      }
    }
    function exitMdxJsxTagAttributeNamePrimary(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const node2 = tag.attributes[tag.attributes.length - 1];
      ok(node2.type === "mdxJsxAttribute");
      node2.name = this.sliceSerialize(token);
      ok(node2.position !== void 0);
      node2.position.end = point3(token.end);
    }
    function exitMdxJsxTagAttributeNameLocal(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const node2 = tag.attributes[tag.attributes.length - 1];
      ok(node2.type === "mdxJsxAttribute");
      node2.name += ":" + this.sliceSerialize(token);
      ok(node2.position !== void 0);
      node2.position.end = point3(token.end);
    }
    function exitMdxJsxTagAttributeValueLiteral(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const node2 = tag.attributes[tag.attributes.length - 1];
      node2.value = parseEntities(this.resume(), { nonTerminated: false });
      ok(node2.position !== void 0);
      node2.position.end = point3(token.end);
    }
    function exitMdxJsxTagAttributeValueExpression(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const tail = tag.attributes[tag.attributes.length - 1];
      ok(tail.type === "mdxJsxAttribute");
      const node2 = { type: "mdxJsxAttributeValueExpression", value: this.resume() };
      const estree = token.estree;
      if (estree) {
        node2.data = { estree };
      }
      tail.value = node2;
      ok(tail.position !== void 0);
      tail.position.end = point3(token.end);
    }
    function exitMdxJsxTagSelfClosingMarker() {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      tag.selfClosing = true;
    }
    function exitMdxJsxTag(token) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      const stack = this.data.mdxJsxTagStack;
      ok(stack, "expected `mdxJsxTagStack`");
      const tail = stack[stack.length - 1];
      if (tag.close && tail.name !== tag.name) {
        throw new VFileMessage(
          "Unexpected closing tag `" + serializeAbbreviatedTag(tag) + "`, expected corresponding closing tag for `" + serializeAbbreviatedTag(tail) + "` (" + stringifyPosition(tail) + ")",
          { start: token.start, end: token.end },
          "mdast-util-mdx-jsx:end-tag-mismatch"
        );
      }
      this.resume();
      if (tag.close) {
        stack.pop();
      } else {
        this.enter(
          {
            type: token.type === "mdxJsxTextTag" ? "mdxJsxTextElement" : "mdxJsxFlowElement",
            name: tag.name || null,
            attributes: tag.attributes,
            children: []
          },
          token,
          onErrorRightIsTag
        );
      }
      if (tag.selfClosing || tag.close) {
        this.exit(token, onErrorLeftIsTag);
      } else {
        stack.push(tag);
      }
    }
    function onErrorRightIsTag(closing, open) {
      const stack = this.data.mdxJsxTagStack;
      ok(stack, "expected `mdxJsxTagStack`");
      const tag = stack[stack.length - 1];
      ok(tag, "expected `mdxJsxTag`");
      const place = closing ? " before the end of `" + closing.type + "`" : "";
      const position2 = closing ? { start: closing.start, end: closing.end } : void 0;
      throw new VFileMessage(
        "Expected a closing tag for `" + serializeAbbreviatedTag(tag) + "` (" + stringifyPosition({ start: open.start, end: open.end }) + ")" + place,
        position2,
        "mdast-util-mdx-jsx:end-tag-mismatch"
      );
    }
    function onErrorLeftIsTag(a, b) {
      const tag = this.data.mdxJsxTag;
      ok(tag, "expected `mdxJsxTag`");
      throw new VFileMessage(
        "Expected the closing tag `" + serializeAbbreviatedTag(tag) + "` either after the end of `" + b.type + "` (" + stringifyPosition(b.end) + ") or another opening tag after the start of `" + b.type + "` (" + stringifyPosition(b.start) + ")",
        { start: a.start, end: a.end },
        "mdast-util-mdx-jsx:end-tag-mismatch"
      );
    }
    function serializeAbbreviatedTag(tag) {
      return "<" + (tag.close ? "/" : "") + (tag.name || "") + ">";
    }
  }
  function mdxJsxToMarkdown(options) {
    const options_ = options || {};
    const quote = options_.quote || '"';
    const quoteSmart = options_.quoteSmart || false;
    const tightSelfClosing = options_.tightSelfClosing || false;
    const printWidth = options_.printWidth || Number.POSITIVE_INFINITY;
    const alternative = quote === '"' ? "'" : '"';
    if (quote !== '"' && quote !== "'") {
      throw new Error(
        "Cannot serialize attribute values with `" + quote + "` for `options.quote`, expected `\"`, or `'`"
      );
    }
    mdxElement.peek = peekElement;
    return {
      handlers: {
        mdxJsxFlowElement: mdxElement,
        mdxJsxTextElement: mdxElement
      },
      unsafe: [
        { character: "<", inConstruct: ["phrasing"] },
        { atBreak: true, character: "<" }
      ],
      // Always generate fenced code (never indented code).
      fences: true,
      // Always generate links with resources (never autolinks).
      resourceLink: true
    };
    function mdxElement(node2, _, state, info) {
      const flow3 = node2.type === "mdxJsxFlowElement";
      const selfClosing = node2.name ? !node2.children || node2.children.length === 0 : false;
      const depth = inferDepth(state);
      const currentIndent = createIndent(depth);
      const trackerOneLine = state.createTracker(info);
      const trackerMultiLine = state.createTracker(info);
      const serializedAttributes = [];
      const prefix = (flow3 ? currentIndent : "") + "<" + (node2.name || "");
      const exit3 = state.enter(node2.type);
      trackerOneLine.move(prefix);
      trackerMultiLine.move(prefix);
      if (node2.attributes && node2.attributes.length > 0) {
        if (!node2.name) {
          throw new Error("Cannot serialize fragment w/ attributes");
        }
        let index2 = -1;
        while (++index2 < node2.attributes.length) {
          const attribute = node2.attributes[index2];
          let result;
          if (attribute.type === "mdxJsxExpressionAttribute") {
            result = "{" + (attribute.value || "") + "}";
          } else {
            if (!attribute.name) {
              throw new Error("Cannot serialize attribute w/o name");
            }
            const value2 = attribute.value;
            const left = attribute.name;
            let right = "";
            if (value2 === null || value2 === void 0) {
            } else if (typeof value2 === "object") {
              right = "{" + (value2.value || "") + "}";
            } else {
              const appliedQuote = quoteSmart && ccount(value2, quote) > ccount(value2, alternative) ? alternative : quote;
              right = appliedQuote + stringifyEntitiesLight(value2, { subset: [appliedQuote] }) + appliedQuote;
            }
            result = left + (right ? "=" : "") + right;
          }
          serializedAttributes.push(result);
        }
      }
      let attributesOnTheirOwnLine = false;
      const attributesOnOneLine = serializedAttributes.join(" ");
      if (
        // Block:
        flow3 && // Including a line ending (expressions).
        (/\r?\n|\r/.test(attributesOnOneLine) || // Current position (including `<tag`).
        trackerOneLine.current().now.column + // -1 because columns, +1 for ` ` before attributes.
        // Attributes joined by spaces.
        attributesOnOneLine.length + // ` />`.
        (selfClosing ? tightSelfClosing ? 2 : 3 : 1) > printWidth)
      ) {
        attributesOnTheirOwnLine = true;
      }
      let tracker = trackerOneLine;
      let value = prefix;
      if (attributesOnTheirOwnLine) {
        tracker = trackerMultiLine;
        let index2 = -1;
        while (++index2 < serializedAttributes.length) {
          serializedAttributes[index2] = currentIndent + indent + serializedAttributes[index2];
        }
        value += tracker.move(
          "\n" + serializedAttributes.join("\n") + "\n" + currentIndent
        );
      } else if (attributesOnOneLine) {
        value += tracker.move(" " + attributesOnOneLine);
      }
      if (selfClosing) {
        value += tracker.move(
          (tightSelfClosing || attributesOnTheirOwnLine ? "" : " ") + "/"
        );
      }
      value += tracker.move(">");
      if (node2.children && node2.children.length > 0) {
        if (node2.type === "mdxJsxTextElement") {
          value += tracker.move(
            state.containerPhrasing(node2, {
              ...tracker.current(),
              before: ">",
              after: "<"
            })
          );
        } else {
          tracker.shift(2);
          value += tracker.move("\n");
          value += tracker.move(containerFlow(node2, state, tracker.current()));
          value += tracker.move("\n");
        }
      }
      if (!selfClosing) {
        value += tracker.move(
          (flow3 ? currentIndent : "") + "</" + (node2.name || "") + ">"
        );
      }
      exit3();
      return value;
    }
  }
  function containerFlow(parent, state, info) {
    const indexStack = state.indexStack;
    const children = parent.children;
    const tracker = state.createTracker(info);
    const currentIndent = createIndent(inferDepth(state));
    const results = [];
    let index2 = -1;
    indexStack.push(-1);
    while (++index2 < children.length) {
      const child = children[index2];
      indexStack[indexStack.length - 1] = index2;
      const childInfo = { before: "\n", after: "\n", ...tracker.current() };
      const result = state.handle(child, parent, state, childInfo);
      const serializedChild = child.type === "mdxJsxFlowElement" ? result : state.indentLines(result, function(line, _, blank) {
        return (blank ? "" : currentIndent) + line;
      });
      results.push(tracker.move(serializedChild));
      if (child.type !== "list") {
        state.bulletLastUsed = void 0;
      }
      if (index2 < children.length - 1) {
        results.push(tracker.move("\n\n"));
      }
    }
    indexStack.pop();
    return results.join("");
  }
  function inferDepth(state) {
    let depth = 0;
    let index2 = state.stack.length;
    while (--index2 > -1) {
      const name2 = state.stack[index2];
      if (name2 === "blockquote" || name2 === "listItem") break;
      if (name2 === "mdxJsxFlowElement") depth++;
    }
    return depth;
  }
  function createIndent(depth) {
    return indent.repeat(depth);
  }
  function peekElement() {
    return "<";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-mdxjs-esm@2.0.1/node_modules/mdast-util-mdxjs-esm/lib/index.js
  function mdxjsEsmFromMarkdown() {
    return {
      enter: { mdxjsEsm: enterMdxjsEsm },
      exit: { mdxjsEsm: exitMdxjsEsm, mdxjsEsmData: exitMdxjsEsmData }
    };
  }
  function mdxjsEsmToMarkdown() {
    return { handlers: { mdxjsEsm: handleMdxjsEsm } };
  }
  function enterMdxjsEsm(token) {
    this.enter({ type: "mdxjsEsm", value: "" }, token);
    this.buffer();
  }
  function exitMdxjsEsm(token) {
    const value = this.resume();
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "mdxjsEsm");
    this.exit(token);
    const estree = token.estree;
    node2.value = value;
    if (estree) {
      node2.data = { estree };
    }
  }
  function exitMdxjsEsmData(token) {
    this.config.enter.data.call(this, token);
    this.config.exit.data.call(this, token);
  }
  function handleMdxjsEsm(node2) {
    return node2.value || "";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-mdx@3.0.0/node_modules/mdast-util-mdx/lib/index.js
  function mdxFromMarkdown() {
    return [
      mdxExpressionFromMarkdown(),
      mdxJsxFromMarkdown(),
      mdxjsEsmFromMarkdown()
    ];
  }
  function mdxToMarkdown(options) {
    return {
      extensions: [
        mdxExpressionToMarkdown(),
        mdxJsxToMarkdown(options),
        mdxjsEsmToMarkdown()
      ]
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdxjs@3.0.0/node_modules/micromark-extension-mdxjs/index.js
  init_acorn();
  var import_acorn_jsx = __toESM(require_acorn_jsx(), 1);

  // ../../../ui/node_modules/.pnpm/micromark-util-character@2.1.1/node_modules/micromark-util-character/index.js
  var asciiAlpha = regexCheck(/[A-Za-z]/);
  var asciiAlphanumeric = regexCheck(/[\dA-Za-z]/);
  var asciiAtext = regexCheck(/[#-'*+\--9=?A-Z^-~]/);
  function asciiControl(code3) {
    return (
      // Special whitespace codes (which have negative values), C0 and Control
      // character DEL
      code3 !== null && (code3 < 32 || code3 === 127)
    );
  }
  var asciiDigit = regexCheck(/\d/);
  var asciiHexDigit = regexCheck(/[\dA-Fa-f]/);
  var asciiPunctuation = regexCheck(/[!-/:-@[-`{-~]/);
  function markdownLineEnding(code3) {
    return code3 !== null && code3 < -2;
  }
  function markdownLineEndingOrSpace(code3) {
    return code3 !== null && (code3 < 0 || code3 === 32);
  }
  function markdownSpace(code3) {
    return code3 === -2 || code3 === -1 || code3 === 32;
  }
  var unicodePunctuation = regexCheck(/[\u0021-\u002f\u003a-\u0040\u005b-\u0060\u007b-\u007e\u00a1-\u00a9\u00ab-\u00ac\u00ae-\u00b1\u00b4\u00b6-\u00b8\u00bb\u00bf\u00d7\u00f7\u02c2-\u02c5\u02d2-\u02df\u02e5-\u02eb\u02ed\u02ef-\u02ff\u0375\u037e\u0384-\u0385\u0387\u03f6\u0482\u055a-\u055f\u0589-\u058a\u058d-\u058f\u05be\u05c0\u05c3\u05c6\u05f3-\u05f4\u0606-\u060f\u061b\u061d-\u061f\u066a-\u066d\u06d4\u06de\u06e9\u06fd-\u06fe\u0700-\u070d\u07f6-\u07f9\u07fe-\u07ff\u0830-\u083e\u085e\u0888\u0964-\u0965\u0970\u09f2-\u09f3\u09fa-\u09fb\u09fd\u0a76\u0af0-\u0af1\u0b70\u0bf3-\u0bfa\u0c77\u0c7f\u0c84\u0d4f\u0d79\u0df4\u0e3f\u0e4f\u0e5a-\u0e5b\u0f01-\u0f17\u0f1a-\u0f1f\u0f34\u0f36\u0f38\u0f3a-\u0f3d\u0f85\u0fbe-\u0fc5\u0fc7-\u0fcc\u0fce-\u0fda\u104a-\u104f\u109e-\u109f\u10fb\u1360-\u1368\u1390-\u1399\u1400\u166d-\u166e\u169b-\u169c\u16eb-\u16ed\u1735-\u1736\u17d4-\u17d6\u17d8-\u17db\u1800-\u180a\u1940\u1944-\u1945\u19de-\u19ff\u1a1e-\u1a1f\u1aa0-\u1aa6\u1aa8-\u1aad\u1b4e-\u1b4f\u1b5a-\u1b6a\u1b74-\u1b7f\u1bfc-\u1bff\u1c3b-\u1c3f\u1c7e-\u1c7f\u1cc0-\u1cc7\u1cd3\u1fbd\u1fbf-\u1fc1\u1fcd-\u1fcf\u1fdd-\u1fdf\u1fed-\u1fef\u1ffd-\u1ffe\u2010-\u2027\u2030-\u205e\u207a-\u207e\u208a-\u208e\u20a0-\u20c1\u2100-\u2101\u2103-\u2106\u2108-\u2109\u2114\u2116-\u2118\u211e-\u2123\u2125\u2127\u2129\u212e\u213a-\u213b\u2140-\u2144\u214a-\u214d\u214f\u218a-\u218b\u2190-\u2429\u2440-\u244a\u249c-\u24e9\u2500-\u2775\u2794-\u2b73\u2b76-\u2bff\u2ce5-\u2cea\u2cf9-\u2cfc\u2cfe-\u2cff\u2d70\u2e00-\u2e2e\u2e30-\u2e5d\u2e80-\u2e99\u2e9b-\u2ef3\u2f00-\u2fd5\u2ff0-\u2fff\u3001-\u3004\u3008-\u3020\u3030\u3036-\u3037\u303d-\u303f\u309b-\u309c\u30a0\u30fb\u3190-\u3191\u3196-\u319f\u31c0-\u31e5\u31ef\u3200-\u321e\u322a-\u3247\u3250\u3260-\u327f\u328a-\u32b0\u32c0-\u33ff\u4dc0-\u4dff\ua490-\ua4c6\ua4fe-\ua4ff\ua60d-\ua60f\ua673\ua67e\ua6f2-\ua6f7\ua700-\ua716\ua720-\ua721\ua789-\ua78a\ua828-\ua82b\ua836-\ua839\ua874-\ua877\ua8ce-\ua8cf\ua8f8-\ua8fa\ua8fc\ua92e-\ua92f\ua95f\ua9c1-\ua9cd\ua9de-\ua9df\uaa5c-\uaa5f\uaa77-\uaa79\uaade-\uaadf\uaaf0-\uaaf1\uab5b\uab6a-\uab6b\uabeb\ufb29\ufbb2-\ufbd2\ufd3e-\ufd4f\ufd90-\ufd91\ufdc8-\ufdcf\ufdfc-\ufdff\ufe10-\ufe19\ufe30-\ufe52\ufe54-\ufe66\ufe68-\ufe6b\uff01-\uff0f\uff1a-\uff20\uff3b-\uff40\uff5b-\uff65\uffe0-\uffe6\uffe8-\uffee\ufffc-\ufffd]/);
  var unicodeWhitespace = regexCheck(/\s/);
  function regexCheck(regex2) {
    return check;
    function check(code3) {
      return code3 !== null && code3 > -1 && regex2.test(String.fromCharCode(code3));
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-space@2.0.1/node_modules/micromark-factory-space/index.js
  function factorySpace(effects, ok3, type, max) {
    const limit = max ? max - 1 : Number.POSITIVE_INFINITY;
    let size = 0;
    return start2;
    function start2(code3) {
      if (markdownSpace(code3)) {
        effects.enter(type);
        return prefix(code3);
      }
      return ok3(code3);
    }
    function prefix(code3) {
      if (markdownSpace(code3) && size++ < limit) {
        effects.consume(code3);
        return prefix;
      }
      effects.exit(type);
      return ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/estree-util-visit@2.0.0/node_modules/estree-util-visit/lib/color.default.js
  function color(d) {
    return d;
  }

  // ../../../ui/node_modules/.pnpm/estree-util-visit@2.0.0/node_modules/estree-util-visit/lib/index.js
  var own3 = {}.hasOwnProperty;
  var CONTINUE = Symbol("continue");
  var EXIT = Symbol("exit");
  var SKIP = Symbol("skip");
  function visit(tree, visitor) {
    let enter;
    let leave;
    if (typeof visitor === "function") {
      enter = visitor;
    } else if (visitor && typeof visitor === "object") {
      if (visitor.enter) enter = visitor.enter;
      if (visitor.leave) leave = visitor.leave;
    }
    build(tree, void 0, void 0, [])();
    function build(node2, key, index2, parents) {
      if (nodelike(node2)) {
        visit3.displayName = "node (" + color(node2.type) + ")";
      }
      return visit3;
      function visit3() {
        const result = enter ? toResult(enter(node2, key, index2, parents)) : [];
        if (result[0] === EXIT) {
          return result;
        }
        if (result[0] !== SKIP) {
          let cKey;
          for (cKey in node2) {
            if (own3.call(node2, cKey) && node2[cKey] && typeof node2[cKey] === "object" && // @ts-expect-error: custom esast extension.
            cKey !== "data" && // @ts-expect-error: custom esast extension.
            cKey !== "position") {
              const grandparents = parents.concat(node2);
              const value = node2[cKey];
              if (Array.isArray(value)) {
                const nodes = (
                  /** @type {Array<unknown>} */
                  value
                );
                let cIndex = 0;
                while (cIndex > -1 && cIndex < nodes.length) {
                  const subvalue = nodes[cIndex];
                  if (nodelike(subvalue)) {
                    const subresult = build(
                      subvalue,
                      cKey,
                      cIndex,
                      grandparents
                    )();
                    if (subresult[0] === EXIT) return subresult;
                    cIndex = typeof subresult[1] === "number" ? subresult[1] : cIndex + 1;
                  } else {
                    cIndex++;
                  }
                }
              } else if (nodelike(value)) {
                const subresult = build(value, cKey, void 0, grandparents)();
                if (subresult[0] === EXIT) return subresult;
              }
            }
          }
        }
        return leave ? toResult(leave(node2, key, index2, parents)) : result;
      }
    }
  }
  function toResult(value) {
    if (Array.isArray(value)) {
      return value;
    }
    if (typeof value === "number") {
      return [CONTINUE, value];
    }
    return [value];
  }
  function nodelike(value) {
    return Boolean(
      value && typeof value === "object" && "type" in value && typeof value.type === "string" && value.type.length > 0
    );
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-events-to-acorn@2.0.3/node_modules/micromark-util-events-to-acorn/lib/index.js
  function eventsToAcorn(events, options) {
    const prefix = options.prefix || "";
    const suffix = options.suffix || "";
    const acornOptions = Object.assign({}, options.acornOptions);
    const comments = [];
    const tokens = [];
    const onComment = acornOptions.onComment;
    const onToken = acornOptions.onToken;
    let swallow = false;
    let estree;
    let exception;
    const acornConfig = Object.assign({}, acornOptions, {
      onComment: comments,
      preserveParens: true
    });
    if (onToken) {
      acornConfig.onToken = tokens;
    }
    const collection = collect(events, options.tokenTypes);
    const source = collection.value;
    const value = prefix + source + suffix;
    const isEmptyExpression = options.expression && empty2(source);
    if (isEmptyExpression && !options.allowEmpty) {
      throw new VFileMessage("Unexpected empty expression", {
        place: parseOffsetToUnistPoint(0),
        ruleId: "unexpected-empty-expression",
        source: "micromark-extension-mdx-expression"
      });
    }
    try {
      estree = options.expression && !isEmptyExpression ? options.acorn.parseExpressionAt(value, 0, acornConfig) : options.acorn.parse(value, acornConfig);
    } catch (error_) {
      const error = (
        /** @type {AcornError} */
        error_
      );
      const point3 = parseOffsetToUnistPoint(error.pos);
      error.message = String(error.message).replace(/ \(\d+:\d+\)$/, "");
      error.pos = point3.offset;
      error.loc = {
        line: point3.line,
        column: point3.column - 1
      };
      exception = error;
      swallow = error.raisedAt >= prefix.length + source.length || // Broken comments are raised at their start, not their end.
      error.message === "Unterminated comment";
    }
    if (estree && options.expression && !isEmptyExpression) {
      if (empty2(value.slice(estree.end, value.length - suffix.length))) {
        estree = {
          type: "Program",
          start: 0,
          end: prefix.length + source.length,
          // @ts-expect-error: It’s good.
          body: [{
            type: "ExpressionStatement",
            expression: estree,
            start: 0,
            end: prefix.length + source.length
          }],
          sourceType: "module",
          comments: []
        };
      } else {
        const point3 = parseOffsetToUnistPoint(estree.end);
        const error = (
          /** @type {AcornError} */
          new Error("Unexpected content after expression")
        );
        error.pos = point3.offset;
        error.loc = {
          line: point3.line,
          column: point3.column - 1
        };
        exception = error;
        estree = void 0;
      }
    }
    if (estree) {
      estree.comments = comments;
      visit(estree, function(esnode, field, index2, parents) {
        let context = (
          /** @type {AcornNode | Array<AcornNode>} */
          parents[parents.length - 1]
        );
        let property2 = field;
        if (esnode.type === "ParenthesizedExpression" && context && property2) {
          if (typeof index2 === "number") {
            context = context[property2];
            property2 = index2;
          }
          context[property2] = esnode.expression;
        }
        fixPosition(esnode);
      });
      if (Array.isArray(onComment)) {
        onComment.push(...comments);
      } else if (typeof onComment === "function") {
        for (const comment of comments) {
          onComment(comment.type === "Block", comment.value, comment.start, comment.end, comment.loc.start, comment.loc.end);
        }
      }
      for (const token of tokens) {
        if (token.end <= prefix.length || token.start - prefix.length >= source.length) {
          continue;
        }
        fixPosition(token);
        if (Array.isArray(onToken)) {
          onToken.push(token);
        } else {
          onToken(token);
        }
      }
    }
    return {
      estree,
      error: exception,
      swallow
    };
    function fixPosition(nodeOrToken) {
      const pointStart = parseOffsetToUnistPoint(nodeOrToken.start);
      const pointEnd = parseOffsetToUnistPoint(nodeOrToken.end);
      nodeOrToken.start = pointStart.offset;
      nodeOrToken.end = pointEnd.offset;
      nodeOrToken.loc = {
        start: {
          line: pointStart.line,
          column: pointStart.column - 1,
          // @ts-expect-error: not allowed by acorn types.
          offset: pointStart.offset
        },
        end: {
          line: pointEnd.line,
          column: pointEnd.column - 1,
          // @ts-expect-error: not allowed by acorn types.
          offset: pointEnd.offset
        }
      };
      nodeOrToken.range = [nodeOrToken.start, nodeOrToken.end];
    }
    function parseOffsetToUnistPoint(acornOffset) {
      let sourceOffset = acornOffset - prefix.length;
      if (sourceOffset < 0) {
        sourceOffset = 0;
      } else if (sourceOffset > source.length) {
        sourceOffset = source.length;
      }
      let point3 = relativeToPoint(collection.stops, sourceOffset);
      if (!point3) {
        point3 = {
          line: options.start.line,
          column: options.start.column,
          offset: options.start.offset
        };
      }
      return point3;
    }
  }
  function empty2(value) {
    return /^\s*$/.test(value.replace(/\/\*[\s\S]*?\*\//g, "").replace(/\/\/[^\r\n]*(\r\n|\n|\r)/g, ""));
  }
  function collect(events, tokenTypes) {
    const result = {
      value: "",
      stops: []
    };
    let index2 = -1;
    while (++index2 < events.length) {
      const event = events[index2];
      if (event[0] === "enter") {
        const type = event[1].type;
        if (type === "lineEnding" || tokenTypes.includes(type)) {
          const chunks = event[2].sliceStream(event[1]);
          while (chunks.length > 0 && chunks[0] === -1) {
            chunks.shift();
          }
          const value = serializeChunks(chunks);
          result.stops.push([result.value.length, event[1].start]);
          result.value += value;
          result.stops.push([result.value.length, event[1].end]);
        }
      }
    }
    return result;
  }
  function relativeToPoint(stops, relative) {
    let index2 = 0;
    while (index2 < stops.length && stops[index2][0] <= relative) {
      index2 += 1;
    }
    if (index2 === 0) {
      return void 0;
    }
    const [stopRelative, stopAbsolute] = stops[index2 - 1];
    const rest = relative - stopRelative;
    return {
      line: stopAbsolute.line,
      column: stopAbsolute.column + rest,
      offset: stopAbsolute.offset + rest
    };
  }
  function serializeChunks(chunks) {
    let index2 = -1;
    const result = [];
    let atTab;
    while (++index2 < chunks.length) {
      const chunk = chunks[index2];
      let value;
      if (typeof chunk === "string") {
        value = chunk;
      } else switch (chunk) {
        case -5: {
          value = "\r";
          break;
        }
        case -4: {
          value = "\n";
          break;
        }
        case -3: {
          value = "\r\n";
          break;
        }
        case -2: {
          value = "	";
          break;
        }
        /* c8 ignore next 6 */
        case -1: {
          if (atTab) continue;
          value = " ";
          break;
        }
        default: {
          value = String.fromCharCode(chunk);
        }
      }
      atTab = chunk === -2;
      result.push(value);
    }
    return result.join("");
  }

  // ../../../ui/node_modules/.pnpm/unist-util-position-from-estree@2.0.0/node_modules/unist-util-position-from-estree/lib/index.js
  function positionFromEstree(node2) {
    const nodeLike = node2 || {};
    const loc = nodeLike.loc || {};
    const range = nodeLike.range || [void 0, void 0];
    const start2 = pointOrUndefined(loc.start, range[0] || nodeLike.start);
    const end = pointOrUndefined(loc.end, range[1] || nodeLike.end);
    if (start2 && end) {
      return { start: start2, end };
    }
  }
  function pointOrUndefined(estreePoint, estreeOffset) {
    if (estreePoint && typeof estreePoint === "object") {
      const line = "line" in estreePoint ? numberOrUndefined(estreePoint.line) : void 0;
      const column = "column" in estreePoint ? numberOrUndefined(estreePoint.column) : void 0;
      if (line && column !== void 0) {
        return {
          line,
          column: column + 1,
          offset: numberOrUndefined(estreeOffset)
        };
      }
    }
  }
  function numberOrUndefined(value) {
    return typeof value === "number" && value > -1 ? value : void 0;
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-mdx-expression@2.0.3/node_modules/micromark-factory-mdx-expression/index.js
  var indentSize = 2;
  var trouble = "https://github.com/micromark/micromark-extension-mdx-expression/tree/main/packages/micromark-extension-mdx-expression";
  var unexpectedEndOfFileHash = "#unexpected-end-of-file-in-expression-expected-a-corresponding-closing-brace-for-";
  var unexpectedLazyHash = "#unexpected-lazy-line-in-expression-in-container-expected-line-to-be-prefixed";
  var nonSpreadHash = "#unexpected-type-in-code-expected-an-object-spread-spread";
  var spreadExtraHash = "#unexpected-extra-content-in-spread-only-a-single-spread-is-supported";
  var acornHash = "#could-not-parse-expression-with-acorn";
  function factoryMdxExpression(effects, ok3, type, markerType, chunkType, acorn, acornOptions, addResult, spread, allowEmpty, allowLazy) {
    const self2 = this;
    const eventStart = this.events.length + 3;
    let size = 0;
    let pointStart;
    let lastCrash;
    return start2;
    function start2(code3) {
      effects.enter(type);
      effects.enter(markerType);
      effects.consume(code3);
      effects.exit(markerType);
      pointStart = self2.now();
      return before;
    }
    function before(code3) {
      if (code3 === null) {
        if (lastCrash) throw lastCrash;
        const error = new VFileMessage("Unexpected end of file in expression, expected a corresponding closing brace for `{`", {
          place: self2.now(),
          ruleId: "unexpected-eof",
          source: "micromark-extension-mdx-expression"
        });
        error.url = trouble + unexpectedEndOfFileHash;
        throw error;
      }
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return eolAfter;
      }
      if (code3 === 125 && size === 0) {
        const next = acorn ? mdxExpressionParse.call(self2, acorn, acornOptions, chunkType, eventStart, pointStart, allowEmpty || false, spread || false) : {
          type: "ok",
          estree: void 0
        };
        if (next.type === "ok") {
          effects.enter(markerType);
          effects.consume(code3);
          effects.exit(markerType);
          const token = effects.exit(type);
          if (addResult && next.estree) {
            Object.assign(token, {
              estree: next.estree
            });
          }
          return ok3;
        }
        lastCrash = next.message;
        effects.enter(chunkType);
        effects.consume(code3);
        return inside;
      }
      effects.enter(chunkType);
      return inside(code3);
    }
    function inside(code3) {
      if (code3 === 125 && size === 0 || code3 === null || markdownLineEnding(code3)) {
        effects.exit(chunkType);
        return before(code3);
      }
      if (code3 === 123 && !acorn) {
        size += 1;
      } else if (code3 === 125) {
        size -= 1;
      }
      effects.consume(code3);
      return inside;
    }
    function eolAfter(code3) {
      const now = self2.now();
      if (now.line !== pointStart.line && !allowLazy && self2.parser.lazy[now.line]) {
        const error = new VFileMessage("Unexpected lazy line in expression in container, expected line to be prefixed with `>` when in a block quote, whitespace when in a list, etc", {
          place: self2.now(),
          ruleId: "unexpected-lazy",
          source: "micromark-extension-mdx-expression"
        });
        error.url = trouble + unexpectedLazyHash;
        throw error;
      }
      if (markdownSpace(code3)) {
        return factorySpace(effects, before, "linePrefix", indentSize + 1)(code3);
      }
      return before(code3);
    }
  }
  function mdxExpressionParse(acorn, acornOptions, chunkType, eventStart, pointStart, allowEmpty, spread) {
    const result = eventsToAcorn(this.events.slice(eventStart), {
      acorn,
      tokenTypes: [chunkType],
      acornOptions,
      start: pointStart,
      expression: true,
      allowEmpty,
      prefix: spread ? "({" : "",
      suffix: spread ? "})" : ""
    });
    const estree = result.estree;
    if (spread && estree) {
      const head = estree.body[0];
      if (head.type !== "ExpressionStatement" || head.expression.type !== "ObjectExpression") {
        const place = positionFromEstree(head);
        const error = new VFileMessage("Unexpected `" + head.type + "` in code: expected an object spread (`{...spread}`)", {
          place: place.start,
          ruleId: "non-spread",
          source: "micromark-extension-mdx-expression"
        });
        error.url = trouble + nonSpreadHash;
        throw error;
      }
      if (head.expression.properties[1]) {
        const place = positionFromEstree(head.expression.properties[1]);
        const error = new VFileMessage("Unexpected extra content in spread: only a single spread is supported", {
          place: place.start,
          ruleId: "spread-extra",
          source: "micromark-extension-mdx-expression"
        });
        error.url = trouble + spreadExtraHash;
        throw error;
      }
      if (head.expression.properties[0] && head.expression.properties[0].type !== "SpreadElement") {
        const place = positionFromEstree(head.expression.properties[0]);
        const error = new VFileMessage("Unexpected `" + head.expression.properties[0].type + "` in code: only spread elements are supported", {
          place: place.start,
          ruleId: "non-spread",
          source: "micromark-extension-mdx-expression"
        });
        error.url = trouble + nonSpreadHash;
        throw error;
      }
    }
    if (result.error) {
      const error = new VFileMessage("Could not parse expression with acorn", {
        cause: result.error,
        place: {
          line: result.error.loc.line,
          column: result.error.loc.column + 1,
          offset: result.error.pos
        },
        ruleId: "acorn",
        source: "micromark-extension-mdx-expression"
      });
      error.url = trouble + acornHash;
      return {
        type: "nok",
        message: error
      };
    }
    return {
      type: "ok",
      estree
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-expression@3.0.1/node_modules/micromark-extension-mdx-expression/lib/syntax.js
  function mdxExpression(options) {
    const options_ = options || {};
    const addResult = options_.addResult;
    const acorn = options_.acorn;
    const spread = options_.spread;
    let allowEmpty = options_.allowEmpty;
    let acornOptions;
    if (allowEmpty === null || allowEmpty === void 0) {
      allowEmpty = true;
    }
    if (acorn) {
      if (!acorn.parseExpressionAt) {
        throw new Error("Expected a proper `acorn` instance passed in as `options.acorn`");
      }
      acornOptions = Object.assign({
        ecmaVersion: 2024,
        sourceType: "module"
      }, options_.acornOptions);
    } else if (options_.acornOptions || options_.addResult) {
      throw new Error("Expected an `acorn` instance passed in as `options.acorn`");
    }
    return {
      flow: {
        [123]: {
          name: "mdxFlowExpression",
          tokenize: tokenizeFlowExpression,
          concrete: true
        }
      },
      text: {
        [123]: {
          name: "mdxTextExpression",
          tokenize: tokenizeTextExpression
        }
      }
    };
    function tokenizeFlowExpression(effects, ok3, nok) {
      const self2 = this;
      return start2;
      function start2(code3) {
        return before(code3);
      }
      function before(code3) {
        return factoryMdxExpression.call(self2, effects, after, "mdxFlowExpression", "mdxFlowExpressionMarker", "mdxFlowExpressionChunk", acorn, acornOptions, addResult, spread, allowEmpty)(code3);
      }
      function after(code3) {
        return markdownSpace(code3) ? factorySpace(effects, end, "whitespace")(code3) : end(code3);
      }
      function end(code3) {
        const lessThanValue = self2.parser.constructs.flow[60];
        const constructs2 = Array.isArray(lessThanValue) ? lessThanValue : (
          /* c8 ignore next 3 -- always a list when normalized. */
          lessThanValue ? [lessThanValue] : []
        );
        const jsxTag = constructs2.find(function(d) {
          return d.name === "mdxJsxFlowTag";
        });
        if (code3 === 60 && jsxTag) {
          return effects.attempt(jsxTag, end, nok)(code3);
        }
        return code3 === null || markdownLineEnding(code3) ? ok3(code3) : nok(code3);
      }
    }
    function tokenizeTextExpression(effects, ok3) {
      const self2 = this;
      return start2;
      function start2(code3) {
        return factoryMdxExpression.call(self2, effects, ok3, "mdxTextExpression", "mdxTextExpressionMarker", "mdxTextExpressionChunk", acorn, acornOptions, addResult, spread, allowEmpty, true)(code3);
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/estree-util-is-identifier-name@3.0.0/node_modules/estree-util-is-identifier-name/lib/index.js
  var startRe = /[\u0024\u0041-\u005a\u005f\u0061-\u007a\u00aa\u00b5\u00ba\u00c0-\u00d6\u00d8-\u00f6\u00f8-\u02c1\u02c6-\u02d1\u02e0-\u02e4\u02ec\u02ee\u0370-\u0374\u0376-\u0377\u037a-\u037d\u037f\u0386\u0388-\u038a\u038c\u038e-\u03a1\u03a3-\u03f5\u03f7-\u0481\u048a-\u052f\u0531-\u0556\u0559\u0560-\u0588\u05d0-\u05ea\u05ef-\u05f2\u0620-\u064a\u066e-\u066f\u0671-\u06d3\u06d5\u06e5-\u06e6\u06ee-\u06ef\u06fa-\u06fc\u06ff\u0710\u0712-\u072f\u074d-\u07a5\u07b1\u07ca-\u07ea\u07f4-\u07f5\u07fa\u0800-\u0815\u081a\u0824\u0828\u0840-\u0858\u0860-\u086a\u0870-\u0887\u0889-\u088f\u08a0-\u08c9\u0904-\u0939\u093d\u0950\u0958-\u0961\u0971-\u0980\u0985-\u098c\u098f-\u0990\u0993-\u09a8\u09aa-\u09b0\u09b2\u09b6-\u09b9\u09bd\u09ce\u09dc-\u09dd\u09df-\u09e1\u09f0-\u09f1\u09fc\u0a05-\u0a0a\u0a0f-\u0a10\u0a13-\u0a28\u0a2a-\u0a30\u0a32-\u0a33\u0a35-\u0a36\u0a38-\u0a39\u0a59-\u0a5c\u0a5e\u0a72-\u0a74\u0a85-\u0a8d\u0a8f-\u0a91\u0a93-\u0aa8\u0aaa-\u0ab0\u0ab2-\u0ab3\u0ab5-\u0ab9\u0abd\u0ad0\u0ae0-\u0ae1\u0af9\u0b05-\u0b0c\u0b0f-\u0b10\u0b13-\u0b28\u0b2a-\u0b30\u0b32-\u0b33\u0b35-\u0b39\u0b3d\u0b5c-\u0b5d\u0b5f-\u0b61\u0b71\u0b83\u0b85-\u0b8a\u0b8e-\u0b90\u0b92-\u0b95\u0b99-\u0b9a\u0b9c\u0b9e-\u0b9f\u0ba3-\u0ba4\u0ba8-\u0baa\u0bae-\u0bb9\u0bd0\u0c05-\u0c0c\u0c0e-\u0c10\u0c12-\u0c28\u0c2a-\u0c39\u0c3d\u0c58-\u0c5a\u0c5c-\u0c5d\u0c60-\u0c61\u0c80\u0c85-\u0c8c\u0c8e-\u0c90\u0c92-\u0ca8\u0caa-\u0cb3\u0cb5-\u0cb9\u0cbd\u0cdc-\u0cde\u0ce0-\u0ce1\u0cf1-\u0cf2\u0d04-\u0d0c\u0d0e-\u0d10\u0d12-\u0d3a\u0d3d\u0d4e\u0d54-\u0d56\u0d5f-\u0d61\u0d7a-\u0d7f\u0d85-\u0d96\u0d9a-\u0db1\u0db3-\u0dbb\u0dbd\u0dc0-\u0dc6\u0e01-\u0e30\u0e32-\u0e33\u0e40-\u0e46\u0e81-\u0e82\u0e84\u0e86-\u0e8a\u0e8c-\u0ea3\u0ea5\u0ea7-\u0eb0\u0eb2-\u0eb3\u0ebd\u0ec0-\u0ec4\u0ec6\u0edc-\u0edf\u0f00\u0f40-\u0f47\u0f49-\u0f6c\u0f88-\u0f8c\u1000-\u102a\u103f\u1050-\u1055\u105a-\u105d\u1061\u1065-\u1066\u106e-\u1070\u1075-\u1081\u108e\u10a0-\u10c5\u10c7\u10cd\u10d0-\u10fa\u10fc-\u1248\u124a-\u124d\u1250-\u1256\u1258\u125a-\u125d\u1260-\u1288\u128a-\u128d\u1290-\u12b0\u12b2-\u12b5\u12b8-\u12be\u12c0\u12c2-\u12c5\u12c8-\u12d6\u12d8-\u1310\u1312-\u1315\u1318-\u135a\u1380-\u138f\u13a0-\u13f5\u13f8-\u13fd\u1401-\u166c\u166f-\u167f\u1681-\u169a\u16a0-\u16ea\u16ee-\u16f8\u1700-\u1711\u171f-\u1731\u1740-\u1751\u1760-\u176c\u176e-\u1770\u1780-\u17b3\u17d7\u17dc\u1820-\u1878\u1880-\u18a8\u18aa\u18b0-\u18f5\u1900-\u191e\u1950-\u196d\u1970-\u1974\u1980-\u19ab\u19b0-\u19c9\u1a00-\u1a16\u1a20-\u1a54\u1aa7\u1b05-\u1b33\u1b45-\u1b4c\u1b83-\u1ba0\u1bae-\u1baf\u1bba-\u1be5\u1c00-\u1c23\u1c4d-\u1c4f\u1c5a-\u1c7d\u1c80-\u1c8a\u1c90-\u1cba\u1cbd-\u1cbf\u1ce9-\u1cec\u1cee-\u1cf3\u1cf5-\u1cf6\u1cfa\u1d00-\u1dbf\u1e00-\u1f15\u1f18-\u1f1d\u1f20-\u1f45\u1f48-\u1f4d\u1f50-\u1f57\u1f59\u1f5b\u1f5d\u1f5f-\u1f7d\u1f80-\u1fb4\u1fb6-\u1fbc\u1fbe\u1fc2-\u1fc4\u1fc6-\u1fcc\u1fd0-\u1fd3\u1fd6-\u1fdb\u1fe0-\u1fec\u1ff2-\u1ff4\u1ff6-\u1ffc\u2071\u207f\u2090-\u209c\u2102\u2107\u210a-\u2113\u2115\u2118-\u211d\u2124\u2126\u2128\u212a-\u2139\u213c-\u213f\u2145-\u2149\u214e\u2160-\u2188\u2c00-\u2ce4\u2ceb-\u2cee\u2cf2-\u2cf3\u2d00-\u2d25\u2d27\u2d2d\u2d30-\u2d67\u2d6f\u2d80-\u2d96\u2da0-\u2da6\u2da8-\u2dae\u2db0-\u2db6\u2db8-\u2dbe\u2dc0-\u2dc6\u2dc8-\u2dce\u2dd0-\u2dd6\u2dd8-\u2dde\u3005-\u3007\u3021-\u3029\u3031-\u3035\u3038-\u303c\u3041-\u3096\u309b-\u309f\u30a1-\u30fa\u30fc-\u30ff\u3105-\u312f\u3131-\u318e\u31a0-\u31bf\u31f0-\u31ff\u3400-\u4dbf\u4e00-\ua48c\ua4d0-\ua4fd\ua500-\ua60c\ua610-\ua61f\ua62a-\ua62b\ua640-\ua66e\ua67f-\ua69d\ua6a0-\ua6ef\ua717-\ua71f\ua722-\ua788\ua78b-\ua7dc\ua7f1-\ua801\ua803-\ua805\ua807-\ua80a\ua80c-\ua822\ua840-\ua873\ua882-\ua8b3\ua8f2-\ua8f7\ua8fb\ua8fd-\ua8fe\ua90a-\ua925\ua930-\ua946\ua960-\ua97c\ua984-\ua9b2\ua9cf\ua9e0-\ua9e4\ua9e6-\ua9ef\ua9fa-\ua9fe\uaa00-\uaa28\uaa40-\uaa42\uaa44-\uaa4b\uaa60-\uaa76\uaa7a\uaa7e-\uaaaf\uaab1\uaab5-\uaab6\uaab9-\uaabd\uaac0\uaac2\uaadb-\uaadd\uaae0-\uaaea\uaaf2-\uaaf4\uab01-\uab06\uab09-\uab0e\uab11-\uab16\uab20-\uab26\uab28-\uab2e\uab30-\uab5a\uab5c-\uab69\uab70-\uabe2\uac00-\ud7a3\ud7b0-\ud7c6\ud7cb-\ud7fb\uf900-\ufa6d\ufa70-\ufad9\ufb00-\ufb06\ufb13-\ufb17\ufb1d\ufb1f-\ufb28\ufb2a-\ufb36\ufb38-\ufb3c\ufb3e\ufb40-\ufb41\ufb43-\ufb44\ufb46-\ufbb1\ufbd3-\ufd3d\ufd50-\ufd8f\ufd92-\ufdc7\ufdf0-\ufdfb\ufe70-\ufe74\ufe76-\ufefc\uff21-\uff3a\uff41-\uff5a\uff66-\uffbe\uffc2-\uffc7\uffca-\uffcf\uffd2-\uffd7\uffda-\uffdc]/;
  var contRe = /[\u0024\u0030-\u0039\u0041-\u005a\u005f\u0061-\u007a\u00aa\u00b5\u00b7\u00ba\u00c0-\u00d6\u00d8-\u00f6\u00f8-\u02c1\u02c6-\u02d1\u02e0-\u02e4\u02ec\u02ee\u0300-\u0374\u0376-\u0377\u037a-\u037d\u037f\u0386-\u038a\u038c\u038e-\u03a1\u03a3-\u03f5\u03f7-\u0481\u0483-\u0487\u048a-\u052f\u0531-\u0556\u0559\u0560-\u0588\u0591-\u05bd\u05bf\u05c1-\u05c2\u05c4-\u05c5\u05c7\u05d0-\u05ea\u05ef-\u05f2\u0610-\u061a\u0620-\u0669\u066e-\u06d3\u06d5-\u06dc\u06df-\u06e8\u06ea-\u06fc\u06ff\u0710-\u074a\u074d-\u07b1\u07c0-\u07f5\u07fa\u07fd\u0800-\u082d\u0840-\u085b\u0860-\u086a\u0870-\u0887\u0889-\u088f\u0897-\u08e1\u08e3-\u0963\u0966-\u096f\u0971-\u0983\u0985-\u098c\u098f-\u0990\u0993-\u09a8\u09aa-\u09b0\u09b2\u09b6-\u09b9\u09bc-\u09c4\u09c7-\u09c8\u09cb-\u09ce\u09d7\u09dc-\u09dd\u09df-\u09e3\u09e6-\u09f1\u09fc\u09fe\u0a01-\u0a03\u0a05-\u0a0a\u0a0f-\u0a10\u0a13-\u0a28\u0a2a-\u0a30\u0a32-\u0a33\u0a35-\u0a36\u0a38-\u0a39\u0a3c\u0a3e-\u0a42\u0a47-\u0a48\u0a4b-\u0a4d\u0a51\u0a59-\u0a5c\u0a5e\u0a66-\u0a75\u0a81-\u0a83\u0a85-\u0a8d\u0a8f-\u0a91\u0a93-\u0aa8\u0aaa-\u0ab0\u0ab2-\u0ab3\u0ab5-\u0ab9\u0abc-\u0ac5\u0ac7-\u0ac9\u0acb-\u0acd\u0ad0\u0ae0-\u0ae3\u0ae6-\u0aef\u0af9-\u0aff\u0b01-\u0b03\u0b05-\u0b0c\u0b0f-\u0b10\u0b13-\u0b28\u0b2a-\u0b30\u0b32-\u0b33\u0b35-\u0b39\u0b3c-\u0b44\u0b47-\u0b48\u0b4b-\u0b4d\u0b55-\u0b57\u0b5c-\u0b5d\u0b5f-\u0b63\u0b66-\u0b6f\u0b71\u0b82-\u0b83\u0b85-\u0b8a\u0b8e-\u0b90\u0b92-\u0b95\u0b99-\u0b9a\u0b9c\u0b9e-\u0b9f\u0ba3-\u0ba4\u0ba8-\u0baa\u0bae-\u0bb9\u0bbe-\u0bc2\u0bc6-\u0bc8\u0bca-\u0bcd\u0bd0\u0bd7\u0be6-\u0bef\u0c00-\u0c0c\u0c0e-\u0c10\u0c12-\u0c28\u0c2a-\u0c39\u0c3c-\u0c44\u0c46-\u0c48\u0c4a-\u0c4d\u0c55-\u0c56\u0c58-\u0c5a\u0c5c-\u0c5d\u0c60-\u0c63\u0c66-\u0c6f\u0c80-\u0c83\u0c85-\u0c8c\u0c8e-\u0c90\u0c92-\u0ca8\u0caa-\u0cb3\u0cb5-\u0cb9\u0cbc-\u0cc4\u0cc6-\u0cc8\u0cca-\u0ccd\u0cd5-\u0cd6\u0cdc-\u0cde\u0ce0-\u0ce3\u0ce6-\u0cef\u0cf1-\u0cf3\u0d00-\u0d0c\u0d0e-\u0d10\u0d12-\u0d44\u0d46-\u0d48\u0d4a-\u0d4e\u0d54-\u0d57\u0d5f-\u0d63\u0d66-\u0d6f\u0d7a-\u0d7f\u0d81-\u0d83\u0d85-\u0d96\u0d9a-\u0db1\u0db3-\u0dbb\u0dbd\u0dc0-\u0dc6\u0dca\u0dcf-\u0dd4\u0dd6\u0dd8-\u0ddf\u0de6-\u0def\u0df2-\u0df3\u0e01-\u0e3a\u0e40-\u0e4e\u0e50-\u0e59\u0e81-\u0e82\u0e84\u0e86-\u0e8a\u0e8c-\u0ea3\u0ea5\u0ea7-\u0ebd\u0ec0-\u0ec4\u0ec6\u0ec8-\u0ece\u0ed0-\u0ed9\u0edc-\u0edf\u0f00\u0f18-\u0f19\u0f20-\u0f29\u0f35\u0f37\u0f39\u0f3e-\u0f47\u0f49-\u0f6c\u0f71-\u0f84\u0f86-\u0f97\u0f99-\u0fbc\u0fc6\u1000-\u1049\u1050-\u109d\u10a0-\u10c5\u10c7\u10cd\u10d0-\u10fa\u10fc-\u1248\u124a-\u124d\u1250-\u1256\u1258\u125a-\u125d\u1260-\u1288\u128a-\u128d\u1290-\u12b0\u12b2-\u12b5\u12b8-\u12be\u12c0\u12c2-\u12c5\u12c8-\u12d6\u12d8-\u1310\u1312-\u1315\u1318-\u135a\u135d-\u135f\u1369-\u1371\u1380-\u138f\u13a0-\u13f5\u13f8-\u13fd\u1401-\u166c\u166f-\u167f\u1681-\u169a\u16a0-\u16ea\u16ee-\u16f8\u1700-\u1715\u171f-\u1734\u1740-\u1753\u1760-\u176c\u176e-\u1770\u1772-\u1773\u1780-\u17d3\u17d7\u17dc-\u17dd\u17e0-\u17e9\u180b-\u180d\u180f-\u1819\u1820-\u1878\u1880-\u18aa\u18b0-\u18f5\u1900-\u191e\u1920-\u192b\u1930-\u193b\u1946-\u196d\u1970-\u1974\u1980-\u19ab\u19b0-\u19c9\u19d0-\u19da\u1a00-\u1a1b\u1a20-\u1a5e\u1a60-\u1a7c\u1a7f-\u1a89\u1a90-\u1a99\u1aa7\u1ab0-\u1abd\u1abf-\u1add\u1ae0-\u1aeb\u1b00-\u1b4c\u1b50-\u1b59\u1b6b-\u1b73\u1b80-\u1bf3\u1c00-\u1c37\u1c40-\u1c49\u1c4d-\u1c7d\u1c80-\u1c8a\u1c90-\u1cba\u1cbd-\u1cbf\u1cd0-\u1cd2\u1cd4-\u1cfa\u1d00-\u1f15\u1f18-\u1f1d\u1f20-\u1f45\u1f48-\u1f4d\u1f50-\u1f57\u1f59\u1f5b\u1f5d\u1f5f-\u1f7d\u1f80-\u1fb4\u1fb6-\u1fbc\u1fbe\u1fc2-\u1fc4\u1fc6-\u1fcc\u1fd0-\u1fd3\u1fd6-\u1fdb\u1fe0-\u1fec\u1ff2-\u1ff4\u1ff6-\u1ffc\u200c-\u200d\u203f-\u2040\u2054\u2071\u207f\u2090-\u209c\u20d0-\u20dc\u20e1\u20e5-\u20f0\u2102\u2107\u210a-\u2113\u2115\u2118-\u211d\u2124\u2126\u2128\u212a-\u2139\u213c-\u213f\u2145-\u2149\u214e\u2160-\u2188\u2c00-\u2ce4\u2ceb-\u2cf3\u2d00-\u2d25\u2d27\u2d2d\u2d30-\u2d67\u2d6f\u2d7f-\u2d96\u2da0-\u2da6\u2da8-\u2dae\u2db0-\u2db6\u2db8-\u2dbe\u2dc0-\u2dc6\u2dc8-\u2dce\u2dd0-\u2dd6\u2dd8-\u2dde\u2de0-\u2dff\u3005-\u3007\u3021-\u302f\u3031-\u3035\u3038-\u303c\u3041-\u3096\u3099-\u309f\u30a1-\u30ff\u3105-\u312f\u3131-\u318e\u31a0-\u31bf\u31f0-\u31ff\u3400-\u4dbf\u4e00-\ua48c\ua4d0-\ua4fd\ua500-\ua60c\ua610-\ua62b\ua640-\ua66f\ua674-\ua67d\ua67f-\ua6f1\ua717-\ua71f\ua722-\ua788\ua78b-\ua7dc\ua7f1-\ua827\ua82c\ua840-\ua873\ua880-\ua8c5\ua8d0-\ua8d9\ua8e0-\ua8f7\ua8fb\ua8fd-\ua92d\ua930-\ua953\ua960-\ua97c\ua980-\ua9c0\ua9cf-\ua9d9\ua9e0-\ua9fe\uaa00-\uaa36\uaa40-\uaa4d\uaa50-\uaa59\uaa60-\uaa76\uaa7a-\uaac2\uaadb-\uaadd\uaae0-\uaaef\uaaf2-\uaaf6\uab01-\uab06\uab09-\uab0e\uab11-\uab16\uab20-\uab26\uab28-\uab2e\uab30-\uab5a\uab5c-\uab69\uab70-\uabea\uabec-\uabed\uabf0-\uabf9\uac00-\ud7a3\ud7b0-\ud7c6\ud7cb-\ud7fb\uf900-\ufa6d\ufa70-\ufad9\ufb00-\ufb06\ufb13-\ufb17\ufb1d-\ufb28\ufb2a-\ufb36\ufb38-\ufb3c\ufb3e\ufb40-\ufb41\ufb43-\ufb44\ufb46-\ufbb1\ufbd3-\ufd3d\ufd50-\ufd8f\ufd92-\ufdc7\ufdf0-\ufdfb\ufe00-\ufe0f\ufe20-\ufe2f\ufe33-\ufe34\ufe4d-\ufe4f\ufe70-\ufe74\ufe76-\ufefc\uff10-\uff19\uff21-\uff3a\uff3f\uff41-\uff5a\uff65-\uffbe\uffc2-\uffc7\uffca-\uffcf\uffd2-\uffd7\uffda-\uffdc]/;
  var contReJsx = /[-\u0024\u0030-\u0039\u0041-\u005a\u005f\u0061-\u007a\u00aa\u00b5\u00b7\u00ba\u00c0-\u00d6\u00d8-\u00f6\u00f8-\u02c1\u02c6-\u02d1\u02e0-\u02e4\u02ec\u02ee\u0300-\u0374\u0376-\u0377\u037a-\u037d\u037f\u0386-\u038a\u038c\u038e-\u03a1\u03a3-\u03f5\u03f7-\u0481\u0483-\u0487\u048a-\u052f\u0531-\u0556\u0559\u0560-\u0588\u0591-\u05bd\u05bf\u05c1-\u05c2\u05c4-\u05c5\u05c7\u05d0-\u05ea\u05ef-\u05f2\u0610-\u061a\u0620-\u0669\u066e-\u06d3\u06d5-\u06dc\u06df-\u06e8\u06ea-\u06fc\u06ff\u0710-\u074a\u074d-\u07b1\u07c0-\u07f5\u07fa\u07fd\u0800-\u082d\u0840-\u085b\u0860-\u086a\u0870-\u0887\u0889-\u088f\u0897-\u08e1\u08e3-\u0963\u0966-\u096f\u0971-\u0983\u0985-\u098c\u098f-\u0990\u0993-\u09a8\u09aa-\u09b0\u09b2\u09b6-\u09b9\u09bc-\u09c4\u09c7-\u09c8\u09cb-\u09ce\u09d7\u09dc-\u09dd\u09df-\u09e3\u09e6-\u09f1\u09fc\u09fe\u0a01-\u0a03\u0a05-\u0a0a\u0a0f-\u0a10\u0a13-\u0a28\u0a2a-\u0a30\u0a32-\u0a33\u0a35-\u0a36\u0a38-\u0a39\u0a3c\u0a3e-\u0a42\u0a47-\u0a48\u0a4b-\u0a4d\u0a51\u0a59-\u0a5c\u0a5e\u0a66-\u0a75\u0a81-\u0a83\u0a85-\u0a8d\u0a8f-\u0a91\u0a93-\u0aa8\u0aaa-\u0ab0\u0ab2-\u0ab3\u0ab5-\u0ab9\u0abc-\u0ac5\u0ac7-\u0ac9\u0acb-\u0acd\u0ad0\u0ae0-\u0ae3\u0ae6-\u0aef\u0af9-\u0aff\u0b01-\u0b03\u0b05-\u0b0c\u0b0f-\u0b10\u0b13-\u0b28\u0b2a-\u0b30\u0b32-\u0b33\u0b35-\u0b39\u0b3c-\u0b44\u0b47-\u0b48\u0b4b-\u0b4d\u0b55-\u0b57\u0b5c-\u0b5d\u0b5f-\u0b63\u0b66-\u0b6f\u0b71\u0b82-\u0b83\u0b85-\u0b8a\u0b8e-\u0b90\u0b92-\u0b95\u0b99-\u0b9a\u0b9c\u0b9e-\u0b9f\u0ba3-\u0ba4\u0ba8-\u0baa\u0bae-\u0bb9\u0bbe-\u0bc2\u0bc6-\u0bc8\u0bca-\u0bcd\u0bd0\u0bd7\u0be6-\u0bef\u0c00-\u0c0c\u0c0e-\u0c10\u0c12-\u0c28\u0c2a-\u0c39\u0c3c-\u0c44\u0c46-\u0c48\u0c4a-\u0c4d\u0c55-\u0c56\u0c58-\u0c5a\u0c5c-\u0c5d\u0c60-\u0c63\u0c66-\u0c6f\u0c80-\u0c83\u0c85-\u0c8c\u0c8e-\u0c90\u0c92-\u0ca8\u0caa-\u0cb3\u0cb5-\u0cb9\u0cbc-\u0cc4\u0cc6-\u0cc8\u0cca-\u0ccd\u0cd5-\u0cd6\u0cdc-\u0cde\u0ce0-\u0ce3\u0ce6-\u0cef\u0cf1-\u0cf3\u0d00-\u0d0c\u0d0e-\u0d10\u0d12-\u0d44\u0d46-\u0d48\u0d4a-\u0d4e\u0d54-\u0d57\u0d5f-\u0d63\u0d66-\u0d6f\u0d7a-\u0d7f\u0d81-\u0d83\u0d85-\u0d96\u0d9a-\u0db1\u0db3-\u0dbb\u0dbd\u0dc0-\u0dc6\u0dca\u0dcf-\u0dd4\u0dd6\u0dd8-\u0ddf\u0de6-\u0def\u0df2-\u0df3\u0e01-\u0e3a\u0e40-\u0e4e\u0e50-\u0e59\u0e81-\u0e82\u0e84\u0e86-\u0e8a\u0e8c-\u0ea3\u0ea5\u0ea7-\u0ebd\u0ec0-\u0ec4\u0ec6\u0ec8-\u0ece\u0ed0-\u0ed9\u0edc-\u0edf\u0f00\u0f18-\u0f19\u0f20-\u0f29\u0f35\u0f37\u0f39\u0f3e-\u0f47\u0f49-\u0f6c\u0f71-\u0f84\u0f86-\u0f97\u0f99-\u0fbc\u0fc6\u1000-\u1049\u1050-\u109d\u10a0-\u10c5\u10c7\u10cd\u10d0-\u10fa\u10fc-\u1248\u124a-\u124d\u1250-\u1256\u1258\u125a-\u125d\u1260-\u1288\u128a-\u128d\u1290-\u12b0\u12b2-\u12b5\u12b8-\u12be\u12c0\u12c2-\u12c5\u12c8-\u12d6\u12d8-\u1310\u1312-\u1315\u1318-\u135a\u135d-\u135f\u1369-\u1371\u1380-\u138f\u13a0-\u13f5\u13f8-\u13fd\u1401-\u166c\u166f-\u167f\u1681-\u169a\u16a0-\u16ea\u16ee-\u16f8\u1700-\u1715\u171f-\u1734\u1740-\u1753\u1760-\u176c\u176e-\u1770\u1772-\u1773\u1780-\u17d3\u17d7\u17dc-\u17dd\u17e0-\u17e9\u180b-\u180d\u180f-\u1819\u1820-\u1878\u1880-\u18aa\u18b0-\u18f5\u1900-\u191e\u1920-\u192b\u1930-\u193b\u1946-\u196d\u1970-\u1974\u1980-\u19ab\u19b0-\u19c9\u19d0-\u19da\u1a00-\u1a1b\u1a20-\u1a5e\u1a60-\u1a7c\u1a7f-\u1a89\u1a90-\u1a99\u1aa7\u1ab0-\u1abd\u1abf-\u1add\u1ae0-\u1aeb\u1b00-\u1b4c\u1b50-\u1b59\u1b6b-\u1b73\u1b80-\u1bf3\u1c00-\u1c37\u1c40-\u1c49\u1c4d-\u1c7d\u1c80-\u1c8a\u1c90-\u1cba\u1cbd-\u1cbf\u1cd0-\u1cd2\u1cd4-\u1cfa\u1d00-\u1f15\u1f18-\u1f1d\u1f20-\u1f45\u1f48-\u1f4d\u1f50-\u1f57\u1f59\u1f5b\u1f5d\u1f5f-\u1f7d\u1f80-\u1fb4\u1fb6-\u1fbc\u1fbe\u1fc2-\u1fc4\u1fc6-\u1fcc\u1fd0-\u1fd3\u1fd6-\u1fdb\u1fe0-\u1fec\u1ff2-\u1ff4\u1ff6-\u1ffc\u200c-\u200d\u203f-\u2040\u2054\u2071\u207f\u2090-\u209c\u20d0-\u20dc\u20e1\u20e5-\u20f0\u2102\u2107\u210a-\u2113\u2115\u2118-\u211d\u2124\u2126\u2128\u212a-\u2139\u213c-\u213f\u2145-\u2149\u214e\u2160-\u2188\u2c00-\u2ce4\u2ceb-\u2cf3\u2d00-\u2d25\u2d27\u2d2d\u2d30-\u2d67\u2d6f\u2d7f-\u2d96\u2da0-\u2da6\u2da8-\u2dae\u2db0-\u2db6\u2db8-\u2dbe\u2dc0-\u2dc6\u2dc8-\u2dce\u2dd0-\u2dd6\u2dd8-\u2dde\u2de0-\u2dff\u3005-\u3007\u3021-\u302f\u3031-\u3035\u3038-\u303c\u3041-\u3096\u3099-\u309f\u30a1-\u30ff\u3105-\u312f\u3131-\u318e\u31a0-\u31bf\u31f0-\u31ff\u3400-\u4dbf\u4e00-\ua48c\ua4d0-\ua4fd\ua500-\ua60c\ua610-\ua62b\ua640-\ua66f\ua674-\ua67d\ua67f-\ua6f1\ua717-\ua71f\ua722-\ua788\ua78b-\ua7dc\ua7f1-\ua827\ua82c\ua840-\ua873\ua880-\ua8c5\ua8d0-\ua8d9\ua8e0-\ua8f7\ua8fb\ua8fd-\ua92d\ua930-\ua953\ua960-\ua97c\ua980-\ua9c0\ua9cf-\ua9d9\ua9e0-\ua9fe\uaa00-\uaa36\uaa40-\uaa4d\uaa50-\uaa59\uaa60-\uaa76\uaa7a-\uaac2\uaadb-\uaadd\uaae0-\uaaef\uaaf2-\uaaf6\uab01-\uab06\uab09-\uab0e\uab11-\uab16\uab20-\uab26\uab28-\uab2e\uab30-\uab5a\uab5c-\uab69\uab70-\uabea\uabec-\uabed\uabf0-\uabf9\uac00-\ud7a3\ud7b0-\ud7c6\ud7cb-\ud7fb\uf900-\ufa6d\ufa70-\ufad9\ufb00-\ufb06\ufb13-\ufb17\ufb1d-\ufb28\ufb2a-\ufb36\ufb38-\ufb3c\ufb3e\ufb40-\ufb41\ufb43-\ufb44\ufb46-\ufbb1\ufbd3-\ufd3d\ufd50-\ufd8f\ufd92-\ufdc7\ufdf0-\ufdfb\ufe00-\ufe0f\ufe20-\ufe2f\ufe33-\ufe34\ufe4d-\ufe4f\ufe70-\ufe74\ufe76-\ufefc\uff10-\uff19\uff21-\uff3a\uff3f\uff41-\uff5a\uff65-\uffbe\uffc2-\uffc7\uffca-\uffcf\uffd2-\uffd7\uffda-\uffdc]/;
  var emptyOptions = {};
  function start(code3) {
    return code3 ? startRe.test(String.fromCodePoint(code3)) : false;
  }
  function cont(code3, options) {
    const settings = options || emptyOptions;
    const re = settings.jsx ? contReJsx : contRe;
    return code3 ? re.test(String.fromCodePoint(code3)) : false;
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-jsx@3.0.2/node_modules/micromark-extension-mdx-jsx/lib/factory-tag.js
  var trouble2 = "https://github.com/micromark/micromark-extension-mdx-jsx";
  function factoryTag(effects, ok3, nok, acorn, acornOptions, addResult, allowLazy, tagType, tagMarkerType, tagClosingMarkerType, tagSelfClosingMarker, tagNameType, tagNamePrimaryType, tagNameMemberMarkerType, tagNameMemberType, tagNamePrefixMarkerType, tagNameLocalType, tagExpressionAttributeType, tagExpressionAttributeMarkerType, tagExpressionAttributeValueType, tagAttributeType, tagAttributeNameType, tagAttributeNamePrimaryType, tagAttributeNamePrefixMarkerType, tagAttributeNameLocalType, tagAttributeInitializerMarkerType, tagAttributeValueLiteralType, tagAttributeValueLiteralMarkerType, tagAttributeValueLiteralValueType, tagAttributeValueExpressionType, tagAttributeValueExpressionMarkerType, tagAttributeValueExpressionValueType) {
    const self2 = this;
    let returnState;
    let marker;
    return start2;
    function start2(code3) {
      effects.enter(tagType);
      effects.enter(tagMarkerType);
      effects.consume(code3);
      effects.exit(tagMarkerType);
      return startAfter;
    }
    function startAfter(code3) {
      if (markdownLineEndingOrSpace(code3)) {
        return nok(code3);
      }
      returnState = nameBefore;
      return esWhitespaceStart(code3);
    }
    function nameBefore(code3) {
      if (code3 === 47) {
        effects.enter(tagClosingMarkerType);
        effects.consume(code3);
        effects.exit(tagClosingMarkerType);
        returnState = closingTagNameBefore;
        return esWhitespaceStart;
      }
      if (code3 === 62) {
        return tagEnd(code3);
      }
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagNameType);
        effects.enter(tagNamePrimaryType);
        effects.consume(code3);
        return primaryName;
      }
      crash(code3, "before name", "a character that can start a name, such as a letter, `$`, or `_`" + (code3 === 33 ? " (note: to create a comment in MDX, use `{/* text */}`)" : ""));
    }
    function closingTagNameBefore(code3) {
      if (code3 === 62) {
        return tagEnd(code3);
      }
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagNameType);
        effects.enter(tagNamePrimaryType);
        effects.consume(code3);
        return primaryName;
      }
      crash(code3, "before name", "a character that can start a name, such as a letter, `$`, or `_`" + (code3 === 42 || code3 === 47 ? " (note: JS comments in JSX tags are not supported in MDX)" : ""));
    }
    function primaryName(code3) {
      if (code3 !== null && code3 >= 0 && cont(code3, {
        jsx: true
      })) {
        effects.consume(code3);
        return primaryName;
      }
      if (code3 === 46 || code3 === 47 || code3 === 58 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        effects.exit(tagNamePrimaryType);
        returnState = primaryNameAfter;
        return esWhitespaceStart(code3);
      }
      crash(code3, "in name", "a name character such as letters, digits, `$`, or `_`; whitespace before attributes; or the end of the tag" + (code3 === 64 ? " (note: to create a link in MDX, use `[text](url)`)" : ""));
    }
    function primaryNameAfter(code3) {
      if (code3 === 46) {
        effects.enter(tagNameMemberMarkerType);
        effects.consume(code3);
        effects.exit(tagNameMemberMarkerType);
        returnState = memberNameBefore;
        return esWhitespaceStart;
      }
      if (code3 === 58) {
        effects.enter(tagNamePrefixMarkerType);
        effects.consume(code3);
        effects.exit(tagNamePrefixMarkerType);
        returnState = localNameBefore;
        return esWhitespaceStart;
      }
      if (code3 === 47 || code3 === 62 || code3 === 123 || code3 !== null && code3 >= 0 && start(code3)) {
        effects.exit(tagNameType);
        return attributeBefore(code3);
      }
      crash(code3, "after name", "a character that can start an attribute name, such as a letter, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function memberNameBefore(code3) {
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagNameMemberType);
        effects.consume(code3);
        return memberName;
      }
      crash(code3, "before member name", "a character that can start an attribute name, such as a letter, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function memberName(code3) {
      if (code3 !== null && code3 >= 0 && cont(code3, {
        jsx: true
      })) {
        effects.consume(code3);
        return memberName;
      }
      if (code3 === 46 || code3 === 47 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        effects.exit(tagNameMemberType);
        returnState = memberNameAfter;
        return esWhitespaceStart(code3);
      }
      crash(code3, "in member name", "a name character such as letters, digits, `$`, or `_`; whitespace before attributes; or the end of the tag" + (code3 === 64 ? " (note: to create a link in MDX, use `[text](url)`)" : ""));
    }
    function memberNameAfter(code3) {
      if (code3 === 46) {
        effects.enter(tagNameMemberMarkerType);
        effects.consume(code3);
        effects.exit(tagNameMemberMarkerType);
        returnState = memberNameBefore;
        return esWhitespaceStart;
      }
      if (code3 === 47 || code3 === 62 || code3 === 123 || code3 !== null && code3 >= 0 && start(code3)) {
        effects.exit(tagNameType);
        return attributeBefore(code3);
      }
      crash(code3, "after member name", "a character that can start an attribute name, such as a letter, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function localNameBefore(code3) {
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagNameLocalType);
        effects.consume(code3);
        return localName;
      }
      crash(code3, "before local name", "a character that can start a name, such as a letter, `$`, or `_`" + (code3 === 43 || code3 !== null && code3 > 46 && code3 < 58 ? " (note: to create a link in MDX, use `[text](url)`)" : ""));
    }
    function localName(code3) {
      if (code3 !== null && code3 >= 0 && cont(code3, {
        jsx: true
      })) {
        effects.consume(code3);
        return localName;
      }
      if (code3 === 47 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        effects.exit(tagNameLocalType);
        returnState = localNameAfter;
        return esWhitespaceStart(code3);
      }
      crash(code3, "in local name", "a name character such as letters, digits, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function localNameAfter(code3) {
      if (code3 === 47 || code3 === 62 || code3 === 123 || code3 !== null && code3 >= 0 && start(code3)) {
        effects.exit(tagNameType);
        return attributeBefore(code3);
      }
      crash(code3, "after local name", "a character that can start an attribute name, such as a letter, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function attributeBefore(code3) {
      if (code3 === 47) {
        effects.enter(tagSelfClosingMarker);
        effects.consume(code3);
        effects.exit(tagSelfClosingMarker);
        returnState = selfClosing;
        return esWhitespaceStart;
      }
      if (code3 === 62) {
        return tagEnd(code3);
      }
      if (code3 === 123) {
        return factoryMdxExpression.call(self2, effects, attributeExpressionAfter, tagExpressionAttributeType, tagExpressionAttributeMarkerType, tagExpressionAttributeValueType, acorn, acornOptions, addResult, true, false, allowLazy)(code3);
      }
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagAttributeType);
        effects.enter(tagAttributeNameType);
        effects.enter(tagAttributeNamePrimaryType);
        effects.consume(code3);
        return attributePrimaryName;
      }
      crash(code3, "before attribute name", "a character that can start an attribute name, such as a letter, `$`, or `_`; whitespace before attributes; or the end of the tag");
    }
    function attributeExpressionAfter(code3) {
      returnState = attributeBefore;
      return esWhitespaceStart(code3);
    }
    function attributePrimaryName(code3) {
      if (code3 !== null && code3 >= 0 && cont(code3, {
        jsx: true
      })) {
        effects.consume(code3);
        return attributePrimaryName;
      }
      if (code3 === 47 || code3 === 58 || code3 === 61 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        effects.exit(tagAttributeNamePrimaryType);
        returnState = attributePrimaryNameAfter;
        return esWhitespaceStart(code3);
      }
      crash(code3, "in attribute name", "an attribute name character such as letters, digits, `$`, or `_`; `=` to initialize a value; whitespace before attributes; or the end of the tag");
    }
    function attributePrimaryNameAfter(code3) {
      if (code3 === 58) {
        effects.enter(tagAttributeNamePrefixMarkerType);
        effects.consume(code3);
        effects.exit(tagAttributeNamePrefixMarkerType);
        returnState = attributeLocalNameBefore;
        return esWhitespaceStart;
      }
      if (code3 === 61) {
        effects.exit(tagAttributeNameType);
        effects.enter(tagAttributeInitializerMarkerType);
        effects.consume(code3);
        effects.exit(tagAttributeInitializerMarkerType);
        returnState = attributeValueBefore;
        return esWhitespaceStart;
      }
      if (code3 === 47 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3) || code3 !== null && code3 >= 0 && start(code3)) {
        effects.exit(tagAttributeNameType);
        effects.exit(tagAttributeType);
        returnState = attributeBefore;
        return esWhitespaceStart(code3);
      }
      crash(code3, "after attribute name", "a character that can start an attribute name, such as a letter, `$`, or `_`; `=` to initialize a value; or the end of the tag");
    }
    function attributeLocalNameBefore(code3) {
      if (code3 !== null && code3 >= 0 && start(code3)) {
        effects.enter(tagAttributeNameLocalType);
        effects.consume(code3);
        return attributeLocalName;
      }
      crash(code3, "before local attribute name", "a character that can start an attribute name, such as a letter, `$`, or `_`; `=` to initialize a value; or the end of the tag");
    }
    function attributeLocalName(code3) {
      if (code3 !== null && code3 >= 0 && cont(code3, {
        jsx: true
      })) {
        effects.consume(code3);
        return attributeLocalName;
      }
      if (code3 === 47 || code3 === 61 || code3 === 62 || code3 === 123 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        effects.exit(tagAttributeNameLocalType);
        effects.exit(tagAttributeNameType);
        returnState = attributeLocalNameAfter;
        return esWhitespaceStart(code3);
      }
      crash(code3, "in local attribute name", "an attribute name character such as letters, digits, `$`, or `_`; `=` to initialize a value; whitespace before attributes; or the end of the tag");
    }
    function attributeLocalNameAfter(code3) {
      if (code3 === 61) {
        effects.enter(tagAttributeInitializerMarkerType);
        effects.consume(code3);
        effects.exit(tagAttributeInitializerMarkerType);
        returnState = attributeValueBefore;
        return esWhitespaceStart;
      }
      if (code3 === 47 || code3 === 62 || code3 === 123 || code3 !== null && code3 >= 0 && start(code3)) {
        effects.exit(tagAttributeType);
        return attributeBefore(code3);
      }
      crash(code3, "after local attribute name", "a character that can start an attribute name, such as a letter, `$`, or `_`; `=` to initialize a value; or the end of the tag");
    }
    function attributeValueBefore(code3) {
      if (code3 === 34 || code3 === 39) {
        effects.enter(tagAttributeValueLiteralType);
        effects.enter(tagAttributeValueLiteralMarkerType);
        effects.consume(code3);
        effects.exit(tagAttributeValueLiteralMarkerType);
        marker = code3;
        return attributeValueQuotedStart;
      }
      if (code3 === 123) {
        return factoryMdxExpression.call(self2, effects, attributeValueExpressionAfter, tagAttributeValueExpressionType, tagAttributeValueExpressionMarkerType, tagAttributeValueExpressionValueType, acorn, acornOptions, addResult, false, false, allowLazy)(code3);
      }
      crash(code3, "before attribute value", "a character that can start an attribute value, such as `\"`, `'`, or `{`" + (code3 === 60 ? " (note: to use an element or fragment as a prop value in MDX, use `{<element />}`)" : ""));
    }
    function attributeValueExpressionAfter(code3) {
      effects.exit(tagAttributeType);
      returnState = attributeBefore;
      return esWhitespaceStart(code3);
    }
    function attributeValueQuotedStart(code3) {
      if (code3 === null) {
        crash(code3, "in attribute value", "a corresponding closing quote `" + String.fromCodePoint(marker) + "`");
      }
      if (code3 === marker) {
        effects.enter(tagAttributeValueLiteralMarkerType);
        effects.consume(code3);
        effects.exit(tagAttributeValueLiteralMarkerType);
        effects.exit(tagAttributeValueLiteralType);
        effects.exit(tagAttributeType);
        marker = void 0;
        returnState = attributeBefore;
        return esWhitespaceStart;
      }
      if (markdownLineEnding(code3)) {
        returnState = attributeValueQuotedStart;
        return esWhitespaceStart(code3);
      }
      effects.enter(tagAttributeValueLiteralValueType);
      return attributeValueQuoted(code3);
    }
    function attributeValueQuoted(code3) {
      if (code3 === null || code3 === marker || markdownLineEnding(code3)) {
        effects.exit(tagAttributeValueLiteralValueType);
        return attributeValueQuotedStart(code3);
      }
      effects.consume(code3);
      return attributeValueQuoted;
    }
    function selfClosing(code3) {
      if (code3 === 62) {
        return tagEnd(code3);
      }
      crash(code3, "after self-closing slash", "`>` to end the tag" + (code3 === 42 || code3 === 47 ? " (note: JS comments in JSX tags are not supported in MDX)" : ""));
    }
    function tagEnd(code3) {
      effects.enter(tagMarkerType);
      effects.consume(code3);
      effects.exit(tagMarkerType);
      effects.exit(tagType);
      return ok3;
    }
    function esWhitespaceStart(code3) {
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return esWhitespaceEolAfter;
      }
      if (markdownSpace(code3) || unicodeWhitespace(code3)) {
        effects.enter("esWhitespace");
        return esWhitespaceInside(code3);
      }
      return returnState(code3);
    }
    function esWhitespaceInside(code3) {
      if (markdownLineEnding(code3)) {
        effects.exit("esWhitespace");
        return esWhitespaceStart(code3);
      }
      if (markdownSpace(code3) || unicodeWhitespace(code3)) {
        effects.consume(code3);
        return esWhitespaceInside;
      }
      effects.exit("esWhitespace");
      return returnState(code3);
    }
    function esWhitespaceEolAfter(code3) {
      if (!allowLazy && self2.parser.lazy[self2.now().line]) {
        const error = new VFileMessage("Unexpected lazy line in container, expected line to be prefixed with `>` when in a block quote, whitespace when in a list, etc", self2.now(), "micromark-extension-mdx-jsx:unexpected-lazy");
        error.url = trouble2 + "#unexpected-lazy-line-in-container-expected-line-to-be";
        throw error;
      }
      return esWhitespaceStart(code3);
    }
    function crash(code3, at2, expect) {
      const error = new VFileMessage("Unexpected " + (code3 === null ? "end of file" : "character `" + (code3 === 96 ? "` ` `" : String.fromCodePoint(code3)) + "` (" + serializeCharCode(code3) + ")") + " " + at2 + ", expected " + expect, self2.now(), "micromark-extension-mdx-jsx:unexpected-" + (code3 === null ? "eof" : "character"));
      error.url = trouble2 + (code3 === null ? "#unexpected-end-of-file-at-expected-expect" : "#unexpected-character-at-expected-expect");
      throw error;
    }
  }
  function serializeCharCode(code3) {
    return "U+" + code3.toString(16).toUpperCase().padStart(4, "0");
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-jsx@3.0.2/node_modules/micromark-extension-mdx-jsx/lib/jsx-text.js
  function jsxText(acorn, options) {
    return {
      name: "mdxJsxTextTag",
      tokenize: tokenizeJsxText
    };
    function tokenizeJsxText(effects, ok3, nok) {
      return factoryTag.call(this, effects, ok3, nok, acorn, options.acornOptions, options.addResult, true, "mdxJsxTextTag", "mdxJsxTextTagMarker", "mdxJsxTextTagClosingMarker", "mdxJsxTextTagSelfClosingMarker", "mdxJsxTextTagName", "mdxJsxTextTagNamePrimary", "mdxJsxTextTagNameMemberMarker", "mdxJsxTextTagNameMember", "mdxJsxTextTagNamePrefixMarker", "mdxJsxTextTagNameLocal", "mdxJsxTextTagExpressionAttribute", "mdxJsxTextTagExpressionAttributeMarker", "mdxJsxTextTagExpressionAttributeValue", "mdxJsxTextTagAttribute", "mdxJsxTextTagAttributeName", "mdxJsxTextTagAttributeNamePrimary", "mdxJsxTextTagAttributeNamePrefixMarker", "mdxJsxTextTagAttributeNameLocal", "mdxJsxTextTagAttributeInitializerMarker", "mdxJsxTextTagAttributeValueLiteral", "mdxJsxTextTagAttributeValueLiteralMarker", "mdxJsxTextTagAttributeValueLiteralValue", "mdxJsxTextTagAttributeValueExpression", "mdxJsxTextTagAttributeValueExpressionMarker", "mdxJsxTextTagAttributeValueExpressionValue");
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-jsx@3.0.2/node_modules/micromark-extension-mdx-jsx/lib/jsx-flow.js
  function jsxFlow(acorn, options) {
    return {
      concrete: true,
      name: "mdxJsxFlowTag",
      tokenize: tokenizeJsxFlow
    };
    function tokenizeJsxFlow(effects, ok3, nok) {
      const self2 = this;
      return start2;
      function start2(code3) {
        return before(code3);
      }
      function before(code3) {
        return factoryTag.call(self2, effects, after, nok, acorn, options.acornOptions, options.addResult, false, "mdxJsxFlowTag", "mdxJsxFlowTagMarker", "mdxJsxFlowTagClosingMarker", "mdxJsxFlowTagSelfClosingMarker", "mdxJsxFlowTagName", "mdxJsxFlowTagNamePrimary", "mdxJsxFlowTagNameMemberMarker", "mdxJsxFlowTagNameMember", "mdxJsxFlowTagNamePrefixMarker", "mdxJsxFlowTagNameLocal", "mdxJsxFlowTagExpressionAttribute", "mdxJsxFlowTagExpressionAttributeMarker", "mdxJsxFlowTagExpressionAttributeValue", "mdxJsxFlowTagAttribute", "mdxJsxFlowTagAttributeName", "mdxJsxFlowTagAttributeNamePrimary", "mdxJsxFlowTagAttributeNamePrefixMarker", "mdxJsxFlowTagAttributeNameLocal", "mdxJsxFlowTagAttributeInitializerMarker", "mdxJsxFlowTagAttributeValueLiteral", "mdxJsxFlowTagAttributeValueLiteralMarker", "mdxJsxFlowTagAttributeValueLiteralValue", "mdxJsxFlowTagAttributeValueExpression", "mdxJsxFlowTagAttributeValueExpressionMarker", "mdxJsxFlowTagAttributeValueExpressionValue")(code3);
      }
      function after(code3) {
        return markdownSpace(code3) ? factorySpace(effects, end, "whitespace")(code3) : end(code3);
      }
      function end(code3) {
        const leftBraceValue = self2.parser.constructs.flow[123];
        const constructs2 = Array.isArray(leftBraceValue) ? leftBraceValue : leftBraceValue ? [leftBraceValue] : [];
        let expression;
        for (const construct of constructs2) {
          if (construct.name === "mdxFlowExpression") {
            expression = construct;
            break;
          }
        }
        return code3 === 60 ? (
          // We can’t just say: fine. Lines of blocks have to be parsed until an eol/eof.
          start2(code3)
        ) : code3 === 123 && expression ? effects.attempt(expression, end, nok)(code3) : code3 === null || markdownLineEnding(code3) ? ok3(code3) : nok(code3);
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-jsx@3.0.2/node_modules/micromark-extension-mdx-jsx/lib/syntax.js
  function mdxJsx(options) {
    const settings = options || {};
    const acorn = settings.acorn;
    let acornOptions;
    if (acorn) {
      if (!acorn.parse || !acorn.parseExpressionAt) {
        throw new Error("Expected a proper `acorn` instance passed in as `options.acorn`");
      }
      acornOptions = Object.assign({
        ecmaVersion: 2024,
        sourceType: "module"
      }, settings.acornOptions, {
        locations: true
      });
    } else if (settings.acornOptions || settings.addResult) {
      throw new Error("Expected an `acorn` instance passed in as `options.acorn`");
    }
    return {
      flow: {
        [60]: jsxFlow(acorn || void 0, {
          acornOptions,
          addResult: settings.addResult || void 0
        })
      },
      text: {
        [60]: jsxText(acorn || void 0, {
          acornOptions,
          addResult: settings.addResult || void 0
        })
      }
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdx-md@2.0.0/node_modules/micromark-extension-mdx-md/index.js
  function mdxMd() {
    return {
      disable: { null: ["autolink", "codeIndented", "htmlFlow", "htmlText"] }
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-chunked@2.0.1/node_modules/micromark-util-chunked/index.js
  function splice(list4, start2, remove5, items) {
    const end = list4.length;
    let chunkStart = 0;
    let parameters;
    if (start2 < 0) {
      start2 = -start2 > end ? 0 : end + start2;
    } else {
      start2 = start2 > end ? end : start2;
    }
    remove5 = remove5 > 0 ? remove5 : 0;
    if (items.length < 1e4) {
      parameters = Array.from(items);
      parameters.unshift(start2, remove5);
      list4.splice(...parameters);
    } else {
      if (remove5) list4.splice(start2, remove5);
      while (chunkStart < items.length) {
        parameters = items.slice(chunkStart, chunkStart + 1e4);
        parameters.unshift(start2, 0);
        list4.splice(...parameters);
        chunkStart += 1e4;
        start2 += 1e4;
      }
    }
  }
  function push(list4, items) {
    if (list4.length > 0) {
      splice(list4, list4.length, 0, items);
      return list4;
    }
    return items;
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-classify-character@2.0.1/node_modules/micromark-util-classify-character/index.js
  function classifyCharacter(code3) {
    if (code3 === null || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
      return 1;
    }
    if (unicodePunctuation(code3)) {
      return 2;
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-resolve-all@2.0.1/node_modules/micromark-util-resolve-all/index.js
  function resolveAll(constructs2, events, context) {
    const called = [];
    let index2 = -1;
    while (++index2 < constructs2.length) {
      const resolve = constructs2[index2].resolveAll;
      if (resolve && !called.includes(resolve)) {
        events = resolve(events, context);
        called.push(resolve);
      }
    }
    return events;
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/attention.js
  var attention = {
    name: "attention",
    resolveAll: resolveAllAttention,
    tokenize: tokenizeAttention
  };
  function resolveAllAttention(events, context) {
    let index2 = -1;
    let open;
    let group;
    let text5;
    let openingSequence;
    let closingSequence;
    let use;
    let nextEvents;
    let offset2;
    while (++index2 < events.length) {
      if (events[index2][0] === "enter" && events[index2][1].type === "attentionSequence" && events[index2][1]._close) {
        open = index2;
        while (open--) {
          if (events[open][0] === "exit" && events[open][1].type === "attentionSequence" && events[open][1]._open && // If the markers are the same:
          context.sliceSerialize(events[open][1]).charCodeAt(0) === context.sliceSerialize(events[index2][1]).charCodeAt(0)) {
            if ((events[open][1]._close || events[index2][1]._open) && (events[index2][1].end.offset - events[index2][1].start.offset) % 3 && !((events[open][1].end.offset - events[open][1].start.offset + events[index2][1].end.offset - events[index2][1].start.offset) % 3)) {
              continue;
            }
            use = events[open][1].end.offset - events[open][1].start.offset > 1 && events[index2][1].end.offset - events[index2][1].start.offset > 1 ? 2 : 1;
            const start2 = {
              ...events[open][1].end
            };
            const end = {
              ...events[index2][1].start
            };
            movePoint(start2, -use);
            movePoint(end, use);
            openingSequence = {
              type: use > 1 ? "strongSequence" : "emphasisSequence",
              start: start2,
              end: {
                ...events[open][1].end
              }
            };
            closingSequence = {
              type: use > 1 ? "strongSequence" : "emphasisSequence",
              start: {
                ...events[index2][1].start
              },
              end
            };
            text5 = {
              type: use > 1 ? "strongText" : "emphasisText",
              start: {
                ...events[open][1].end
              },
              end: {
                ...events[index2][1].start
              }
            };
            group = {
              type: use > 1 ? "strong" : "emphasis",
              start: {
                ...openingSequence.start
              },
              end: {
                ...closingSequence.end
              }
            };
            events[open][1].end = {
              ...openingSequence.start
            };
            events[index2][1].start = {
              ...closingSequence.end
            };
            nextEvents = [];
            if (events[open][1].end.offset - events[open][1].start.offset) {
              nextEvents = push(nextEvents, [["enter", events[open][1], context], ["exit", events[open][1], context]]);
            }
            nextEvents = push(nextEvents, [["enter", group, context], ["enter", openingSequence, context], ["exit", openingSequence, context], ["enter", text5, context]]);
            nextEvents = push(nextEvents, resolveAll(context.parser.constructs.insideSpan.null, events.slice(open + 1, index2), context));
            nextEvents = push(nextEvents, [["exit", text5, context], ["enter", closingSequence, context], ["exit", closingSequence, context], ["exit", group, context]]);
            if (events[index2][1].end.offset - events[index2][1].start.offset) {
              offset2 = 2;
              nextEvents = push(nextEvents, [["enter", events[index2][1], context], ["exit", events[index2][1], context]]);
            } else {
              offset2 = 0;
            }
            splice(events, open - 1, index2 - open + 3, nextEvents);
            index2 = open + nextEvents.length - offset2 - 2;
            break;
          }
        }
      }
    }
    index2 = -1;
    while (++index2 < events.length) {
      if (events[index2][1].type === "attentionSequence") {
        events[index2][1].type = "data";
      }
    }
    return events;
  }
  function tokenizeAttention(effects, ok3) {
    const attentionMarkers2 = this.parser.constructs.attentionMarkers.null;
    const previous3 = this.previous;
    const before = classifyCharacter(previous3);
    let marker;
    return start2;
    function start2(code3) {
      marker = code3;
      effects.enter("attentionSequence");
      return inside(code3);
    }
    function inside(code3) {
      if (code3 === marker) {
        effects.consume(code3);
        return inside;
      }
      const token = effects.exit("attentionSequence");
      const after = classifyCharacter(code3);
      const open = !after || after === 2 && before || attentionMarkers2.includes(code3);
      const close = !before || before === 2 && after || attentionMarkers2.includes(previous3);
      token._open = Boolean(marker === 42 ? open : open && (before || !close));
      token._close = Boolean(marker === 42 ? close : close && (after || !open));
      return ok3(code3);
    }
  }
  function movePoint(point3, offset2) {
    point3.column += offset2;
    point3.offset += offset2;
    point3._bufferIndex += offset2;
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/autolink.js
  var autolink = {
    name: "autolink",
    tokenize: tokenizeAutolink
  };
  function tokenizeAutolink(effects, ok3, nok) {
    let size = 0;
    return start2;
    function start2(code3) {
      effects.enter("autolink");
      effects.enter("autolinkMarker");
      effects.consume(code3);
      effects.exit("autolinkMarker");
      effects.enter("autolinkProtocol");
      return open;
    }
    function open(code3) {
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        return schemeOrEmailAtext;
      }
      if (code3 === 64) {
        return nok(code3);
      }
      return emailAtext(code3);
    }
    function schemeOrEmailAtext(code3) {
      if (code3 === 43 || code3 === 45 || code3 === 46 || asciiAlphanumeric(code3)) {
        size = 1;
        return schemeInsideOrEmailAtext(code3);
      }
      return emailAtext(code3);
    }
    function schemeInsideOrEmailAtext(code3) {
      if (code3 === 58) {
        effects.consume(code3);
        size = 0;
        return urlInside;
      }
      if ((code3 === 43 || code3 === 45 || code3 === 46 || asciiAlphanumeric(code3)) && size++ < 32) {
        effects.consume(code3);
        return schemeInsideOrEmailAtext;
      }
      size = 0;
      return emailAtext(code3);
    }
    function urlInside(code3) {
      if (code3 === 62) {
        effects.exit("autolinkProtocol");
        effects.enter("autolinkMarker");
        effects.consume(code3);
        effects.exit("autolinkMarker");
        effects.exit("autolink");
        return ok3;
      }
      if (code3 === null || code3 === 32 || code3 === 60 || asciiControl(code3)) {
        return nok(code3);
      }
      effects.consume(code3);
      return urlInside;
    }
    function emailAtext(code3) {
      if (code3 === 64) {
        effects.consume(code3);
        return emailAtSignOrDot;
      }
      if (asciiAtext(code3)) {
        effects.consume(code3);
        return emailAtext;
      }
      return nok(code3);
    }
    function emailAtSignOrDot(code3) {
      return asciiAlphanumeric(code3) ? emailLabel(code3) : nok(code3);
    }
    function emailLabel(code3) {
      if (code3 === 46) {
        effects.consume(code3);
        size = 0;
        return emailAtSignOrDot;
      }
      if (code3 === 62) {
        effects.exit("autolinkProtocol").type = "autolinkEmail";
        effects.enter("autolinkMarker");
        effects.consume(code3);
        effects.exit("autolinkMarker");
        effects.exit("autolink");
        return ok3;
      }
      return emailValue(code3);
    }
    function emailValue(code3) {
      if ((code3 === 45 || asciiAlphanumeric(code3)) && size++ < 63) {
        const next = code3 === 45 ? emailValue : emailLabel;
        effects.consume(code3);
        return next;
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/blank-line.js
  var blankLine = {
    partial: true,
    tokenize: tokenizeBlankLine
  };
  function tokenizeBlankLine(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      return markdownSpace(code3) ? factorySpace(effects, after, "linePrefix")(code3) : after(code3);
    }
    function after(code3) {
      return code3 === null || markdownLineEnding(code3) ? ok3(code3) : nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/block-quote.js
  var blockQuote = {
    continuation: {
      tokenize: tokenizeBlockQuoteContinuation
    },
    exit,
    name: "blockQuote",
    tokenize: tokenizeBlockQuoteStart
  };
  function tokenizeBlockQuoteStart(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      if (code3 === 62) {
        const state = self2.containerState;
        if (!state.open) {
          effects.enter("blockQuote", {
            _container: true
          });
          state.open = true;
        }
        effects.enter("blockQuotePrefix");
        effects.enter("blockQuoteMarker");
        effects.consume(code3);
        effects.exit("blockQuoteMarker");
        return after;
      }
      return nok(code3);
    }
    function after(code3) {
      if (markdownSpace(code3)) {
        effects.enter("blockQuotePrefixWhitespace");
        effects.consume(code3);
        effects.exit("blockQuotePrefixWhitespace");
        effects.exit("blockQuotePrefix");
        return ok3;
      }
      effects.exit("blockQuotePrefix");
      return ok3(code3);
    }
  }
  function tokenizeBlockQuoteContinuation(effects, ok3, nok) {
    const self2 = this;
    return contStart;
    function contStart(code3) {
      if (markdownSpace(code3)) {
        return factorySpace(effects, contBefore, "linePrefix", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4)(code3);
      }
      return contBefore(code3);
    }
    function contBefore(code3) {
      return effects.attempt(blockQuote, ok3, nok)(code3);
    }
  }
  function exit(effects) {
    effects.exit("blockQuote");
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/character-escape.js
  var characterEscape = {
    name: "characterEscape",
    tokenize: tokenizeCharacterEscape
  };
  function tokenizeCharacterEscape(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      effects.enter("characterEscape");
      effects.enter("escapeMarker");
      effects.consume(code3);
      effects.exit("escapeMarker");
      return inside;
    }
    function inside(code3) {
      if (asciiPunctuation(code3)) {
        effects.enter("characterEscapeValue");
        effects.consume(code3);
        effects.exit("characterEscapeValue");
        effects.exit("characterEscape");
        return ok3;
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/character-reference.js
  var characterReference = {
    name: "characterReference",
    tokenize: tokenizeCharacterReference
  };
  function tokenizeCharacterReference(effects, ok3, nok) {
    const self2 = this;
    let size = 0;
    let max;
    let test;
    return start2;
    function start2(code3) {
      effects.enter("characterReference");
      effects.enter("characterReferenceMarker");
      effects.consume(code3);
      effects.exit("characterReferenceMarker");
      return open;
    }
    function open(code3) {
      if (code3 === 35) {
        effects.enter("characterReferenceMarkerNumeric");
        effects.consume(code3);
        effects.exit("characterReferenceMarkerNumeric");
        return numeric;
      }
      effects.enter("characterReferenceValue");
      max = 31;
      test = asciiAlphanumeric;
      return value(code3);
    }
    function numeric(code3) {
      if (code3 === 88 || code3 === 120) {
        effects.enter("characterReferenceMarkerHexadecimal");
        effects.consume(code3);
        effects.exit("characterReferenceMarkerHexadecimal");
        effects.enter("characterReferenceValue");
        max = 6;
        test = asciiHexDigit;
        return value;
      }
      effects.enter("characterReferenceValue");
      max = 7;
      test = asciiDigit;
      return value(code3);
    }
    function value(code3) {
      if (code3 === 59 && size) {
        const token = effects.exit("characterReferenceValue");
        if (test === asciiAlphanumeric && !decodeNamedCharacterReference(self2.sliceSerialize(token))) {
          return nok(code3);
        }
        effects.enter("characterReferenceMarker");
        effects.consume(code3);
        effects.exit("characterReferenceMarker");
        effects.exit("characterReference");
        return ok3;
      }
      if (test(code3) && size++ < max) {
        effects.consume(code3);
        return value;
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/code-fenced.js
  var nonLazyContinuation = {
    partial: true,
    tokenize: tokenizeNonLazyContinuation
  };
  var codeFenced = {
    concrete: true,
    name: "codeFenced",
    tokenize: tokenizeCodeFenced
  };
  function tokenizeCodeFenced(effects, ok3, nok) {
    const self2 = this;
    const closeStart = {
      partial: true,
      tokenize: tokenizeCloseStart
    };
    let initialPrefix = 0;
    let sizeOpen = 0;
    let marker;
    return start2;
    function start2(code3) {
      return beforeSequenceOpen(code3);
    }
    function beforeSequenceOpen(code3) {
      const tail = self2.events[self2.events.length - 1];
      initialPrefix = tail && tail[1].type === "linePrefix" ? tail[2].sliceSerialize(tail[1], true).length : 0;
      marker = code3;
      effects.enter("codeFenced");
      effects.enter("codeFencedFence");
      effects.enter("codeFencedFenceSequence");
      return sequenceOpen(code3);
    }
    function sequenceOpen(code3) {
      if (code3 === marker) {
        sizeOpen++;
        effects.consume(code3);
        return sequenceOpen;
      }
      if (sizeOpen < 3) {
        return nok(code3);
      }
      effects.exit("codeFencedFenceSequence");
      return markdownSpace(code3) ? factorySpace(effects, infoBefore, "whitespace")(code3) : infoBefore(code3);
    }
    function infoBefore(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("codeFencedFence");
        return self2.interrupt ? ok3(code3) : effects.check(nonLazyContinuation, atNonLazyBreak, after)(code3);
      }
      effects.enter("codeFencedFenceInfo");
      effects.enter("chunkString", {
        contentType: "string"
      });
      return info(code3);
    }
    function info(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("chunkString");
        effects.exit("codeFencedFenceInfo");
        return infoBefore(code3);
      }
      if (markdownSpace(code3)) {
        effects.exit("chunkString");
        effects.exit("codeFencedFenceInfo");
        return factorySpace(effects, metaBefore, "whitespace")(code3);
      }
      if (code3 === 96 && code3 === marker) {
        return nok(code3);
      }
      effects.consume(code3);
      return info;
    }
    function metaBefore(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        return infoBefore(code3);
      }
      effects.enter("codeFencedFenceMeta");
      effects.enter("chunkString", {
        contentType: "string"
      });
      return meta(code3);
    }
    function meta(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("chunkString");
        effects.exit("codeFencedFenceMeta");
        return infoBefore(code3);
      }
      if (code3 === 96 && code3 === marker) {
        return nok(code3);
      }
      effects.consume(code3);
      return meta;
    }
    function atNonLazyBreak(code3) {
      return effects.attempt(closeStart, after, contentBefore)(code3);
    }
    function contentBefore(code3) {
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return contentStart;
    }
    function contentStart(code3) {
      return initialPrefix > 0 && markdownSpace(code3) ? factorySpace(effects, beforeContentChunk, "linePrefix", initialPrefix + 1)(code3) : beforeContentChunk(code3);
    }
    function beforeContentChunk(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        return effects.check(nonLazyContinuation, atNonLazyBreak, after)(code3);
      }
      effects.enter("codeFlowValue");
      return contentChunk(code3);
    }
    function contentChunk(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("codeFlowValue");
        return beforeContentChunk(code3);
      }
      effects.consume(code3);
      return contentChunk;
    }
    function after(code3) {
      effects.exit("codeFenced");
      return ok3(code3);
    }
    function tokenizeCloseStart(effects2, ok4, nok2) {
      let size = 0;
      return startBefore;
      function startBefore(code3) {
        effects2.enter("lineEnding");
        effects2.consume(code3);
        effects2.exit("lineEnding");
        return start3;
      }
      function start3(code3) {
        effects2.enter("codeFencedFence");
        return markdownSpace(code3) ? factorySpace(effects2, beforeSequenceClose, "linePrefix", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4)(code3) : beforeSequenceClose(code3);
      }
      function beforeSequenceClose(code3) {
        if (code3 === marker) {
          effects2.enter("codeFencedFenceSequence");
          return sequenceClose(code3);
        }
        return nok2(code3);
      }
      function sequenceClose(code3) {
        if (code3 === marker) {
          size++;
          effects2.consume(code3);
          return sequenceClose;
        }
        if (size >= sizeOpen) {
          effects2.exit("codeFencedFenceSequence");
          return markdownSpace(code3) ? factorySpace(effects2, sequenceCloseAfter, "whitespace")(code3) : sequenceCloseAfter(code3);
        }
        return nok2(code3);
      }
      function sequenceCloseAfter(code3) {
        if (code3 === null || markdownLineEnding(code3)) {
          effects2.exit("codeFencedFence");
          return ok4(code3);
        }
        return nok2(code3);
      }
    }
  }
  function tokenizeNonLazyContinuation(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return lineStart;
    }
    function lineStart(code3) {
      return self2.parser.lazy[self2.now().line] ? nok(code3) : ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/code-indented.js
  var codeIndented = {
    name: "codeIndented",
    tokenize: tokenizeCodeIndented
  };
  var furtherStart = {
    partial: true,
    tokenize: tokenizeFurtherStart
  };
  function tokenizeCodeIndented(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      effects.enter("codeIndented");
      return factorySpace(effects, afterPrefix, "linePrefix", 4 + 1)(code3);
    }
    function afterPrefix(code3) {
      const tail = self2.events[self2.events.length - 1];
      return tail && tail[1].type === "linePrefix" && tail[2].sliceSerialize(tail[1], true).length >= 4 ? atBreak(code3) : nok(code3);
    }
    function atBreak(code3) {
      if (code3 === null) {
        return after(code3);
      }
      if (markdownLineEnding(code3)) {
        return effects.attempt(furtherStart, atBreak, after)(code3);
      }
      effects.enter("codeFlowValue");
      return inside(code3);
    }
    function inside(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("codeFlowValue");
        return atBreak(code3);
      }
      effects.consume(code3);
      return inside;
    }
    function after(code3) {
      effects.exit("codeIndented");
      return ok3(code3);
    }
  }
  function tokenizeFurtherStart(effects, ok3, nok) {
    const self2 = this;
    return furtherStart2;
    function furtherStart2(code3) {
      if (self2.parser.lazy[self2.now().line]) {
        return nok(code3);
      }
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return furtherStart2;
      }
      return factorySpace(effects, afterPrefix, "linePrefix", 4 + 1)(code3);
    }
    function afterPrefix(code3) {
      const tail = self2.events[self2.events.length - 1];
      return tail && tail[1].type === "linePrefix" && tail[2].sliceSerialize(tail[1], true).length >= 4 ? ok3(code3) : markdownLineEnding(code3) ? furtherStart2(code3) : nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/code-text.js
  var codeText = {
    name: "codeText",
    previous,
    resolve: resolveCodeText,
    tokenize: tokenizeCodeText
  };
  function resolveCodeText(events) {
    let tailExitIndex = events.length - 4;
    let headEnterIndex = 3;
    let index2;
    let enter;
    if ((events[headEnterIndex][1].type === "lineEnding" || events[headEnterIndex][1].type === "space") && (events[tailExitIndex][1].type === "lineEnding" || events[tailExitIndex][1].type === "space")) {
      index2 = headEnterIndex;
      while (++index2 < tailExitIndex) {
        if (events[index2][1].type === "codeTextData") {
          events[headEnterIndex][1].type = "codeTextPadding";
          events[tailExitIndex][1].type = "codeTextPadding";
          headEnterIndex += 2;
          tailExitIndex -= 2;
          break;
        }
      }
    }
    index2 = headEnterIndex - 1;
    tailExitIndex++;
    while (++index2 <= tailExitIndex) {
      if (enter === void 0) {
        if (index2 !== tailExitIndex && events[index2][1].type !== "lineEnding") {
          enter = index2;
        }
      } else if (index2 === tailExitIndex || events[index2][1].type === "lineEnding") {
        events[enter][1].type = "codeTextData";
        if (index2 !== enter + 2) {
          events[enter][1].end = events[index2 - 1][1].end;
          events.splice(enter + 2, index2 - enter - 2);
          tailExitIndex -= index2 - enter - 2;
          index2 = enter + 2;
        }
        enter = void 0;
      }
    }
    return events;
  }
  function previous(code3) {
    return code3 !== 96 || this.events[this.events.length - 1][1].type === "characterEscape";
  }
  function tokenizeCodeText(effects, ok3, nok) {
    const self2 = this;
    let sizeOpen = 0;
    let size;
    let token;
    return start2;
    function start2(code3) {
      effects.enter("codeText");
      effects.enter("codeTextSequence");
      return sequenceOpen(code3);
    }
    function sequenceOpen(code3) {
      if (code3 === 96) {
        effects.consume(code3);
        sizeOpen++;
        return sequenceOpen;
      }
      effects.exit("codeTextSequence");
      return between2(code3);
    }
    function between2(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      if (code3 === 32) {
        effects.enter("space");
        effects.consume(code3);
        effects.exit("space");
        return between2;
      }
      if (code3 === 96) {
        token = effects.enter("codeTextSequence");
        size = 0;
        return sequenceClose(code3);
      }
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return between2;
      }
      effects.enter("codeTextData");
      return data2(code3);
    }
    function data2(code3) {
      if (code3 === null || code3 === 32 || code3 === 96 || markdownLineEnding(code3)) {
        effects.exit("codeTextData");
        return between2(code3);
      }
      effects.consume(code3);
      return data2;
    }
    function sequenceClose(code3) {
      if (code3 === 96) {
        effects.consume(code3);
        size++;
        return sequenceClose;
      }
      if (size === sizeOpen) {
        effects.exit("codeTextSequence");
        effects.exit("codeText");
        return ok3(code3);
      }
      token.type = "codeTextData";
      return data2(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-subtokenize@2.1.0/node_modules/micromark-util-subtokenize/lib/splice-buffer.js
  var SpliceBuffer = class {
    /**
     * @param {ReadonlyArray<T> | null | undefined} [initial]
     *   Initial items (optional).
     * @returns
     *   Splice buffer.
     */
    constructor(initial) {
      this.left = initial ? [...initial] : [];
      this.right = [];
    }
    /**
     * Array access;
     * does not move the cursor.
     *
     * @param {number} index
     *   Index.
     * @return {T}
     *   Item.
     */
    get(index2) {
      if (index2 < 0 || index2 >= this.left.length + this.right.length) {
        throw new RangeError("Cannot access index `" + index2 + "` in a splice buffer of size `" + (this.left.length + this.right.length) + "`");
      }
      if (index2 < this.left.length) return this.left[index2];
      return this.right[this.right.length - index2 + this.left.length - 1];
    }
    /**
     * The length of the splice buffer, one greater than the largest index in the
     * array.
     */
    get length() {
      return this.left.length + this.right.length;
    }
    /**
     * Remove and return `list[0]`;
     * moves the cursor to `0`.
     *
     * @returns {T | undefined}
     *   Item, optional.
     */
    shift() {
      this.setCursor(0);
      return this.right.pop();
    }
    /**
     * Slice the buffer to get an array;
     * does not move the cursor.
     *
     * @param {number} start
     *   Start.
     * @param {number | null | undefined} [end]
     *   End (optional).
     * @returns {Array<T>}
     *   Array of items.
     */
    slice(start2, end) {
      const stop = end === null || end === void 0 ? Number.POSITIVE_INFINITY : end;
      if (stop < this.left.length) {
        return this.left.slice(start2, stop);
      }
      if (start2 > this.left.length) {
        return this.right.slice(this.right.length - stop + this.left.length, this.right.length - start2 + this.left.length).reverse();
      }
      return this.left.slice(start2).concat(this.right.slice(this.right.length - stop + this.left.length).reverse());
    }
    /**
     * Mimics the behavior of Array.prototype.splice() except for the change of
     * interface necessary to avoid segfaults when patching in very large arrays.
     *
     * This operation moves cursor is moved to `start` and results in the cursor
     * placed after any inserted items.
     *
     * @param {number} start
     *   Start;
     *   zero-based index at which to start changing the array;
     *   negative numbers count backwards from the end of the array and values
     *   that are out-of bounds are clamped to the appropriate end of the array.
     * @param {number | null | undefined} [deleteCount=0]
     *   Delete count (default: `0`);
     *   maximum number of elements to delete, starting from start.
     * @param {Array<T> | null | undefined} [items=[]]
     *   Items to include in place of the deleted items (default: `[]`).
     * @return {Array<T>}
     *   Any removed items.
     */
    splice(start2, deleteCount, items) {
      const count3 = deleteCount || 0;
      this.setCursor(Math.trunc(start2));
      const removed = this.right.splice(this.right.length - count3, Number.POSITIVE_INFINITY);
      if (items) chunkedPush(this.left, items);
      return removed.reverse();
    }
    /**
     * Remove and return the highest-numbered item in the array, so
     * `list[list.length - 1]`;
     * Moves the cursor to `length`.
     *
     * @returns {T | undefined}
     *   Item, optional.
     */
    pop() {
      this.setCursor(Number.POSITIVE_INFINITY);
      return this.left.pop();
    }
    /**
     * Inserts a single item to the high-numbered side of the array;
     * moves the cursor to `length`.
     *
     * @param {T} item
     *   Item.
     * @returns {undefined}
     *   Nothing.
     */
    push(item) {
      this.setCursor(Number.POSITIVE_INFINITY);
      this.left.push(item);
    }
    /**
     * Inserts many items to the high-numbered side of the array.
     * Moves the cursor to `length`.
     *
     * @param {Array<T>} items
     *   Items.
     * @returns {undefined}
     *   Nothing.
     */
    pushMany(items) {
      this.setCursor(Number.POSITIVE_INFINITY);
      chunkedPush(this.left, items);
    }
    /**
     * Inserts a single item to the low-numbered side of the array;
     * Moves the cursor to `0`.
     *
     * @param {T} item
     *   Item.
     * @returns {undefined}
     *   Nothing.
     */
    unshift(item) {
      this.setCursor(0);
      this.right.push(item);
    }
    /**
     * Inserts many items to the low-numbered side of the array;
     * moves the cursor to `0`.
     *
     * @param {Array<T>} items
     *   Items.
     * @returns {undefined}
     *   Nothing.
     */
    unshiftMany(items) {
      this.setCursor(0);
      chunkedPush(this.right, items.reverse());
    }
    /**
     * Move the cursor to a specific position in the array. Requires
     * time proportional to the distance moved.
     *
     * If `n < 0`, the cursor will end up at the beginning.
     * If `n > length`, the cursor will end up at the end.
     *
     * @param {number} n
     *   Position.
     * @return {undefined}
     *   Nothing.
     */
    setCursor(n) {
      if (n === this.left.length || n > this.left.length && this.right.length === 0 || n < 0 && this.left.length === 0) return;
      if (n < this.left.length) {
        const removed = this.left.splice(n, Number.POSITIVE_INFINITY);
        chunkedPush(this.right, removed.reverse());
      } else {
        const removed = this.right.splice(this.left.length + this.right.length - n, Number.POSITIVE_INFINITY);
        chunkedPush(this.left, removed.reverse());
      }
    }
  };
  function chunkedPush(list4, right) {
    let chunkStart = 0;
    if (right.length < 1e4) {
      list4.push(...right);
    } else {
      while (chunkStart < right.length) {
        list4.push(...right.slice(chunkStart, chunkStart + 1e4));
        chunkStart += 1e4;
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-subtokenize@2.1.0/node_modules/micromark-util-subtokenize/index.js
  function subtokenize(eventsArray) {
    const jumps = {};
    let index2 = -1;
    let event;
    let lineIndex;
    let otherIndex;
    let otherEvent;
    let parameters;
    let subevents;
    let more;
    const events = new SpliceBuffer(eventsArray);
    while (++index2 < events.length) {
      while (index2 in jumps) {
        index2 = jumps[index2];
      }
      event = events.get(index2);
      if (index2 && event[1].type === "chunkFlow" && events.get(index2 - 1)[1].type === "listItemPrefix") {
        subevents = event[1]._tokenizer.events;
        otherIndex = 0;
        if (otherIndex < subevents.length && subevents[otherIndex][1].type === "lineEndingBlank") {
          otherIndex += 2;
        }
        if (otherIndex < subevents.length && subevents[otherIndex][1].type === "content") {
          while (++otherIndex < subevents.length) {
            if (subevents[otherIndex][1].type === "content") {
              break;
            }
            if (subevents[otherIndex][1].type === "chunkText") {
              subevents[otherIndex][1]._isInFirstContentOfListItem = true;
              otherIndex++;
            }
          }
        }
      }
      if (event[0] === "enter") {
        if (event[1].contentType) {
          Object.assign(jumps, subcontent(events, index2));
          index2 = jumps[index2];
          more = true;
        }
      } else if (event[1]._container) {
        otherIndex = index2;
        lineIndex = void 0;
        while (otherIndex--) {
          otherEvent = events.get(otherIndex);
          if (otherEvent[1].type === "lineEnding" || otherEvent[1].type === "lineEndingBlank") {
            if (otherEvent[0] === "enter") {
              if (lineIndex) {
                events.get(lineIndex)[1].type = "lineEndingBlank";
              }
              otherEvent[1].type = "lineEnding";
              lineIndex = otherIndex;
            }
          } else if (otherEvent[1].type === "linePrefix" || otherEvent[1].type === "listItemIndent") {
          } else {
            break;
          }
        }
        if (lineIndex) {
          event[1].end = {
            ...events.get(lineIndex)[1].start
          };
          parameters = events.slice(lineIndex, index2);
          parameters.unshift(event);
          events.splice(lineIndex, index2 - lineIndex + 1, parameters);
        }
      }
    }
    splice(eventsArray, 0, Number.POSITIVE_INFINITY, events.slice(0));
    return !more;
  }
  function subcontent(events, eventIndex) {
    const token = events.get(eventIndex)[1];
    const context = events.get(eventIndex)[2];
    let startPosition = eventIndex - 1;
    const startPositions = [];
    let tokenizer3 = token._tokenizer;
    if (!tokenizer3) {
      tokenizer3 = context.parser[token.contentType](token.start);
      if (token._contentTypeTextTrailing) {
        tokenizer3._contentTypeTextTrailing = true;
      }
    }
    const childEvents = tokenizer3.events;
    const jumps = [];
    const gaps = {};
    let stream;
    let previous3;
    let index2 = -1;
    let current2 = token;
    let adjust = 0;
    let start2 = 0;
    const breaks = [start2];
    while (current2) {
      while (events.get(++startPosition)[1] !== current2) {
      }
      startPositions.push(startPosition);
      if (!current2._tokenizer) {
        stream = context.sliceStream(current2);
        if (!current2.next) {
          stream.push(null);
        }
        if (previous3) {
          tokenizer3.defineSkip(current2.start);
        }
        if (current2._isInFirstContentOfListItem) {
          tokenizer3._gfmTasklistFirstContentOfListItem = true;
        }
        tokenizer3.write(stream);
        if (current2._isInFirstContentOfListItem) {
          tokenizer3._gfmTasklistFirstContentOfListItem = void 0;
        }
      }
      previous3 = current2;
      current2 = current2.next;
    }
    current2 = token;
    while (++index2 < childEvents.length) {
      if (
        // Find a void token that includes a break.
        childEvents[index2][0] === "exit" && childEvents[index2 - 1][0] === "enter" && childEvents[index2][1].type === childEvents[index2 - 1][1].type && childEvents[index2][1].start.line !== childEvents[index2][1].end.line
      ) {
        start2 = index2 + 1;
        breaks.push(start2);
        current2._tokenizer = void 0;
        current2.previous = void 0;
        current2 = current2.next;
      }
    }
    tokenizer3.events = [];
    if (current2) {
      current2._tokenizer = void 0;
      current2.previous = void 0;
    } else {
      breaks.pop();
    }
    index2 = breaks.length;
    while (index2--) {
      const slice = childEvents.slice(breaks[index2], breaks[index2 + 1]);
      const start3 = startPositions.pop();
      jumps.push([start3, start3 + slice.length - 1]);
      events.splice(start3, 2, slice);
    }
    jumps.reverse();
    index2 = -1;
    while (++index2 < jumps.length) {
      gaps[adjust + jumps[index2][0]] = adjust + jumps[index2][1];
      adjust += jumps[index2][1] - jumps[index2][0] - 1;
    }
    return gaps;
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/content.js
  var content = {
    resolve: resolveContent,
    tokenize: tokenizeContent
  };
  var continuationConstruct = {
    partial: true,
    tokenize: tokenizeContinuation
  };
  function resolveContent(events) {
    subtokenize(events);
    return events;
  }
  function tokenizeContent(effects, ok3) {
    let previous3;
    return chunkStart;
    function chunkStart(code3) {
      effects.enter("content");
      previous3 = effects.enter("chunkContent", {
        contentType: "content"
      });
      return chunkInside(code3);
    }
    function chunkInside(code3) {
      if (code3 === null) {
        return contentEnd(code3);
      }
      if (markdownLineEnding(code3)) {
        return effects.check(continuationConstruct, contentContinue, contentEnd)(code3);
      }
      effects.consume(code3);
      return chunkInside;
    }
    function contentEnd(code3) {
      effects.exit("chunkContent");
      effects.exit("content");
      return ok3(code3);
    }
    function contentContinue(code3) {
      effects.consume(code3);
      effects.exit("chunkContent");
      previous3.next = effects.enter("chunkContent", {
        contentType: "content",
        previous: previous3
      });
      previous3 = previous3.next;
      return chunkInside;
    }
  }
  function tokenizeContinuation(effects, ok3, nok) {
    const self2 = this;
    return startLookahead;
    function startLookahead(code3) {
      effects.exit("chunkContent");
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return factorySpace(effects, prefixed, "linePrefix");
    }
    function prefixed(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        return nok(code3);
      }
      const tail = self2.events[self2.events.length - 1];
      if (!self2.parser.constructs.disable.null.includes("codeIndented") && tail && tail[1].type === "linePrefix" && tail[2].sliceSerialize(tail[1], true).length >= 4) {
        return ok3(code3);
      }
      return effects.interrupt(self2.parser.constructs.flow, nok, ok3)(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-destination@2.0.1/node_modules/micromark-factory-destination/index.js
  function factoryDestination(effects, ok3, nok, type, literalType, literalMarkerType, rawType, stringType, max) {
    const limit = max || Number.POSITIVE_INFINITY;
    let balance = 0;
    return start2;
    function start2(code3) {
      if (code3 === 60) {
        effects.enter(type);
        effects.enter(literalType);
        effects.enter(literalMarkerType);
        effects.consume(code3);
        effects.exit(literalMarkerType);
        return enclosedBefore;
      }
      if (code3 === null || code3 === 32 || code3 === 41 || asciiControl(code3)) {
        return nok(code3);
      }
      effects.enter(type);
      effects.enter(rawType);
      effects.enter(stringType);
      effects.enter("chunkString", {
        contentType: "string"
      });
      return raw(code3);
    }
    function enclosedBefore(code3) {
      if (code3 === 62) {
        effects.enter(literalMarkerType);
        effects.consume(code3);
        effects.exit(literalMarkerType);
        effects.exit(literalType);
        effects.exit(type);
        return ok3;
      }
      effects.enter(stringType);
      effects.enter("chunkString", {
        contentType: "string"
      });
      return enclosed(code3);
    }
    function enclosed(code3) {
      if (code3 === 62) {
        effects.exit("chunkString");
        effects.exit(stringType);
        return enclosedBefore(code3);
      }
      if (code3 === null || code3 === 60 || markdownLineEnding(code3)) {
        return nok(code3);
      }
      effects.consume(code3);
      return code3 === 92 ? enclosedEscape : enclosed;
    }
    function enclosedEscape(code3) {
      if (code3 === 60 || code3 === 62 || code3 === 92) {
        effects.consume(code3);
        return enclosed;
      }
      return enclosed(code3);
    }
    function raw(code3) {
      if (!balance && (code3 === null || code3 === 41 || markdownLineEndingOrSpace(code3))) {
        effects.exit("chunkString");
        effects.exit(stringType);
        effects.exit(rawType);
        effects.exit(type);
        return ok3(code3);
      }
      if (balance < limit && code3 === 40) {
        effects.consume(code3);
        balance++;
        return raw;
      }
      if (code3 === 41) {
        effects.consume(code3);
        balance--;
        return raw;
      }
      if (code3 === null || code3 === 32 || code3 === 40 || asciiControl(code3)) {
        return nok(code3);
      }
      effects.consume(code3);
      return code3 === 92 ? rawEscape : raw;
    }
    function rawEscape(code3) {
      if (code3 === 40 || code3 === 41 || code3 === 92) {
        effects.consume(code3);
        return raw;
      }
      return raw(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-label@2.0.1/node_modules/micromark-factory-label/index.js
  function factoryLabel(effects, ok3, nok, type, markerType, stringType) {
    const self2 = this;
    let size = 0;
    let seen;
    return start2;
    function start2(code3) {
      effects.enter(type);
      effects.enter(markerType);
      effects.consume(code3);
      effects.exit(markerType);
      effects.enter(stringType);
      return atBreak;
    }
    function atBreak(code3) {
      if (size > 999 || code3 === null || code3 === 91 || code3 === 93 && !seen || // To do: remove in the future once we’ve switched from
      // `micromark-extension-footnote` to `micromark-extension-gfm-footnote`,
      // which doesn’t need this.
      // Hidden footnotes hook.
      /* c8 ignore next 3 */
      code3 === 94 && !size && "_hiddenFootnoteSupport" in self2.parser.constructs) {
        return nok(code3);
      }
      if (code3 === 93) {
        effects.exit(stringType);
        effects.enter(markerType);
        effects.consume(code3);
        effects.exit(markerType);
        effects.exit(type);
        return ok3;
      }
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return atBreak;
      }
      effects.enter("chunkString", {
        contentType: "string"
      });
      return labelInside(code3);
    }
    function labelInside(code3) {
      if (code3 === null || code3 === 91 || code3 === 93 || markdownLineEnding(code3) || size++ > 999) {
        effects.exit("chunkString");
        return atBreak(code3);
      }
      effects.consume(code3);
      if (!seen) seen = !markdownSpace(code3);
      return code3 === 92 ? labelEscape : labelInside;
    }
    function labelEscape(code3) {
      if (code3 === 91 || code3 === 92 || code3 === 93) {
        effects.consume(code3);
        size++;
        return labelInside;
      }
      return labelInside(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-title@2.0.1/node_modules/micromark-factory-title/index.js
  function factoryTitle(effects, ok3, nok, type, markerType, stringType) {
    let marker;
    return start2;
    function start2(code3) {
      if (code3 === 34 || code3 === 39 || code3 === 40) {
        effects.enter(type);
        effects.enter(markerType);
        effects.consume(code3);
        effects.exit(markerType);
        marker = code3 === 40 ? 41 : code3;
        return begin;
      }
      return nok(code3);
    }
    function begin(code3) {
      if (code3 === marker) {
        effects.enter(markerType);
        effects.consume(code3);
        effects.exit(markerType);
        effects.exit(type);
        return ok3;
      }
      effects.enter(stringType);
      return atBreak(code3);
    }
    function atBreak(code3) {
      if (code3 === marker) {
        effects.exit(stringType);
        return begin(marker);
      }
      if (code3 === null) {
        return nok(code3);
      }
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return factorySpace(effects, atBreak, "linePrefix");
      }
      effects.enter("chunkString", {
        contentType: "string"
      });
      return inside(code3);
    }
    function inside(code3) {
      if (code3 === marker || code3 === null || markdownLineEnding(code3)) {
        effects.exit("chunkString");
        return atBreak(code3);
      }
      effects.consume(code3);
      return code3 === 92 ? escape : inside;
    }
    function escape(code3) {
      if (code3 === marker || code3 === 92) {
        effects.consume(code3);
        return inside;
      }
      return inside(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-factory-whitespace@2.0.1/node_modules/micromark-factory-whitespace/index.js
  function factoryWhitespace(effects, ok3) {
    let seen;
    return start2;
    function start2(code3) {
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        seen = true;
        return start2;
      }
      if (markdownSpace(code3)) {
        return factorySpace(effects, start2, seen ? "linePrefix" : "lineSuffix")(code3);
      }
      return ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-normalize-identifier@2.0.1/node_modules/micromark-util-normalize-identifier/index.js
  function normalizeIdentifier(value) {
    return value.replace(/[\t\n\r ]+/g, " ").replace(/^ | $/g, "").toLowerCase().toUpperCase();
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/definition.js
  var definition = {
    name: "definition",
    tokenize: tokenizeDefinition
  };
  var titleBefore = {
    partial: true,
    tokenize: tokenizeTitleBefore
  };
  function tokenizeDefinition(effects, ok3, nok) {
    const self2 = this;
    let identifier2;
    return start2;
    function start2(code3) {
      effects.enter("definition");
      return before(code3);
    }
    function before(code3) {
      return factoryLabel.call(
        self2,
        effects,
        labelAfter,
        // Note: we don’t need to reset the way `markdown-rs` does.
        nok,
        "definitionLabel",
        "definitionLabelMarker",
        "definitionLabelString"
      )(code3);
    }
    function labelAfter(code3) {
      identifier2 = normalizeIdentifier(self2.sliceSerialize(self2.events[self2.events.length - 1][1]).slice(1, -1));
      if (code3 === 58) {
        effects.enter("definitionMarker");
        effects.consume(code3);
        effects.exit("definitionMarker");
        return markerAfter;
      }
      return nok(code3);
    }
    function markerAfter(code3) {
      return markdownLineEndingOrSpace(code3) ? factoryWhitespace(effects, destinationBefore)(code3) : destinationBefore(code3);
    }
    function destinationBefore(code3) {
      return factoryDestination(
        effects,
        destinationAfter,
        // Note: we don’t need to reset the way `markdown-rs` does.
        nok,
        "definitionDestination",
        "definitionDestinationLiteral",
        "definitionDestinationLiteralMarker",
        "definitionDestinationRaw",
        "definitionDestinationString"
      )(code3);
    }
    function destinationAfter(code3) {
      return effects.attempt(titleBefore, after, after)(code3);
    }
    function after(code3) {
      return markdownSpace(code3) ? factorySpace(effects, afterWhitespace, "whitespace")(code3) : afterWhitespace(code3);
    }
    function afterWhitespace(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("definition");
        self2.parser.defined.push(identifier2);
        return ok3(code3);
      }
      return nok(code3);
    }
  }
  function tokenizeTitleBefore(effects, ok3, nok) {
    return titleBefore2;
    function titleBefore2(code3) {
      return markdownLineEndingOrSpace(code3) ? factoryWhitespace(effects, beforeMarker)(code3) : nok(code3);
    }
    function beforeMarker(code3) {
      return factoryTitle(effects, titleAfter, nok, "definitionTitle", "definitionTitleMarker", "definitionTitleString")(code3);
    }
    function titleAfter(code3) {
      return markdownSpace(code3) ? factorySpace(effects, titleAfterOptionalWhitespace, "whitespace")(code3) : titleAfterOptionalWhitespace(code3);
    }
    function titleAfterOptionalWhitespace(code3) {
      return code3 === null || markdownLineEnding(code3) ? ok3(code3) : nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/hard-break-escape.js
  var hardBreakEscape = {
    name: "hardBreakEscape",
    tokenize: tokenizeHardBreakEscape
  };
  function tokenizeHardBreakEscape(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      effects.enter("hardBreakEscape");
      effects.consume(code3);
      return after;
    }
    function after(code3) {
      if (markdownLineEnding(code3)) {
        effects.exit("hardBreakEscape");
        return ok3(code3);
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/heading-atx.js
  var headingAtx = {
    name: "headingAtx",
    resolve: resolveHeadingAtx,
    tokenize: tokenizeHeadingAtx
  };
  function resolveHeadingAtx(events, context) {
    let contentEnd = events.length - 2;
    let contentStart = 3;
    let content3;
    let text5;
    if (events[contentStart][1].type === "whitespace") {
      contentStart += 2;
    }
    if (contentEnd - 2 > contentStart && events[contentEnd][1].type === "whitespace") {
      contentEnd -= 2;
    }
    if (events[contentEnd][1].type === "atxHeadingSequence" && (contentStart === contentEnd - 1 || contentEnd - 4 > contentStart && events[contentEnd - 2][1].type === "whitespace")) {
      contentEnd -= contentStart + 1 === contentEnd ? 2 : 4;
    }
    if (contentEnd > contentStart) {
      content3 = {
        type: "atxHeadingText",
        start: events[contentStart][1].start,
        end: events[contentEnd][1].end
      };
      text5 = {
        type: "chunkText",
        start: events[contentStart][1].start,
        end: events[contentEnd][1].end,
        contentType: "text"
      };
      splice(events, contentStart, contentEnd - contentStart + 1, [["enter", content3, context], ["enter", text5, context], ["exit", text5, context], ["exit", content3, context]]);
    }
    return events;
  }
  function tokenizeHeadingAtx(effects, ok3, nok) {
    let size = 0;
    return start2;
    function start2(code3) {
      effects.enter("atxHeading");
      return before(code3);
    }
    function before(code3) {
      effects.enter("atxHeadingSequence");
      return sequenceOpen(code3);
    }
    function sequenceOpen(code3) {
      if (code3 === 35 && size++ < 6) {
        effects.consume(code3);
        return sequenceOpen;
      }
      if (code3 === null || markdownLineEndingOrSpace(code3)) {
        effects.exit("atxHeadingSequence");
        return atBreak(code3);
      }
      return nok(code3);
    }
    function atBreak(code3) {
      if (code3 === 35) {
        effects.enter("atxHeadingSequence");
        return sequenceFurther(code3);
      }
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("atxHeading");
        return ok3(code3);
      }
      if (markdownSpace(code3)) {
        return factorySpace(effects, atBreak, "whitespace")(code3);
      }
      effects.enter("atxHeadingText");
      return data2(code3);
    }
    function sequenceFurther(code3) {
      if (code3 === 35) {
        effects.consume(code3);
        return sequenceFurther;
      }
      effects.exit("atxHeadingSequence");
      return atBreak(code3);
    }
    function data2(code3) {
      if (code3 === null || code3 === 35 || markdownLineEndingOrSpace(code3)) {
        effects.exit("atxHeadingText");
        return atBreak(code3);
      }
      effects.consume(code3);
      return data2;
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-html-tag-name@2.0.1/node_modules/micromark-util-html-tag-name/index.js
  var htmlBlockNames = [
    "address",
    "article",
    "aside",
    "base",
    "basefont",
    "blockquote",
    "body",
    "caption",
    "center",
    "col",
    "colgroup",
    "dd",
    "details",
    "dialog",
    "dir",
    "div",
    "dl",
    "dt",
    "fieldset",
    "figcaption",
    "figure",
    "footer",
    "form",
    "frame",
    "frameset",
    "h1",
    "h2",
    "h3",
    "h4",
    "h5",
    "h6",
    "head",
    "header",
    "hr",
    "html",
    "iframe",
    "legend",
    "li",
    "link",
    "main",
    "menu",
    "menuitem",
    "nav",
    "noframes",
    "ol",
    "optgroup",
    "option",
    "p",
    "param",
    "search",
    "section",
    "summary",
    "table",
    "tbody",
    "td",
    "tfoot",
    "th",
    "thead",
    "title",
    "tr",
    "track",
    "ul"
  ];
  var htmlRawNames = ["pre", "script", "style", "textarea"];

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/html-flow.js
  var htmlFlow = {
    concrete: true,
    name: "htmlFlow",
    resolveTo: resolveToHtmlFlow,
    tokenize: tokenizeHtmlFlow
  };
  var blankLineBefore = {
    partial: true,
    tokenize: tokenizeBlankLineBefore
  };
  var nonLazyContinuationStart = {
    partial: true,
    tokenize: tokenizeNonLazyContinuationStart
  };
  function resolveToHtmlFlow(events) {
    let index2 = events.length;
    while (index2--) {
      if (events[index2][0] === "enter" && events[index2][1].type === "htmlFlow") {
        break;
      }
    }
    if (index2 > 1 && events[index2 - 2][1].type === "linePrefix") {
      events[index2][1].start = events[index2 - 2][1].start;
      events[index2 + 1][1].start = events[index2 - 2][1].start;
      events.splice(index2 - 2, 2);
    }
    return events;
  }
  function tokenizeHtmlFlow(effects, ok3, nok) {
    const self2 = this;
    let marker;
    let closingTag;
    let buffer;
    let index2;
    let markerB;
    return start2;
    function start2(code3) {
      return before(code3);
    }
    function before(code3) {
      effects.enter("htmlFlow");
      effects.enter("htmlFlowData");
      effects.consume(code3);
      return open;
    }
    function open(code3) {
      if (code3 === 33) {
        effects.consume(code3);
        return declarationOpen;
      }
      if (code3 === 47) {
        effects.consume(code3);
        closingTag = true;
        return tagCloseStart;
      }
      if (code3 === 63) {
        effects.consume(code3);
        marker = 3;
        return self2.interrupt ? ok3 : continuationDeclarationInside;
      }
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        buffer = String.fromCharCode(code3);
        return tagName;
      }
      return nok(code3);
    }
    function declarationOpen(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        marker = 2;
        return commentOpenInside;
      }
      if (code3 === 91) {
        effects.consume(code3);
        marker = 5;
        index2 = 0;
        return cdataOpenInside;
      }
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        marker = 4;
        return self2.interrupt ? ok3 : continuationDeclarationInside;
      }
      return nok(code3);
    }
    function commentOpenInside(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return self2.interrupt ? ok3 : continuationDeclarationInside;
      }
      return nok(code3);
    }
    function cdataOpenInside(code3) {
      const value = "CDATA[";
      if (code3 === value.charCodeAt(index2++)) {
        effects.consume(code3);
        if (index2 === value.length) {
          return self2.interrupt ? ok3 : continuation;
        }
        return cdataOpenInside;
      }
      return nok(code3);
    }
    function tagCloseStart(code3) {
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        buffer = String.fromCharCode(code3);
        return tagName;
      }
      return nok(code3);
    }
    function tagName(code3) {
      if (code3 === null || code3 === 47 || code3 === 62 || markdownLineEndingOrSpace(code3)) {
        const slash = code3 === 47;
        const name2 = buffer.toLowerCase();
        if (!slash && !closingTag && htmlRawNames.includes(name2)) {
          marker = 1;
          return self2.interrupt ? ok3(code3) : continuation(code3);
        }
        if (htmlBlockNames.includes(buffer.toLowerCase())) {
          marker = 6;
          if (slash) {
            effects.consume(code3);
            return basicSelfClosing;
          }
          return self2.interrupt ? ok3(code3) : continuation(code3);
        }
        marker = 7;
        return self2.interrupt && !self2.parser.lazy[self2.now().line] ? nok(code3) : closingTag ? completeClosingTagAfter(code3) : completeAttributeNameBefore(code3);
      }
      if (code3 === 45 || asciiAlphanumeric(code3)) {
        effects.consume(code3);
        buffer += String.fromCharCode(code3);
        return tagName;
      }
      return nok(code3);
    }
    function basicSelfClosing(code3) {
      if (code3 === 62) {
        effects.consume(code3);
        return self2.interrupt ? ok3 : continuation;
      }
      return nok(code3);
    }
    function completeClosingTagAfter(code3) {
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return completeClosingTagAfter;
      }
      return completeEnd(code3);
    }
    function completeAttributeNameBefore(code3) {
      if (code3 === 47) {
        effects.consume(code3);
        return completeEnd;
      }
      if (code3 === 58 || code3 === 95 || asciiAlpha(code3)) {
        effects.consume(code3);
        return completeAttributeName;
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return completeAttributeNameBefore;
      }
      return completeEnd(code3);
    }
    function completeAttributeName(code3) {
      if (code3 === 45 || code3 === 46 || code3 === 58 || code3 === 95 || asciiAlphanumeric(code3)) {
        effects.consume(code3);
        return completeAttributeName;
      }
      return completeAttributeNameAfter(code3);
    }
    function completeAttributeNameAfter(code3) {
      if (code3 === 61) {
        effects.consume(code3);
        return completeAttributeValueBefore;
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return completeAttributeNameAfter;
      }
      return completeAttributeNameBefore(code3);
    }
    function completeAttributeValueBefore(code3) {
      if (code3 === null || code3 === 60 || code3 === 61 || code3 === 62 || code3 === 96) {
        return nok(code3);
      }
      if (code3 === 34 || code3 === 39) {
        effects.consume(code3);
        markerB = code3;
        return completeAttributeValueQuoted;
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return completeAttributeValueBefore;
      }
      return completeAttributeValueUnquoted(code3);
    }
    function completeAttributeValueQuoted(code3) {
      if (code3 === markerB) {
        effects.consume(code3);
        markerB = null;
        return completeAttributeValueQuotedAfter;
      }
      if (code3 === null || markdownLineEnding(code3)) {
        return nok(code3);
      }
      effects.consume(code3);
      return completeAttributeValueQuoted;
    }
    function completeAttributeValueUnquoted(code3) {
      if (code3 === null || code3 === 34 || code3 === 39 || code3 === 47 || code3 === 60 || code3 === 61 || code3 === 62 || code3 === 96 || markdownLineEndingOrSpace(code3)) {
        return completeAttributeNameAfter(code3);
      }
      effects.consume(code3);
      return completeAttributeValueUnquoted;
    }
    function completeAttributeValueQuotedAfter(code3) {
      if (code3 === 47 || code3 === 62 || markdownSpace(code3)) {
        return completeAttributeNameBefore(code3);
      }
      return nok(code3);
    }
    function completeEnd(code3) {
      if (code3 === 62) {
        effects.consume(code3);
        return completeAfter;
      }
      return nok(code3);
    }
    function completeAfter(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        return continuation(code3);
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return completeAfter;
      }
      return nok(code3);
    }
    function continuation(code3) {
      if (code3 === 45 && marker === 2) {
        effects.consume(code3);
        return continuationCommentInside;
      }
      if (code3 === 60 && marker === 1) {
        effects.consume(code3);
        return continuationRawTagOpen;
      }
      if (code3 === 62 && marker === 4) {
        effects.consume(code3);
        return continuationClose;
      }
      if (code3 === 63 && marker === 3) {
        effects.consume(code3);
        return continuationDeclarationInside;
      }
      if (code3 === 93 && marker === 5) {
        effects.consume(code3);
        return continuationCdataInside;
      }
      if (markdownLineEnding(code3) && (marker === 6 || marker === 7)) {
        effects.exit("htmlFlowData");
        return effects.check(blankLineBefore, continuationAfter, continuationStart)(code3);
      }
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("htmlFlowData");
        return continuationStart(code3);
      }
      effects.consume(code3);
      return continuation;
    }
    function continuationStart(code3) {
      return effects.check(nonLazyContinuationStart, continuationStartNonLazy, continuationAfter)(code3);
    }
    function continuationStartNonLazy(code3) {
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return continuationBefore;
    }
    function continuationBefore(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        return continuationStart(code3);
      }
      effects.enter("htmlFlowData");
      return continuation(code3);
    }
    function continuationCommentInside(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return continuationDeclarationInside;
      }
      return continuation(code3);
    }
    function continuationRawTagOpen(code3) {
      if (code3 === 47) {
        effects.consume(code3);
        buffer = "";
        return continuationRawEndTag;
      }
      return continuation(code3);
    }
    function continuationRawEndTag(code3) {
      if (code3 === 62) {
        const name2 = buffer.toLowerCase();
        if (htmlRawNames.includes(name2)) {
          effects.consume(code3);
          return continuationClose;
        }
        return continuation(code3);
      }
      if (asciiAlpha(code3) && buffer.length < 8) {
        effects.consume(code3);
        buffer += String.fromCharCode(code3);
        return continuationRawEndTag;
      }
      return continuation(code3);
    }
    function continuationCdataInside(code3) {
      if (code3 === 93) {
        effects.consume(code3);
        return continuationDeclarationInside;
      }
      return continuation(code3);
    }
    function continuationDeclarationInside(code3) {
      if (code3 === 62) {
        effects.consume(code3);
        return continuationClose;
      }
      if (code3 === 45 && marker === 2) {
        effects.consume(code3);
        return continuationDeclarationInside;
      }
      return continuation(code3);
    }
    function continuationClose(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("htmlFlowData");
        return continuationAfter(code3);
      }
      effects.consume(code3);
      return continuationClose;
    }
    function continuationAfter(code3) {
      effects.exit("htmlFlow");
      return ok3(code3);
    }
  }
  function tokenizeNonLazyContinuationStart(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      if (markdownLineEnding(code3)) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return after;
      }
      return nok(code3);
    }
    function after(code3) {
      return self2.parser.lazy[self2.now().line] ? nok(code3) : ok3(code3);
    }
  }
  function tokenizeBlankLineBefore(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return effects.attempt(blankLine, ok3, nok);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/html-text.js
  var htmlText = {
    name: "htmlText",
    tokenize: tokenizeHtmlText
  };
  function tokenizeHtmlText(effects, ok3, nok) {
    const self2 = this;
    let marker;
    let index2;
    let returnState;
    return start2;
    function start2(code3) {
      effects.enter("htmlText");
      effects.enter("htmlTextData");
      effects.consume(code3);
      return open;
    }
    function open(code3) {
      if (code3 === 33) {
        effects.consume(code3);
        return declarationOpen;
      }
      if (code3 === 47) {
        effects.consume(code3);
        return tagCloseStart;
      }
      if (code3 === 63) {
        effects.consume(code3);
        return instruction;
      }
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        return tagOpen;
      }
      return nok(code3);
    }
    function declarationOpen(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return commentOpenInside;
      }
      if (code3 === 91) {
        effects.consume(code3);
        index2 = 0;
        return cdataOpenInside;
      }
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        return declaration;
      }
      return nok(code3);
    }
    function commentOpenInside(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return commentEnd;
      }
      return nok(code3);
    }
    function comment(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      if (code3 === 45) {
        effects.consume(code3);
        return commentClose;
      }
      if (markdownLineEnding(code3)) {
        returnState = comment;
        return lineEndingBefore(code3);
      }
      effects.consume(code3);
      return comment;
    }
    function commentClose(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return commentEnd;
      }
      return comment(code3);
    }
    function commentEnd(code3) {
      return code3 === 62 ? end(code3) : code3 === 45 ? commentClose(code3) : comment(code3);
    }
    function cdataOpenInside(code3) {
      const value = "CDATA[";
      if (code3 === value.charCodeAt(index2++)) {
        effects.consume(code3);
        return index2 === value.length ? cdata : cdataOpenInside;
      }
      return nok(code3);
    }
    function cdata(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      if (code3 === 93) {
        effects.consume(code3);
        return cdataClose;
      }
      if (markdownLineEnding(code3)) {
        returnState = cdata;
        return lineEndingBefore(code3);
      }
      effects.consume(code3);
      return cdata;
    }
    function cdataClose(code3) {
      if (code3 === 93) {
        effects.consume(code3);
        return cdataEnd;
      }
      return cdata(code3);
    }
    function cdataEnd(code3) {
      if (code3 === 62) {
        return end(code3);
      }
      if (code3 === 93) {
        effects.consume(code3);
        return cdataEnd;
      }
      return cdata(code3);
    }
    function declaration(code3) {
      if (code3 === null || code3 === 62) {
        return end(code3);
      }
      if (markdownLineEnding(code3)) {
        returnState = declaration;
        return lineEndingBefore(code3);
      }
      effects.consume(code3);
      return declaration;
    }
    function instruction(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      if (code3 === 63) {
        effects.consume(code3);
        return instructionClose;
      }
      if (markdownLineEnding(code3)) {
        returnState = instruction;
        return lineEndingBefore(code3);
      }
      effects.consume(code3);
      return instruction;
    }
    function instructionClose(code3) {
      return code3 === 62 ? end(code3) : instruction(code3);
    }
    function tagCloseStart(code3) {
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        return tagClose;
      }
      return nok(code3);
    }
    function tagClose(code3) {
      if (code3 === 45 || asciiAlphanumeric(code3)) {
        effects.consume(code3);
        return tagClose;
      }
      return tagCloseBetween(code3);
    }
    function tagCloseBetween(code3) {
      if (markdownLineEnding(code3)) {
        returnState = tagCloseBetween;
        return lineEndingBefore(code3);
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return tagCloseBetween;
      }
      return end(code3);
    }
    function tagOpen(code3) {
      if (code3 === 45 || asciiAlphanumeric(code3)) {
        effects.consume(code3);
        return tagOpen;
      }
      if (code3 === 47 || code3 === 62 || markdownLineEndingOrSpace(code3)) {
        return tagOpenBetween(code3);
      }
      return nok(code3);
    }
    function tagOpenBetween(code3) {
      if (code3 === 47) {
        effects.consume(code3);
        return end;
      }
      if (code3 === 58 || code3 === 95 || asciiAlpha(code3)) {
        effects.consume(code3);
        return tagOpenAttributeName;
      }
      if (markdownLineEnding(code3)) {
        returnState = tagOpenBetween;
        return lineEndingBefore(code3);
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return tagOpenBetween;
      }
      return end(code3);
    }
    function tagOpenAttributeName(code3) {
      if (code3 === 45 || code3 === 46 || code3 === 58 || code3 === 95 || asciiAlphanumeric(code3)) {
        effects.consume(code3);
        return tagOpenAttributeName;
      }
      return tagOpenAttributeNameAfter(code3);
    }
    function tagOpenAttributeNameAfter(code3) {
      if (code3 === 61) {
        effects.consume(code3);
        return tagOpenAttributeValueBefore;
      }
      if (markdownLineEnding(code3)) {
        returnState = tagOpenAttributeNameAfter;
        return lineEndingBefore(code3);
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return tagOpenAttributeNameAfter;
      }
      return tagOpenBetween(code3);
    }
    function tagOpenAttributeValueBefore(code3) {
      if (code3 === null || code3 === 60 || code3 === 61 || code3 === 62 || code3 === 96) {
        return nok(code3);
      }
      if (code3 === 34 || code3 === 39) {
        effects.consume(code3);
        marker = code3;
        return tagOpenAttributeValueQuoted;
      }
      if (markdownLineEnding(code3)) {
        returnState = tagOpenAttributeValueBefore;
        return lineEndingBefore(code3);
      }
      if (markdownSpace(code3)) {
        effects.consume(code3);
        return tagOpenAttributeValueBefore;
      }
      effects.consume(code3);
      return tagOpenAttributeValueUnquoted;
    }
    function tagOpenAttributeValueQuoted(code3) {
      if (code3 === marker) {
        effects.consume(code3);
        marker = void 0;
        return tagOpenAttributeValueQuotedAfter;
      }
      if (code3 === null) {
        return nok(code3);
      }
      if (markdownLineEnding(code3)) {
        returnState = tagOpenAttributeValueQuoted;
        return lineEndingBefore(code3);
      }
      effects.consume(code3);
      return tagOpenAttributeValueQuoted;
    }
    function tagOpenAttributeValueUnquoted(code3) {
      if (code3 === null || code3 === 34 || code3 === 39 || code3 === 60 || code3 === 61 || code3 === 96) {
        return nok(code3);
      }
      if (code3 === 47 || code3 === 62 || markdownLineEndingOrSpace(code3)) {
        return tagOpenBetween(code3);
      }
      effects.consume(code3);
      return tagOpenAttributeValueUnquoted;
    }
    function tagOpenAttributeValueQuotedAfter(code3) {
      if (code3 === 47 || code3 === 62 || markdownLineEndingOrSpace(code3)) {
        return tagOpenBetween(code3);
      }
      return nok(code3);
    }
    function end(code3) {
      if (code3 === 62) {
        effects.consume(code3);
        effects.exit("htmlTextData");
        effects.exit("htmlText");
        return ok3;
      }
      return nok(code3);
    }
    function lineEndingBefore(code3) {
      effects.exit("htmlTextData");
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return lineEndingAfter;
    }
    function lineEndingAfter(code3) {
      return markdownSpace(code3) ? factorySpace(effects, lineEndingAfterPrefix, "linePrefix", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4)(code3) : lineEndingAfterPrefix(code3);
    }
    function lineEndingAfterPrefix(code3) {
      effects.enter("htmlTextData");
      return returnState(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/label-end.js
  var labelEnd = {
    name: "labelEnd",
    resolveAll: resolveAllLabelEnd,
    resolveTo: resolveToLabelEnd,
    tokenize: tokenizeLabelEnd
  };
  var resourceConstruct = {
    tokenize: tokenizeResource
  };
  var referenceFullConstruct = {
    tokenize: tokenizeReferenceFull
  };
  var referenceCollapsedConstruct = {
    tokenize: tokenizeReferenceCollapsed
  };
  function resolveAllLabelEnd(events) {
    let index2 = -1;
    const newEvents = [];
    while (++index2 < events.length) {
      const token = events[index2][1];
      newEvents.push(events[index2]);
      if (token.type === "labelImage" || token.type === "labelLink" || token.type === "labelEnd") {
        const offset2 = token.type === "labelImage" ? 4 : 2;
        token.type = "data";
        index2 += offset2;
      }
    }
    if (events.length !== newEvents.length) {
      splice(events, 0, events.length, newEvents);
    }
    return events;
  }
  function resolveToLabelEnd(events, context) {
    let index2 = events.length;
    let offset2 = 0;
    let token;
    let open;
    let close;
    let media;
    while (index2--) {
      token = events[index2][1];
      if (open) {
        if (token.type === "link" || token.type === "labelLink" && token._inactive) {
          break;
        }
        if (events[index2][0] === "enter" && token.type === "labelLink") {
          token._inactive = true;
        }
      } else if (close) {
        if (events[index2][0] === "enter" && (token.type === "labelImage" || token.type === "labelLink") && !token._balanced) {
          open = index2;
          if (token.type !== "labelLink") {
            offset2 = 2;
            break;
          }
        }
      } else if (token.type === "labelEnd") {
        close = index2;
      }
    }
    const group = {
      type: events[open][1].type === "labelLink" ? "link" : "image",
      start: {
        ...events[open][1].start
      },
      end: {
        ...events[events.length - 1][1].end
      }
    };
    const label = {
      type: "label",
      start: {
        ...events[open][1].start
      },
      end: {
        ...events[close][1].end
      }
    };
    const text5 = {
      type: "labelText",
      start: {
        ...events[open + offset2 + 2][1].end
      },
      end: {
        ...events[close - 2][1].start
      }
    };
    media = [["enter", group, context], ["enter", label, context]];
    media = push(media, events.slice(open + 1, open + offset2 + 3));
    media = push(media, [["enter", text5, context]]);
    media = push(media, resolveAll(context.parser.constructs.insideSpan.null, events.slice(open + offset2 + 4, close - 3), context));
    media = push(media, [["exit", text5, context], events[close - 2], events[close - 1], ["exit", label, context]]);
    media = push(media, events.slice(close + 1));
    media = push(media, [["exit", group, context]]);
    splice(events, open, events.length, media);
    return events;
  }
  function tokenizeLabelEnd(effects, ok3, nok) {
    const self2 = this;
    let index2 = self2.events.length;
    let labelStart;
    let defined;
    while (index2--) {
      if ((self2.events[index2][1].type === "labelImage" || self2.events[index2][1].type === "labelLink") && !self2.events[index2][1]._balanced) {
        labelStart = self2.events[index2][1];
        break;
      }
    }
    return start2;
    function start2(code3) {
      if (!labelStart) {
        return nok(code3);
      }
      if (labelStart._inactive) {
        return labelEndNok(code3);
      }
      defined = self2.parser.defined.includes(normalizeIdentifier(self2.sliceSerialize({
        start: labelStart.end,
        end: self2.now()
      })));
      effects.enter("labelEnd");
      effects.enter("labelMarker");
      effects.consume(code3);
      effects.exit("labelMarker");
      effects.exit("labelEnd");
      return after;
    }
    function after(code3) {
      if (code3 === 40) {
        return effects.attempt(resourceConstruct, labelEndOk, defined ? labelEndOk : labelEndNok)(code3);
      }
      if (code3 === 91) {
        return effects.attempt(referenceFullConstruct, labelEndOk, defined ? referenceNotFull : labelEndNok)(code3);
      }
      return defined ? labelEndOk(code3) : labelEndNok(code3);
    }
    function referenceNotFull(code3) {
      return effects.attempt(referenceCollapsedConstruct, labelEndOk, labelEndNok)(code3);
    }
    function labelEndOk(code3) {
      return ok3(code3);
    }
    function labelEndNok(code3) {
      labelStart._balanced = true;
      return nok(code3);
    }
  }
  function tokenizeResource(effects, ok3, nok) {
    return resourceStart;
    function resourceStart(code3) {
      effects.enter("resource");
      effects.enter("resourceMarker");
      effects.consume(code3);
      effects.exit("resourceMarker");
      return resourceBefore;
    }
    function resourceBefore(code3) {
      return markdownLineEndingOrSpace(code3) ? factoryWhitespace(effects, resourceOpen)(code3) : resourceOpen(code3);
    }
    function resourceOpen(code3) {
      if (code3 === 41) {
        return resourceEnd(code3);
      }
      return factoryDestination(effects, resourceDestinationAfter, resourceDestinationMissing, "resourceDestination", "resourceDestinationLiteral", "resourceDestinationLiteralMarker", "resourceDestinationRaw", "resourceDestinationString", 32)(code3);
    }
    function resourceDestinationAfter(code3) {
      return markdownLineEndingOrSpace(code3) ? factoryWhitespace(effects, resourceBetween)(code3) : resourceEnd(code3);
    }
    function resourceDestinationMissing(code3) {
      return nok(code3);
    }
    function resourceBetween(code3) {
      if (code3 === 34 || code3 === 39 || code3 === 40) {
        return factoryTitle(effects, resourceTitleAfter, nok, "resourceTitle", "resourceTitleMarker", "resourceTitleString")(code3);
      }
      return resourceEnd(code3);
    }
    function resourceTitleAfter(code3) {
      return markdownLineEndingOrSpace(code3) ? factoryWhitespace(effects, resourceEnd)(code3) : resourceEnd(code3);
    }
    function resourceEnd(code3) {
      if (code3 === 41) {
        effects.enter("resourceMarker");
        effects.consume(code3);
        effects.exit("resourceMarker");
        effects.exit("resource");
        return ok3;
      }
      return nok(code3);
    }
  }
  function tokenizeReferenceFull(effects, ok3, nok) {
    const self2 = this;
    return referenceFull;
    function referenceFull(code3) {
      return factoryLabel.call(self2, effects, referenceFullAfter, referenceFullMissing, "reference", "referenceMarker", "referenceString")(code3);
    }
    function referenceFullAfter(code3) {
      return self2.parser.defined.includes(normalizeIdentifier(self2.sliceSerialize(self2.events[self2.events.length - 1][1]).slice(1, -1))) ? ok3(code3) : nok(code3);
    }
    function referenceFullMissing(code3) {
      return nok(code3);
    }
  }
  function tokenizeReferenceCollapsed(effects, ok3, nok) {
    return referenceCollapsedStart;
    function referenceCollapsedStart(code3) {
      effects.enter("reference");
      effects.enter("referenceMarker");
      effects.consume(code3);
      effects.exit("referenceMarker");
      return referenceCollapsedOpen;
    }
    function referenceCollapsedOpen(code3) {
      if (code3 === 93) {
        effects.enter("referenceMarker");
        effects.consume(code3);
        effects.exit("referenceMarker");
        effects.exit("reference");
        return ok3;
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/label-start-image.js
  var labelStartImage = {
    name: "labelStartImage",
    resolveAll: labelEnd.resolveAll,
    tokenize: tokenizeLabelStartImage
  };
  function tokenizeLabelStartImage(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      effects.enter("labelImage");
      effects.enter("labelImageMarker");
      effects.consume(code3);
      effects.exit("labelImageMarker");
      return open;
    }
    function open(code3) {
      if (code3 === 91) {
        effects.enter("labelMarker");
        effects.consume(code3);
        effects.exit("labelMarker");
        effects.exit("labelImage");
        return after;
      }
      return nok(code3);
    }
    function after(code3) {
      return code3 === 94 && "_hiddenFootnoteSupport" in self2.parser.constructs ? nok(code3) : ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/label-start-link.js
  var labelStartLink = {
    name: "labelStartLink",
    resolveAll: labelEnd.resolveAll,
    tokenize: tokenizeLabelStartLink
  };
  function tokenizeLabelStartLink(effects, ok3, nok) {
    const self2 = this;
    return start2;
    function start2(code3) {
      effects.enter("labelLink");
      effects.enter("labelMarker");
      effects.consume(code3);
      effects.exit("labelMarker");
      effects.exit("labelLink");
      return after;
    }
    function after(code3) {
      return code3 === 94 && "_hiddenFootnoteSupport" in self2.parser.constructs ? nok(code3) : ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/line-ending.js
  var lineEnding = {
    name: "lineEnding",
    tokenize: tokenizeLineEnding
  };
  function tokenizeLineEnding(effects, ok3) {
    return start2;
    function start2(code3) {
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return factorySpace(effects, ok3, "linePrefix");
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/thematic-break.js
  var thematicBreak = {
    name: "thematicBreak",
    tokenize: tokenizeThematicBreak
  };
  function tokenizeThematicBreak(effects, ok3, nok) {
    let size = 0;
    let marker;
    return start2;
    function start2(code3) {
      effects.enter("thematicBreak");
      return before(code3);
    }
    function before(code3) {
      marker = code3;
      return atBreak(code3);
    }
    function atBreak(code3) {
      if (code3 === marker) {
        effects.enter("thematicBreakSequence");
        return sequence(code3);
      }
      if (size >= 3 && (code3 === null || markdownLineEnding(code3))) {
        effects.exit("thematicBreak");
        return ok3(code3);
      }
      return nok(code3);
    }
    function sequence(code3) {
      if (code3 === marker) {
        effects.consume(code3);
        size++;
        return sequence;
      }
      effects.exit("thematicBreakSequence");
      return markdownSpace(code3) ? factorySpace(effects, atBreak, "whitespace")(code3) : atBreak(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/list.js
  var list = {
    continuation: {
      tokenize: tokenizeListContinuation
    },
    exit: tokenizeListEnd,
    name: "list",
    tokenize: tokenizeListStart
  };
  var listItemPrefixWhitespaceConstruct = {
    partial: true,
    tokenize: tokenizeListItemPrefixWhitespace
  };
  var indentConstruct = {
    partial: true,
    tokenize: tokenizeIndent
  };
  function tokenizeListStart(effects, ok3, nok) {
    const self2 = this;
    const tail = self2.events[self2.events.length - 1];
    let initialSize = tail && tail[1].type === "linePrefix" ? tail[2].sliceSerialize(tail[1], true).length : 0;
    let size = 0;
    return start2;
    function start2(code3) {
      const kind = self2.containerState.type || (code3 === 42 || code3 === 43 || code3 === 45 ? "listUnordered" : "listOrdered");
      if (kind === "listUnordered" ? !self2.containerState.marker || code3 === self2.containerState.marker : asciiDigit(code3)) {
        if (!self2.containerState.type) {
          self2.containerState.type = kind;
          effects.enter(kind, {
            _container: true
          });
        }
        if (kind === "listUnordered") {
          effects.enter("listItemPrefix");
          return code3 === 42 || code3 === 45 ? effects.check(thematicBreak, nok, atMarker)(code3) : atMarker(code3);
        }
        if (!self2.interrupt || code3 === 49) {
          effects.enter("listItemPrefix");
          effects.enter("listItemValue");
          return inside(code3);
        }
      }
      return nok(code3);
    }
    function inside(code3) {
      if (asciiDigit(code3) && ++size < 10) {
        effects.consume(code3);
        return inside;
      }
      if ((!self2.interrupt || size < 2) && (self2.containerState.marker ? code3 === self2.containerState.marker : code3 === 41 || code3 === 46)) {
        effects.exit("listItemValue");
        return atMarker(code3);
      }
      return nok(code3);
    }
    function atMarker(code3) {
      effects.enter("listItemMarker");
      effects.consume(code3);
      effects.exit("listItemMarker");
      self2.containerState.marker = self2.containerState.marker || code3;
      return effects.check(
        blankLine,
        // Can’t be empty when interrupting.
        self2.interrupt ? nok : onBlank,
        effects.attempt(listItemPrefixWhitespaceConstruct, endOfPrefix, otherPrefix)
      );
    }
    function onBlank(code3) {
      self2.containerState.initialBlankLine = true;
      initialSize++;
      return endOfPrefix(code3);
    }
    function otherPrefix(code3) {
      if (markdownSpace(code3)) {
        effects.enter("listItemPrefixWhitespace");
        effects.consume(code3);
        effects.exit("listItemPrefixWhitespace");
        return endOfPrefix;
      }
      return nok(code3);
    }
    function endOfPrefix(code3) {
      self2.containerState.size = initialSize + self2.sliceSerialize(effects.exit("listItemPrefix"), true).length;
      return ok3(code3);
    }
  }
  function tokenizeListContinuation(effects, ok3, nok) {
    const self2 = this;
    self2.containerState._closeFlow = void 0;
    return effects.check(blankLine, onBlank, notBlank);
    function onBlank(code3) {
      self2.containerState.furtherBlankLines = self2.containerState.furtherBlankLines || self2.containerState.initialBlankLine;
      return factorySpace(effects, ok3, "listItemIndent", self2.containerState.size + 1)(code3);
    }
    function notBlank(code3) {
      if (self2.containerState.furtherBlankLines || !markdownSpace(code3)) {
        self2.containerState.furtherBlankLines = void 0;
        self2.containerState.initialBlankLine = void 0;
        return notInCurrentItem(code3);
      }
      self2.containerState.furtherBlankLines = void 0;
      self2.containerState.initialBlankLine = void 0;
      return effects.attempt(indentConstruct, ok3, notInCurrentItem)(code3);
    }
    function notInCurrentItem(code3) {
      self2.containerState._closeFlow = true;
      self2.interrupt = void 0;
      return factorySpace(effects, effects.attempt(list, ok3, nok), "linePrefix", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4)(code3);
    }
  }
  function tokenizeIndent(effects, ok3, nok) {
    const self2 = this;
    return factorySpace(effects, afterPrefix, "listItemIndent", self2.containerState.size + 1);
    function afterPrefix(code3) {
      const tail = self2.events[self2.events.length - 1];
      return tail && tail[1].type === "listItemIndent" && tail[2].sliceSerialize(tail[1], true).length === self2.containerState.size ? ok3(code3) : nok(code3);
    }
  }
  function tokenizeListEnd(effects) {
    effects.exit(this.containerState.type);
  }
  function tokenizeListItemPrefixWhitespace(effects, ok3, nok) {
    const self2 = this;
    return factorySpace(effects, afterPrefix, "listItemPrefixWhitespace", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4 + 1);
    function afterPrefix(code3) {
      const tail = self2.events[self2.events.length - 1];
      return !markdownSpace(code3) && tail && tail[1].type === "listItemPrefixWhitespace" ? ok3(code3) : nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-core-commonmark@2.0.3/node_modules/micromark-core-commonmark/lib/setext-underline.js
  var setextUnderline = {
    name: "setextUnderline",
    resolveTo: resolveToSetextUnderline,
    tokenize: tokenizeSetextUnderline
  };
  function resolveToSetextUnderline(events, context) {
    let index2 = events.length;
    let content3;
    let text5;
    let definition3;
    while (index2--) {
      if (events[index2][0] === "enter") {
        if (events[index2][1].type === "content") {
          content3 = index2;
          break;
        }
        if (events[index2][1].type === "paragraph") {
          text5 = index2;
        }
      } else {
        if (events[index2][1].type === "content") {
          events.splice(index2, 1);
        }
        if (!definition3 && events[index2][1].type === "definition") {
          definition3 = index2;
        }
      }
    }
    const heading2 = {
      type: "setextHeading",
      start: {
        ...events[content3][1].start
      },
      end: {
        ...events[events.length - 1][1].end
      }
    };
    events[text5][1].type = "setextHeadingText";
    if (definition3) {
      events.splice(text5, 0, ["enter", heading2, context]);
      events.splice(definition3 + 1, 0, ["exit", events[content3][1], context]);
      events[content3][1].end = {
        ...events[definition3][1].end
      };
    } else {
      events[content3][1] = heading2;
    }
    events.push(["exit", heading2, context]);
    return events;
  }
  function tokenizeSetextUnderline(effects, ok3, nok) {
    const self2 = this;
    let marker;
    return start2;
    function start2(code3) {
      let index2 = self2.events.length;
      let paragraph2;
      while (index2--) {
        if (self2.events[index2][1].type !== "lineEnding" && self2.events[index2][1].type !== "linePrefix" && self2.events[index2][1].type !== "content") {
          paragraph2 = self2.events[index2][1].type === "paragraph";
          break;
        }
      }
      if (!self2.parser.lazy[self2.now().line] && (self2.interrupt || paragraph2)) {
        effects.enter("setextHeadingLine");
        marker = code3;
        return before(code3);
      }
      return nok(code3);
    }
    function before(code3) {
      effects.enter("setextHeadingLineSequence");
      return inside(code3);
    }
    function inside(code3) {
      if (code3 === marker) {
        effects.consume(code3);
        return inside;
      }
      effects.exit("setextHeadingLineSequence");
      return markdownSpace(code3) ? factorySpace(effects, after, "lineSuffix")(code3) : after(code3);
    }
    function after(code3) {
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("setextHeadingLine");
        return ok3(code3);
      }
      return nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdxjs-esm@3.0.0/node_modules/micromark-extension-mdxjs-esm/lib/syntax.js
  var blankLineBefore2 = {
    tokenize: tokenizeNextBlank,
    partial: true
  };
  var trouble3 = "https://github.com/micromark/micromark-extension-mdxjs-esm";
  var allowedAcornTypes = /* @__PURE__ */ new Set(["ExportAllDeclaration", "ExportDefaultDeclaration", "ExportNamedDeclaration", "ImportDeclaration"]);
  function mdxjsEsm(options) {
    const exportImportConstruct = {
      tokenize: tokenizeExportImport,
      concrete: true
    };
    if (!options || !options.acorn || !options.acorn.parse) {
      throw new Error("Expected an `acorn` instance passed in as `options.acorn`");
    }
    const acorn = options.acorn;
    const acornOptions = Object.assign({
      ecmaVersion: 2024,
      sourceType: "module"
    }, options.acornOptions, {
      locations: true
    });
    return {
      flow: {
        [101]: exportImportConstruct,
        [105]: exportImportConstruct
      }
    };
    function tokenizeExportImport(effects, ok3, nok) {
      const self2 = this;
      const definedModuleSpecifiers = self2.parser.definedModuleSpecifiers || (self2.parser.definedModuleSpecifiers = []);
      const eventStart = this.events.length + 1;
      let buffer = "";
      return self2.interrupt ? nok : start2;
      function start2(code3) {
        if (self2.now().column > 1) return nok(code3);
        effects.enter("mdxjsEsm");
        effects.enter("mdxjsEsmData");
        effects.consume(code3);
        buffer += String.fromCharCode(code3);
        return word;
      }
      function word(code3) {
        if (asciiAlpha(code3)) {
          effects.consume(code3);
          buffer += String.fromCharCode(code3);
          return word;
        }
        if ((buffer === "import" || buffer === "export") && code3 === 32) {
          effects.consume(code3);
          return inside;
        }
        return nok(code3);
      }
      function inside(code3) {
        if (code3 === null || markdownLineEnding(code3)) {
          effects.exit("mdxjsEsmData");
          return lineStart(code3);
        }
        effects.consume(code3);
        return inside;
      }
      function lineStart(code3) {
        if (code3 === null) {
          return atEnd(code3);
        }
        if (markdownLineEnding(code3)) {
          return effects.check(blankLineBefore2, atEnd, continuationStart)(code3);
        }
        effects.enter("mdxjsEsmData");
        return inside(code3);
      }
      function continuationStart(code3) {
        effects.enter("lineEnding");
        effects.consume(code3);
        effects.exit("lineEnding");
        return lineStart;
      }
      function atEnd(code3) {
        const result = eventsToAcorn(self2.events.slice(eventStart), {
          acorn,
          acornOptions,
          tokenTypes: ["mdxjsEsmData"],
          prefix: definedModuleSpecifiers.length > 0 ? "var " + definedModuleSpecifiers.join(",") + "\n" : ""
        });
        if (result.error) {
          if (code3 !== null && result.swallow) {
            return continuationStart(code3);
          }
          const error = new VFileMessage("Could not parse import/exports with acorn", {
            cause: result.error,
            place: {
              line: result.error.loc.line,
              column: result.error.loc.column + 1,
              offset: result.error.pos
            },
            ruleId: "acorn",
            source: "micromark-extension-mdxjs-esm"
          });
          error.url = trouble3 + "#could-not-parse-importexports-with-acorn";
          throw error;
        }
        if (definedModuleSpecifiers.length > 0) {
          const declaration = result.estree.body.shift();
        }
        let index2 = -1;
        while (++index2 < result.estree.body.length) {
          const node2 = result.estree.body[index2];
          if (!allowedAcornTypes.has(node2.type)) {
            const error = new VFileMessage("Unexpected `" + node2.type + "` in code: only import/exports are supported", {
              place: positionFromEstree(node2),
              ruleId: "non-esm",
              source: "micromark-extension-mdxjs-esm"
            });
            error.url = trouble3 + "#unexpected-type-in-code-only-importexports-are-supported";
            throw error;
          }
          if (node2.type === "ImportDeclaration" && !self2.interrupt) {
            let index3 = -1;
            while (++index3 < node2.specifiers.length) {
              const specifier = node2.specifiers[index3];
              definedModuleSpecifiers.push(specifier.local.name);
            }
          }
        }
        Object.assign(effects.exit("mdxjsEsm"), options.addResult ? {
          estree: result.estree
        } : void 0);
        return ok3(code3);
      }
    }
  }
  function tokenizeNextBlank(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      effects.enter("lineEndingBlank");
      effects.consume(code3);
      effects.exit("lineEndingBlank");
      return effects.attempt(blankLine, ok3, nok);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-combine-extensions@2.0.1/node_modules/micromark-util-combine-extensions/index.js
  var hasOwnProperty2 = {}.hasOwnProperty;
  function combineExtensions(extensions) {
    const all2 = {};
    let index2 = -1;
    while (++index2 < extensions.length) {
      syntaxExtension(all2, extensions[index2]);
    }
    return all2;
  }
  function syntaxExtension(all2, extension2) {
    let hook;
    for (hook in extension2) {
      const maybe = hasOwnProperty2.call(all2, hook) ? all2[hook] : void 0;
      const left = maybe || (all2[hook] = {});
      const right = extension2[hook];
      let code3;
      if (right) {
        for (code3 in right) {
          if (!hasOwnProperty2.call(left, code3)) left[code3] = [];
          const value = right[code3];
          constructs(
            // @ts-expect-error Looks like a list.
            left[code3],
            Array.isArray(value) ? value : value ? [value] : []
          );
        }
      }
    }
  }
  function constructs(existing, list4) {
    let index2 = -1;
    const before = [];
    while (++index2 < list4.length) {
      ;
      (list4[index2].add === "after" ? existing : before).push(list4[index2]);
    }
    splice(existing, 0, 0, before);
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-mdxjs@3.0.0/node_modules/micromark-extension-mdxjs/index.js
  function mdxjs(options) {
    const settings = Object.assign(
      {
        acorn: Parser.extend((0, import_acorn_jsx.default)()),
        acornOptions: { ecmaVersion: 2024, sourceType: "module" },
        addResult: true
      },
      options
    );
    return combineExtensions([
      mdxjsEsm(settings),
      mdxExpression(settings),
      mdxJsx(settings),
      mdxMd()
    ]);
  }

  // ../../../ui/node_modules/.pnpm/remark-mdx@3.1.1/node_modules/remark-mdx/lib/index.js
  var emptyOptions2 = {};
  function remarkMdx(options) {
    const self2 = this;
    const settings = options || emptyOptions2;
    const data2 = self2.data();
    const micromarkExtensions = data2.micromarkExtensions || (data2.micromarkExtensions = []);
    const fromMarkdownExtensions = data2.fromMarkdownExtensions || (data2.fromMarkdownExtensions = []);
    const toMarkdownExtensions = data2.toMarkdownExtensions || (data2.toMarkdownExtensions = []);
    micromarkExtensions.push(mdxjs(settings));
    fromMarkdownExtensions.push(mdxFromMarkdown());
    toMarkdownExtensions.push(mdxToMarkdown(settings));
  }

  // ../../../ui/node_modules/.pnpm/escape-string-regexp@5.0.0/node_modules/escape-string-regexp/index.js
  function escapeStringRegexp(string3) {
    if (typeof string3 !== "string") {
      throw new TypeError("Expected a string");
    }
    return string3.replace(/[|\\{}()[\]^$+*?.]/g, "\\$&").replace(/-/g, "\\x2d");
  }

  // ../../../ui/node_modules/.pnpm/unist-util-is@6.0.0/node_modules/unist-util-is/lib/index.js
  var convert = (
    // Note: overloads in JSDoc can’t yet use different `@template`s.
    /**
     * @type {(
     *   (<Condition extends string>(test: Condition) => (node: unknown, index?: number | null | undefined, parent?: Parent | null | undefined, context?: unknown) => node is Node & {type: Condition}) &
     *   (<Condition extends Props>(test: Condition) => (node: unknown, index?: number | null | undefined, parent?: Parent | null | undefined, context?: unknown) => node is Node & Condition) &
     *   (<Condition extends TestFunction>(test: Condition) => (node: unknown, index?: number | null | undefined, parent?: Parent | null | undefined, context?: unknown) => node is Node & Predicate<Condition, Node>) &
     *   ((test?: null | undefined) => (node?: unknown, index?: number | null | undefined, parent?: Parent | null | undefined, context?: unknown) => node is Node) &
     *   ((test?: Test) => Check)
     * )}
     */
    /**
     * @param {Test} [test]
     * @returns {Check}
     */
    (function(test) {
      if (test === null || test === void 0) {
        return ok2;
      }
      if (typeof test === "function") {
        return castFactory(test);
      }
      if (typeof test === "object") {
        return Array.isArray(test) ? anyFactory(test) : propsFactory(test);
      }
      if (typeof test === "string") {
        return typeFactory(test);
      }
      throw new Error("Expected function, string, or object as test");
    })
  );
  function anyFactory(tests) {
    const checks = [];
    let index2 = -1;
    while (++index2 < tests.length) {
      checks[index2] = convert(tests[index2]);
    }
    return castFactory(any);
    function any(...parameters) {
      let index3 = -1;
      while (++index3 < checks.length) {
        if (checks[index3].apply(this, parameters)) return true;
      }
      return false;
    }
  }
  function propsFactory(check) {
    const checkAsRecord = (
      /** @type {Record<string, unknown>} */
      check
    );
    return castFactory(all2);
    function all2(node2) {
      const nodeAsRecord = (
        /** @type {Record<string, unknown>} */
        /** @type {unknown} */
        node2
      );
      let key;
      for (key in check) {
        if (nodeAsRecord[key] !== checkAsRecord[key]) return false;
      }
      return true;
    }
  }
  function typeFactory(check) {
    return castFactory(type);
    function type(node2) {
      return node2 && node2.type === check;
    }
  }
  function castFactory(testFunction) {
    return check;
    function check(value, index2, parent) {
      return Boolean(
        looksLikeANode(value) && testFunction.call(
          this,
          value,
          typeof index2 === "number" ? index2 : void 0,
          parent || void 0
        )
      );
    }
  }
  function ok2() {
    return true;
  }
  function looksLikeANode(value) {
    return value !== null && typeof value === "object" && "type" in value;
  }

  // ../../../ui/node_modules/.pnpm/unist-util-visit-parents@6.0.1/node_modules/unist-util-visit-parents/lib/color.js
  function color2(d) {
    return d;
  }

  // ../../../ui/node_modules/.pnpm/unist-util-visit-parents@6.0.1/node_modules/unist-util-visit-parents/lib/index.js
  var empty3 = [];
  var CONTINUE2 = true;
  var EXIT2 = false;
  var SKIP2 = "skip";
  function visitParents(tree, test, visitor, reverse) {
    let check;
    if (typeof test === "function" && typeof visitor !== "function") {
      reverse = visitor;
      visitor = test;
    } else {
      check = test;
    }
    const is2 = convert(check);
    const step = reverse ? -1 : 1;
    factory(tree, void 0, [])();
    function factory(node2, index2, parents) {
      const value = (
        /** @type {Record<string, unknown>} */
        node2 && typeof node2 === "object" ? node2 : {}
      );
      if (typeof value.type === "string") {
        const name2 = (
          // `hast`
          typeof value.tagName === "string" ? value.tagName : (
            // `xast`
            typeof value.name === "string" ? value.name : void 0
          )
        );
        Object.defineProperty(visit3, "name", {
          value: "node (" + color2(node2.type + (name2 ? "<" + name2 + ">" : "")) + ")"
        });
      }
      return visit3;
      function visit3() {
        let result = empty3;
        let subresult;
        let offset2;
        let grandparents;
        if (!test || is2(node2, index2, parents[parents.length - 1] || void 0)) {
          result = toResult2(visitor(node2, parents));
          if (result[0] === EXIT2) {
            return result;
          }
        }
        if ("children" in node2 && node2.children) {
          const nodeAsParent = (
            /** @type {UnistParent} */
            node2
          );
          if (nodeAsParent.children && result[0] !== SKIP2) {
            offset2 = (reverse ? nodeAsParent.children.length : -1) + step;
            grandparents = parents.concat(nodeAsParent);
            while (offset2 > -1 && offset2 < nodeAsParent.children.length) {
              const child = nodeAsParent.children[offset2];
              subresult = factory(child, offset2, grandparents)();
              if (subresult[0] === EXIT2) {
                return subresult;
              }
              offset2 = typeof subresult[1] === "number" ? subresult[1] : offset2 + step;
            }
          }
        }
        return result;
      }
    }
  }
  function toResult2(value) {
    if (Array.isArray(value)) {
      return value;
    }
    if (typeof value === "number") {
      return [CONTINUE2, value];
    }
    return value === null || value === void 0 ? empty3 : [value];
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-find-and-replace@3.0.2/node_modules/mdast-util-find-and-replace/lib/index.js
  function findAndReplace(tree, list4, options) {
    const settings = options || {};
    const ignored = convert(settings.ignore || []);
    const pairs = toPairs(list4);
    let pairIndex = -1;
    while (++pairIndex < pairs.length) {
      visitParents(tree, "text", visitor);
    }
    function visitor(node2, parents) {
      let index2 = -1;
      let grandparent;
      while (++index2 < parents.length) {
        const parent = parents[index2];
        const siblings = grandparent ? grandparent.children : void 0;
        if (ignored(
          parent,
          siblings ? siblings.indexOf(parent) : void 0,
          grandparent
        )) {
          return;
        }
        grandparent = parent;
      }
      if (grandparent) {
        return handler(node2, parents);
      }
    }
    function handler(node2, parents) {
      const parent = parents[parents.length - 1];
      const find = pairs[pairIndex][0];
      const replace2 = pairs[pairIndex][1];
      let start2 = 0;
      const siblings = parent.children;
      const index2 = siblings.indexOf(node2);
      let change = false;
      let nodes = [];
      find.lastIndex = 0;
      let match = find.exec(node2.value);
      while (match) {
        const position2 = match.index;
        const matchObject = {
          index: match.index,
          input: match.input,
          stack: [...parents, node2]
        };
        let value = replace2(...match, matchObject);
        if (typeof value === "string") {
          value = value.length > 0 ? { type: "text", value } : void 0;
        }
        if (value === false) {
          find.lastIndex = position2 + 1;
        } else {
          if (start2 !== position2) {
            nodes.push({
              type: "text",
              value: node2.value.slice(start2, position2)
            });
          }
          if (Array.isArray(value)) {
            nodes.push(...value);
          } else if (value) {
            nodes.push(value);
          }
          start2 = position2 + match[0].length;
          change = true;
        }
        if (!find.global) {
          break;
        }
        match = find.exec(node2.value);
      }
      if (change) {
        if (start2 < node2.value.length) {
          nodes.push({ type: "text", value: node2.value.slice(start2) });
        }
        parent.children.splice(index2, 1, ...nodes);
      } else {
        nodes = [node2];
      }
      return index2 + nodes.length;
    }
  }
  function toPairs(tupleOrList) {
    const result = [];
    if (!Array.isArray(tupleOrList)) {
      throw new TypeError("Expected find and replace tuple or list of tuples");
    }
    const list4 = !tupleOrList[0] || Array.isArray(tupleOrList[0]) ? tupleOrList : [tupleOrList];
    let index2 = -1;
    while (++index2 < list4.length) {
      const tuple = list4[index2];
      result.push([toExpression(tuple[0]), toFunction(tuple[1])]);
    }
    return result;
  }
  function toExpression(find) {
    return typeof find === "string" ? new RegExp(escapeStringRegexp(find), "g") : find;
  }
  function toFunction(replace2) {
    return typeof replace2 === "function" ? replace2 : function() {
      return replace2;
    };
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm-autolink-literal@2.0.1/node_modules/mdast-util-gfm-autolink-literal/lib/index.js
  var inConstruct = "phrasing";
  var notInConstruct = ["autolink", "link", "image", "label"];
  function gfmAutolinkLiteralFromMarkdown() {
    return {
      transforms: [transformGfmAutolinkLiterals],
      enter: {
        literalAutolink: enterLiteralAutolink,
        literalAutolinkEmail: enterLiteralAutolinkValue,
        literalAutolinkHttp: enterLiteralAutolinkValue,
        literalAutolinkWww: enterLiteralAutolinkValue
      },
      exit: {
        literalAutolink: exitLiteralAutolink,
        literalAutolinkEmail: exitLiteralAutolinkEmail,
        literalAutolinkHttp: exitLiteralAutolinkHttp,
        literalAutolinkWww: exitLiteralAutolinkWww
      }
    };
  }
  function gfmAutolinkLiteralToMarkdown() {
    return {
      unsafe: [
        {
          character: "@",
          before: "[+\\-.\\w]",
          after: "[\\-.\\w]",
          inConstruct,
          notInConstruct
        },
        {
          character: ".",
          before: "[Ww]",
          after: "[\\-.\\w]",
          inConstruct,
          notInConstruct
        },
        {
          character: ":",
          before: "[ps]",
          after: "\\/",
          inConstruct,
          notInConstruct
        }
      ]
    };
  }
  function enterLiteralAutolink(token) {
    this.enter({ type: "link", title: null, url: "", children: [] }, token);
  }
  function enterLiteralAutolinkValue(token) {
    this.config.enter.autolinkProtocol.call(this, token);
  }
  function exitLiteralAutolinkHttp(token) {
    this.config.exit.autolinkProtocol.call(this, token);
  }
  function exitLiteralAutolinkWww(token) {
    this.config.exit.data.call(this, token);
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "link");
    node2.url = "http://" + this.sliceSerialize(token);
  }
  function exitLiteralAutolinkEmail(token) {
    this.config.exit.autolinkEmail.call(this, token);
  }
  function exitLiteralAutolink(token) {
    this.exit(token);
  }
  function transformGfmAutolinkLiterals(tree) {
    findAndReplace(
      tree,
      [
        [/(https?:\/\/|www(?=\.))([-.\w]+)([^ \t\r\n]*)/gi, findUrl],
        [/(?<=^|\s|[\u0021-\u002f\u003a-\u0040\u005b-\u0060\u007b-\u007e\u00a1-\u00a9\u00ab-\u00ac\u00ae-\u00b1\u00b4\u00b6-\u00b8\u00bb\u00bf\u00d7\u00f7\u02c2-\u02c5\u02d2-\u02df\u02e5-\u02eb\u02ed\u02ef-\u02ff\u0375\u037e\u0384-\u0385\u0387\u03f6\u0482\u055a-\u055f\u0589-\u058a\u058d-\u058f\u05be\u05c0\u05c3\u05c6\u05f3-\u05f4\u0606-\u060f\u061b\u061d-\u061f\u066a-\u066d\u06d4\u06de\u06e9\u06fd-\u06fe\u0700-\u070d\u07f6-\u07f9\u07fe-\u07ff\u0830-\u083e\u085e\u0888\u0964-\u0965\u0970\u09f2-\u09f3\u09fa-\u09fb\u09fd\u0a76\u0af0-\u0af1\u0b70\u0bf3-\u0bfa\u0c77\u0c7f\u0c84\u0d4f\u0d79\u0df4\u0e3f\u0e4f\u0e5a-\u0e5b\u0f01-\u0f17\u0f1a-\u0f1f\u0f34\u0f36\u0f38\u0f3a-\u0f3d\u0f85\u0fbe-\u0fc5\u0fc7-\u0fcc\u0fce-\u0fda\u104a-\u104f\u109e-\u109f\u10fb\u1360-\u1368\u1390-\u1399\u1400\u166d-\u166e\u169b-\u169c\u16eb-\u16ed\u1735-\u1736\u17d4-\u17d6\u17d8-\u17db\u1800-\u180a\u1940\u1944-\u1945\u19de-\u19ff\u1a1e-\u1a1f\u1aa0-\u1aa6\u1aa8-\u1aad\u1b4e-\u1b4f\u1b5a-\u1b6a\u1b74-\u1b7f\u1bfc-\u1bff\u1c3b-\u1c3f\u1c7e-\u1c7f\u1cc0-\u1cc7\u1cd3\u1fbd\u1fbf-\u1fc1\u1fcd-\u1fcf\u1fdd-\u1fdf\u1fed-\u1fef\u1ffd-\u1ffe\u2010-\u2027\u2030-\u205e\u207a-\u207e\u208a-\u208e\u20a0-\u20c1\u2100-\u2101\u2103-\u2106\u2108-\u2109\u2114\u2116-\u2118\u211e-\u2123\u2125\u2127\u2129\u212e\u213a-\u213b\u2140-\u2144\u214a-\u214d\u214f\u218a-\u218b\u2190-\u2429\u2440-\u244a\u249c-\u24e9\u2500-\u2775\u2794-\u2b73\u2b76-\u2bff\u2ce5-\u2cea\u2cf9-\u2cfc\u2cfe-\u2cff\u2d70\u2e00-\u2e2e\u2e30-\u2e5d\u2e80-\u2e99\u2e9b-\u2ef3\u2f00-\u2fd5\u2ff0-\u2fff\u3001-\u3004\u3008-\u3020\u3030\u3036-\u3037\u303d-\u303f\u309b-\u309c\u30a0\u30fb\u3190-\u3191\u3196-\u319f\u31c0-\u31e5\u31ef\u3200-\u321e\u322a-\u3247\u3250\u3260-\u327f\u328a-\u32b0\u32c0-\u33ff\u4dc0-\u4dff\ua490-\ua4c6\ua4fe-\ua4ff\ua60d-\ua60f\ua673\ua67e\ua6f2-\ua6f7\ua700-\ua716\ua720-\ua721\ua789-\ua78a\ua828-\ua82b\ua836-\ua839\ua874-\ua877\ua8ce-\ua8cf\ua8f8-\ua8fa\ua8fc\ua92e-\ua92f\ua95f\ua9c1-\ua9cd\ua9de-\ua9df\uaa5c-\uaa5f\uaa77-\uaa79\uaade-\uaadf\uaaf0-\uaaf1\uab5b\uab6a-\uab6b\uabeb\ufb29\ufbb2-\ufbd2\ufd3e-\ufd4f\ufd90-\ufd91\ufdc8-\ufdcf\ufdfc-\ufdff\ufe10-\ufe19\ufe30-\ufe52\ufe54-\ufe66\ufe68-\ufe6b\uff01-\uff0f\uff1a-\uff20\uff3b-\uff40\uff5b-\uff65\uffe0-\uffe6\uffe8-\uffee\ufffc-\ufffd])([-.\w+]+)@([-\w]+(?:\.[-\w]+)+)/g, findEmail]
      ],
      { ignore: ["link", "linkReference"] }
    );
  }
  function findUrl(_, protocol, domain2, path2, match) {
    let prefix = "";
    if (!previous2(match)) {
      return false;
    }
    if (/^w/i.test(protocol)) {
      domain2 = protocol + domain2;
      protocol = "";
      prefix = "http://";
    }
    if (!isCorrectDomain(domain2)) {
      return false;
    }
    const parts = splitUrl(domain2 + path2);
    if (!parts[0]) return false;
    const result = {
      type: "link",
      title: null,
      url: prefix + protocol + parts[0],
      children: [{ type: "text", value: protocol + parts[0] }]
    };
    if (parts[1]) {
      return [result, { type: "text", value: parts[1] }];
    }
    return result;
  }
  function findEmail(_, atext, label, match) {
    if (
      // Not an expected previous character.
      !previous2(match, true) || // Label ends in not allowed character.
      /[-\d_]$/.test(label)
    ) {
      return false;
    }
    return {
      type: "link",
      title: null,
      url: "mailto:" + atext + "@" + label,
      children: [{ type: "text", value: atext + "@" + label }]
    };
  }
  function isCorrectDomain(domain2) {
    const parts = domain2.split(".");
    if (parts.length < 2 || parts[parts.length - 1] && (/_/.test(parts[parts.length - 1]) || !/[a-zA-Z\d]/.test(parts[parts.length - 1])) || parts[parts.length - 2] && (/_/.test(parts[parts.length - 2]) || !/[a-zA-Z\d]/.test(parts[parts.length - 2]))) {
      return false;
    }
    return true;
  }
  function splitUrl(url) {
    const trailExec = /[!"&'),.:;<>?\]}]+$/.exec(url);
    if (!trailExec) {
      return [url, void 0];
    }
    url = url.slice(0, trailExec.index);
    let trail2 = trailExec[0];
    let closingParenIndex = trail2.indexOf(")");
    const openingParens = ccount(url, "(");
    let closingParens = ccount(url, ")");
    while (closingParenIndex !== -1 && openingParens > closingParens) {
      url += trail2.slice(0, closingParenIndex + 1);
      trail2 = trail2.slice(closingParenIndex + 1);
      closingParenIndex = trail2.indexOf(")");
      closingParens++;
    }
    return [url, trail2];
  }
  function previous2(match, email) {
    const code3 = match.input.charCodeAt(match.index - 1);
    return (match.index === 0 || unicodeWhitespace(code3) || unicodePunctuation(code3)) && // If it’s an email, the previous character should not be a slash.
    (!email || code3 !== 47);
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm-footnote@2.1.0/node_modules/mdast-util-gfm-footnote/lib/index.js
  footnoteReference.peek = footnoteReferencePeek;
  function enterFootnoteCallString() {
    this.buffer();
  }
  function enterFootnoteCall(token) {
    this.enter({ type: "footnoteReference", identifier: "", label: "" }, token);
  }
  function enterFootnoteDefinitionLabelString() {
    this.buffer();
  }
  function enterFootnoteDefinition(token) {
    this.enter(
      { type: "footnoteDefinition", identifier: "", label: "", children: [] },
      token
    );
  }
  function exitFootnoteCallString(token) {
    const label = this.resume();
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "footnoteReference");
    node2.identifier = normalizeIdentifier(
      this.sliceSerialize(token)
    ).toLowerCase();
    node2.label = label;
  }
  function exitFootnoteCall(token) {
    this.exit(token);
  }
  function exitFootnoteDefinitionLabelString(token) {
    const label = this.resume();
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "footnoteDefinition");
    node2.identifier = normalizeIdentifier(
      this.sliceSerialize(token)
    ).toLowerCase();
    node2.label = label;
  }
  function exitFootnoteDefinition(token) {
    this.exit(token);
  }
  function footnoteReferencePeek() {
    return "[";
  }
  function footnoteReference(node2, _, state, info) {
    const tracker = state.createTracker(info);
    let value = tracker.move("[^");
    const exit3 = state.enter("footnoteReference");
    const subexit = state.enter("reference");
    value += tracker.move(
      state.safe(state.associationId(node2), { after: "]", before: value })
    );
    subexit();
    exit3();
    value += tracker.move("]");
    return value;
  }
  function gfmFootnoteFromMarkdown() {
    return {
      enter: {
        gfmFootnoteCallString: enterFootnoteCallString,
        gfmFootnoteCall: enterFootnoteCall,
        gfmFootnoteDefinitionLabelString: enterFootnoteDefinitionLabelString,
        gfmFootnoteDefinition: enterFootnoteDefinition
      },
      exit: {
        gfmFootnoteCallString: exitFootnoteCallString,
        gfmFootnoteCall: exitFootnoteCall,
        gfmFootnoteDefinitionLabelString: exitFootnoteDefinitionLabelString,
        gfmFootnoteDefinition: exitFootnoteDefinition
      }
    };
  }
  function gfmFootnoteToMarkdown(options) {
    let firstLineBlank = false;
    if (options && options.firstLineBlank) {
      firstLineBlank = true;
    }
    return {
      handlers: { footnoteDefinition, footnoteReference },
      // This is on by default already.
      unsafe: [{ character: "[", inConstruct: ["label", "phrasing", "reference"] }]
    };
    function footnoteDefinition(node2, _, state, info) {
      const tracker = state.createTracker(info);
      let value = tracker.move("[^");
      const exit3 = state.enter("footnoteDefinition");
      const subexit = state.enter("label");
      value += tracker.move(
        state.safe(state.associationId(node2), { before: value, after: "]" })
      );
      subexit();
      value += tracker.move("]:");
      if (node2.children && node2.children.length > 0) {
        tracker.shift(4);
        value += tracker.move(
          (firstLineBlank ? "\n" : " ") + state.indentLines(
            state.containerFlow(node2, tracker.current()),
            firstLineBlank ? mapAll : mapExceptFirst
          )
        );
      }
      exit3();
      return value;
    }
  }
  function mapExceptFirst(line, index2, blank) {
    return index2 === 0 ? line : mapAll(line, index2, blank);
  }
  function mapAll(line, index2, blank) {
    return (blank ? "" : "    ") + line;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm-strikethrough@2.0.0/node_modules/mdast-util-gfm-strikethrough/lib/index.js
  var constructsWithoutStrikethrough = [
    "autolink",
    "destinationLiteral",
    "destinationRaw",
    "reference",
    "titleQuote",
    "titleApostrophe"
  ];
  handleDelete.peek = peekDelete;
  function gfmStrikethroughFromMarkdown() {
    return {
      canContainEols: ["delete"],
      enter: { strikethrough: enterStrikethrough },
      exit: { strikethrough: exitStrikethrough }
    };
  }
  function gfmStrikethroughToMarkdown() {
    return {
      unsafe: [
        {
          character: "~",
          inConstruct: "phrasing",
          notInConstruct: constructsWithoutStrikethrough
        }
      ],
      handlers: { delete: handleDelete }
    };
  }
  function enterStrikethrough(token) {
    this.enter({ type: "delete", children: [] }, token);
  }
  function exitStrikethrough(token) {
    this.exit(token);
  }
  function handleDelete(node2, _, state, info) {
    const tracker = state.createTracker(info);
    const exit3 = state.enter("strikethrough");
    let value = tracker.move("~~");
    value += state.containerPhrasing(node2, {
      ...tracker.current(),
      before: value,
      after: "~"
    });
    value += tracker.move("~~");
    exit3();
    return value;
  }
  function peekDelete() {
    return "~";
  }

  // ../../../ui/node_modules/.pnpm/markdown-table@3.0.4/node_modules/markdown-table/index.js
  function defaultStringLength(value) {
    return value.length;
  }
  function markdownTable(table, options) {
    const settings = options || {};
    const align = (settings.align || []).concat();
    const stringLength = settings.stringLength || defaultStringLength;
    const alignments = [];
    const cellMatrix = [];
    const sizeMatrix = [];
    const longestCellByColumn = [];
    let mostCellsPerRow = 0;
    let rowIndex = -1;
    while (++rowIndex < table.length) {
      const row2 = [];
      const sizes2 = [];
      let columnIndex2 = -1;
      if (table[rowIndex].length > mostCellsPerRow) {
        mostCellsPerRow = table[rowIndex].length;
      }
      while (++columnIndex2 < table[rowIndex].length) {
        const cell = serialize(table[rowIndex][columnIndex2]);
        if (settings.alignDelimiters !== false) {
          const size = stringLength(cell);
          sizes2[columnIndex2] = size;
          if (longestCellByColumn[columnIndex2] === void 0 || size > longestCellByColumn[columnIndex2]) {
            longestCellByColumn[columnIndex2] = size;
          }
        }
        row2.push(cell);
      }
      cellMatrix[rowIndex] = row2;
      sizeMatrix[rowIndex] = sizes2;
    }
    let columnIndex = -1;
    if (typeof align === "object" && "length" in align) {
      while (++columnIndex < mostCellsPerRow) {
        alignments[columnIndex] = toAlignment(align[columnIndex]);
      }
    } else {
      const code3 = toAlignment(align);
      while (++columnIndex < mostCellsPerRow) {
        alignments[columnIndex] = code3;
      }
    }
    columnIndex = -1;
    const row = [];
    const sizes = [];
    while (++columnIndex < mostCellsPerRow) {
      const code3 = alignments[columnIndex];
      let before = "";
      let after = "";
      if (code3 === 99) {
        before = ":";
        after = ":";
      } else if (code3 === 108) {
        before = ":";
      } else if (code3 === 114) {
        after = ":";
      }
      let size = settings.alignDelimiters === false ? 1 : Math.max(
        1,
        longestCellByColumn[columnIndex] - before.length - after.length
      );
      const cell = before + "-".repeat(size) + after;
      if (settings.alignDelimiters !== false) {
        size = before.length + size + after.length;
        if (size > longestCellByColumn[columnIndex]) {
          longestCellByColumn[columnIndex] = size;
        }
        sizes[columnIndex] = size;
      }
      row[columnIndex] = cell;
    }
    cellMatrix.splice(1, 0, row);
    sizeMatrix.splice(1, 0, sizes);
    rowIndex = -1;
    const lines = [];
    while (++rowIndex < cellMatrix.length) {
      const row2 = cellMatrix[rowIndex];
      const sizes2 = sizeMatrix[rowIndex];
      columnIndex = -1;
      const line = [];
      while (++columnIndex < mostCellsPerRow) {
        const cell = row2[columnIndex] || "";
        let before = "";
        let after = "";
        if (settings.alignDelimiters !== false) {
          const size = longestCellByColumn[columnIndex] - (sizes2[columnIndex] || 0);
          const code3 = alignments[columnIndex];
          if (code3 === 114) {
            before = " ".repeat(size);
          } else if (code3 === 99) {
            if (size % 2) {
              before = " ".repeat(size / 2 + 0.5);
              after = " ".repeat(size / 2 - 0.5);
            } else {
              before = " ".repeat(size / 2);
              after = before;
            }
          } else {
            after = " ".repeat(size);
          }
        }
        if (settings.delimiterStart !== false && !columnIndex) {
          line.push("|");
        }
        if (settings.padding !== false && // Don’t add the opening space if we’re not aligning and the cell is
        // empty: there will be a closing space.
        !(settings.alignDelimiters === false && cell === "") && (settings.delimiterStart !== false || columnIndex)) {
          line.push(" ");
        }
        if (settings.alignDelimiters !== false) {
          line.push(before);
        }
        line.push(cell);
        if (settings.alignDelimiters !== false) {
          line.push(after);
        }
        if (settings.padding !== false) {
          line.push(" ");
        }
        if (settings.delimiterEnd !== false || columnIndex !== mostCellsPerRow - 1) {
          line.push("|");
        }
      }
      lines.push(
        settings.delimiterEnd === false ? line.join("").replace(/ +$/, "") : line.join("")
      );
    }
    return lines.join("\n");
  }
  function serialize(value) {
    return value === null || value === void 0 ? "" : String(value);
  }
  function toAlignment(value) {
    const code3 = typeof value === "string" ? value.codePointAt(0) : 0;
    return code3 === 67 || code3 === 99 ? 99 : code3 === 76 || code3 === 108 ? 108 : code3 === 82 || code3 === 114 ? 114 : 0;
  }

  // ../../../ui/node_modules/.pnpm/zwitch@2.0.4/node_modules/zwitch/index.js
  var own4 = {}.hasOwnProperty;
  function zwitch(key, options) {
    const settings = options || {};
    function one2(value, ...parameters) {
      let fn = one2.invalid;
      const handlers = one2.handlers;
      if (value && own4.call(value, key)) {
        const id = String(value[key]);
        fn = own4.call(handlers, id) ? handlers[id] : one2.unknown;
      }
      if (fn) {
        return fn.call(this, value, ...parameters);
      }
    }
    one2.handlers = settings.handlers || {};
    one2.invalid = settings.invalid;
    one2.unknown = settings.unknown;
    return one2;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/configure.js
  var own5 = {}.hasOwnProperty;
  function configure(base, extension2) {
    let index2 = -1;
    let key;
    if (extension2.extensions) {
      while (++index2 < extension2.extensions.length) {
        configure(base, extension2.extensions[index2]);
      }
    }
    for (key in extension2) {
      if (own5.call(extension2, key)) {
        switch (key) {
          case "extensions": {
            break;
          }
          /* c8 ignore next 4 */
          case "unsafe": {
            list2(base[key], extension2[key]);
            break;
          }
          case "join": {
            list2(base[key], extension2[key]);
            break;
          }
          case "handlers": {
            map(base[key], extension2[key]);
            break;
          }
          default: {
            base.options[key] = extension2[key];
          }
        }
      }
    }
    return base;
  }
  function list2(left, right) {
    if (right) {
      left.push(...right);
    }
  }
  function map(left, right) {
    if (right) {
      Object.assign(left, right);
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/blockquote.js
  function blockquote(node2, _, state, info) {
    const exit3 = state.enter("blockquote");
    const tracker = state.createTracker(info);
    tracker.move("> ");
    tracker.shift(2);
    const value = state.indentLines(
      state.containerFlow(node2, tracker.current()),
      map2
    );
    exit3();
    return value;
  }
  function map2(line, _, blank) {
    return ">" + (blank ? "" : " ") + line;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/pattern-in-scope.js
  function patternInScope(stack, pattern) {
    return listInScope(stack, pattern.inConstruct, true) && !listInScope(stack, pattern.notInConstruct, false);
  }
  function listInScope(stack, list4, none) {
    if (typeof list4 === "string") {
      list4 = [list4];
    }
    if (!list4 || list4.length === 0) {
      return none;
    }
    let index2 = -1;
    while (++index2 < list4.length) {
      if (stack.includes(list4[index2])) {
        return true;
      }
    }
    return false;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/break.js
  function hardBreak(_, _1, state, info) {
    let index2 = -1;
    while (++index2 < state.unsafe.length) {
      if (state.unsafe[index2].character === "\n" && patternInScope(state.stack, state.unsafe[index2])) {
        return /[ \t]/.test(info.before) ? "" : " ";
      }
    }
    return "\\\n";
  }

  // ../../../ui/node_modules/.pnpm/longest-streak@3.1.0/node_modules/longest-streak/index.js
  function longestStreak(value, substring) {
    const source = String(value);
    let index2 = source.indexOf(substring);
    let expected = index2;
    let count3 = 0;
    let max = 0;
    if (typeof substring !== "string") {
      throw new TypeError("Expected substring");
    }
    while (index2 !== -1) {
      if (index2 === expected) {
        if (++count3 > max) {
          max = count3;
        }
      } else {
        count3 = 1;
      }
      expected = index2 + substring.length;
      index2 = source.indexOf(substring, expected);
    }
    return max;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/format-code-as-indented.js
  function formatCodeAsIndented(node2, state) {
    return Boolean(
      state.options.fences === false && node2.value && // If there’s no info…
      !node2.lang && // And there’s a non-whitespace character…
      /[^ \r\n]/.test(node2.value) && // And the value doesn’t start or end in a blank…
      !/^[\t ]*(?:[\r\n]|$)|(?:^|[\r\n])[\t ]*$/.test(node2.value)
    );
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-fence.js
  function checkFence(state) {
    const marker = state.options.fence || "`";
    if (marker !== "`" && marker !== "~") {
      throw new Error(
        "Cannot serialize code with `" + marker + "` for `options.fence`, expected `` ` `` or `~`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/code.js
  function code(node2, _, state, info) {
    const marker = checkFence(state);
    const raw = node2.value || "";
    const suffix = marker === "`" ? "GraveAccent" : "Tilde";
    if (formatCodeAsIndented(node2, state)) {
      const exit4 = state.enter("codeIndented");
      const value2 = state.indentLines(raw, map3);
      exit4();
      return value2;
    }
    const tracker = state.createTracker(info);
    const sequence = marker.repeat(Math.max(longestStreak(raw, marker) + 1, 3));
    const exit3 = state.enter("codeFenced");
    let value = tracker.move(sequence);
    if (node2.lang) {
      const subexit = state.enter(`codeFencedLang${suffix}`);
      value += tracker.move(
        state.safe(node2.lang, {
          before: value,
          after: " ",
          encode: ["`"],
          ...tracker.current()
        })
      );
      subexit();
    }
    if (node2.lang && node2.meta) {
      const subexit = state.enter(`codeFencedMeta${suffix}`);
      value += tracker.move(" ");
      value += tracker.move(
        state.safe(node2.meta, {
          before: value,
          after: "\n",
          encode: ["`"],
          ...tracker.current()
        })
      );
      subexit();
    }
    value += tracker.move("\n");
    if (raw) {
      value += tracker.move(raw + "\n");
    }
    value += tracker.move(sequence);
    exit3();
    return value;
  }
  function map3(line, _, blank) {
    return (blank ? "" : "    ") + line;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-quote.js
  function checkQuote(state) {
    const marker = state.options.quote || '"';
    if (marker !== '"' && marker !== "'") {
      throw new Error(
        "Cannot serialize title with `" + marker + "` for `options.quote`, expected `\"`, or `'`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/definition.js
  function definition2(node2, _, state, info) {
    const quote = checkQuote(state);
    const suffix = quote === '"' ? "Quote" : "Apostrophe";
    const exit3 = state.enter("definition");
    let subexit = state.enter("label");
    const tracker = state.createTracker(info);
    let value = tracker.move("[");
    value += tracker.move(
      state.safe(state.associationId(node2), {
        before: value,
        after: "]",
        ...tracker.current()
      })
    );
    value += tracker.move("]: ");
    subexit();
    if (
      // If there’s no url, or…
      !node2.url || // If there are control characters or whitespace.
      /[\0- \u007F]/.test(node2.url)
    ) {
      subexit = state.enter("destinationLiteral");
      value += tracker.move("<");
      value += tracker.move(
        state.safe(node2.url, { before: value, after: ">", ...tracker.current() })
      );
      value += tracker.move(">");
    } else {
      subexit = state.enter("destinationRaw");
      value += tracker.move(
        state.safe(node2.url, {
          before: value,
          after: node2.title ? " " : "\n",
          ...tracker.current()
        })
      );
    }
    subexit();
    if (node2.title) {
      subexit = state.enter(`title${suffix}`);
      value += tracker.move(" " + quote);
      value += tracker.move(
        state.safe(node2.title, {
          before: value,
          after: quote,
          ...tracker.current()
        })
      );
      value += tracker.move(quote);
      subexit();
    }
    exit3();
    return value;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-emphasis.js
  function checkEmphasis(state) {
    const marker = state.options.emphasis || "*";
    if (marker !== "*" && marker !== "_") {
      throw new Error(
        "Cannot serialize emphasis with `" + marker + "` for `options.emphasis`, expected `*`, or `_`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/encode-character-reference.js
  function encodeCharacterReference(code3) {
    return "&#x" + code3.toString(16).toUpperCase() + ";";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/encode-info.js
  function encodeInfo(outside, inside, marker) {
    const outsideKind = classifyCharacter(outside);
    const insideKind = classifyCharacter(inside);
    if (outsideKind === void 0) {
      return insideKind === void 0 ? (
        // Letter inside:
        // we have to encode *both* letters for `_` as it is looser.
        // it already forms for `*` (and GFMs `~`).
        marker === "_" ? { inside: true, outside: true } : { inside: false, outside: false }
      ) : insideKind === 1 ? (
        // Whitespace inside: encode both (letter, whitespace).
        { inside: true, outside: true }
      ) : (
        // Punctuation inside: encode outer (letter)
        { inside: false, outside: true }
      );
    }
    if (outsideKind === 1) {
      return insideKind === void 0 ? (
        // Letter inside: already forms.
        { inside: false, outside: false }
      ) : insideKind === 1 ? (
        // Whitespace inside: encode both (whitespace).
        { inside: true, outside: true }
      ) : (
        // Punctuation inside: already forms.
        { inside: false, outside: false }
      );
    }
    return insideKind === void 0 ? (
      // Letter inside: already forms.
      { inside: false, outside: false }
    ) : insideKind === 1 ? (
      // Whitespace inside: encode inner (whitespace).
      { inside: true, outside: false }
    ) : (
      // Punctuation inside: already forms.
      { inside: false, outside: false }
    );
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/emphasis.js
  emphasis.peek = emphasisPeek;
  function emphasis(node2, _, state, info) {
    const marker = checkEmphasis(state);
    const exit3 = state.enter("emphasis");
    const tracker = state.createTracker(info);
    const before = tracker.move(marker);
    let between2 = tracker.move(
      state.containerPhrasing(node2, {
        after: marker,
        before,
        ...tracker.current()
      })
    );
    const betweenHead = between2.charCodeAt(0);
    const open = encodeInfo(
      info.before.charCodeAt(info.before.length - 1),
      betweenHead,
      marker
    );
    if (open.inside) {
      between2 = encodeCharacterReference(betweenHead) + between2.slice(1);
    }
    const betweenTail = between2.charCodeAt(between2.length - 1);
    const close = encodeInfo(info.after.charCodeAt(0), betweenTail, marker);
    if (close.inside) {
      between2 = between2.slice(0, -1) + encodeCharacterReference(betweenTail);
    }
    const after = tracker.move(marker);
    exit3();
    state.attentionEncodeSurroundingInfo = {
      after: close.outside,
      before: open.outside
    };
    return before + between2 + after;
  }
  function emphasisPeek(_, _1, state) {
    return state.options.emphasis || "*";
  }

  // ../../../ui/node_modules/.pnpm/unist-util-visit@5.1.0/node_modules/unist-util-visit/lib/index.js
  function visit2(tree, testOrVisitor, visitorOrReverse, maybeReverse) {
    let reverse;
    let test;
    let visitor;
    if (typeof testOrVisitor === "function" && typeof visitorOrReverse !== "function") {
      test = void 0;
      visitor = testOrVisitor;
      reverse = visitorOrReverse;
    } else {
      test = testOrVisitor;
      visitor = visitorOrReverse;
      reverse = maybeReverse;
    }
    visitParents(tree, test, overload, reverse);
    function overload(node2, parents) {
      const parent = parents[parents.length - 1];
      const index2 = parent ? parent.children.indexOf(node2) : void 0;
      return visitor(node2, index2, parent);
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-string@4.0.0/node_modules/mdast-util-to-string/lib/index.js
  var emptyOptions3 = {};
  function toString2(value, options) {
    const settings = options || emptyOptions3;
    const includeImageAlt = typeof settings.includeImageAlt === "boolean" ? settings.includeImageAlt : true;
    const includeHtml = typeof settings.includeHtml === "boolean" ? settings.includeHtml : true;
    return one(value, includeImageAlt, includeHtml);
  }
  function one(value, includeImageAlt, includeHtml) {
    if (node(value)) {
      if ("value" in value) {
        return value.type === "html" && !includeHtml ? "" : value.value;
      }
      if (includeImageAlt && "alt" in value && value.alt) {
        return value.alt;
      }
      if ("children" in value) {
        return all(value.children, includeImageAlt, includeHtml);
      }
    }
    if (Array.isArray(value)) {
      return all(value, includeImageAlt, includeHtml);
    }
    return "";
  }
  function all(values, includeImageAlt, includeHtml) {
    const result = [];
    let index2 = -1;
    while (++index2 < values.length) {
      result[index2] = one(values[index2], includeImageAlt, includeHtml);
    }
    return result.join("");
  }
  function node(value) {
    return Boolean(value && typeof value === "object");
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/format-heading-as-setext.js
  function formatHeadingAsSetext(node2, state) {
    let literalWithBreak = false;
    visit2(node2, function(node3) {
      if ("value" in node3 && /\r?\n|\r/.test(node3.value) || node3.type === "break") {
        literalWithBreak = true;
        return EXIT2;
      }
    });
    return Boolean(
      (!node2.depth || node2.depth < 3) && toString2(node2) && (state.options.setext || literalWithBreak)
    );
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/heading.js
  function heading(node2, _, state, info) {
    const rank = Math.max(Math.min(6, node2.depth || 1), 1);
    const tracker = state.createTracker(info);
    if (formatHeadingAsSetext(node2, state)) {
      const exit4 = state.enter("headingSetext");
      const subexit2 = state.enter("phrasing");
      const value2 = state.containerPhrasing(node2, {
        ...tracker.current(),
        before: "\n",
        after: "\n"
      });
      subexit2();
      exit4();
      return value2 + "\n" + (rank === 1 ? "=" : "-").repeat(
        // The whole size…
        value2.length - // Minus the position of the character after the last EOL (or
        // 0 if there is none)…
        (Math.max(value2.lastIndexOf("\r"), value2.lastIndexOf("\n")) + 1)
      );
    }
    const sequence = "#".repeat(rank);
    const exit3 = state.enter("headingAtx");
    const subexit = state.enter("phrasing");
    tracker.move(sequence + " ");
    let value = state.containerPhrasing(node2, {
      before: "# ",
      after: "\n",
      ...tracker.current()
    });
    if (/^[\t ]/.test(value)) {
      value = encodeCharacterReference(value.charCodeAt(0)) + value.slice(1);
    }
    value = value ? sequence + " " + value : sequence;
    if (state.options.closeAtx) {
      value += " " + sequence;
    }
    subexit();
    exit3();
    return value;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/html.js
  html.peek = htmlPeek;
  function html(node2) {
    return node2.value || "";
  }
  function htmlPeek() {
    return "<";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/image.js
  image.peek = imagePeek;
  function image(node2, _, state, info) {
    const quote = checkQuote(state);
    const suffix = quote === '"' ? "Quote" : "Apostrophe";
    const exit3 = state.enter("image");
    let subexit = state.enter("label");
    const tracker = state.createTracker(info);
    let value = tracker.move("![");
    value += tracker.move(
      state.safe(node2.alt, { before: value, after: "]", ...tracker.current() })
    );
    value += tracker.move("](");
    subexit();
    if (
      // If there’s no url but there is a title…
      !node2.url && node2.title || // If there are control characters or whitespace.
      /[\0- \u007F]/.test(node2.url)
    ) {
      subexit = state.enter("destinationLiteral");
      value += tracker.move("<");
      value += tracker.move(
        state.safe(node2.url, { before: value, after: ">", ...tracker.current() })
      );
      value += tracker.move(">");
    } else {
      subexit = state.enter("destinationRaw");
      value += tracker.move(
        state.safe(node2.url, {
          before: value,
          after: node2.title ? " " : ")",
          ...tracker.current()
        })
      );
    }
    subexit();
    if (node2.title) {
      subexit = state.enter(`title${suffix}`);
      value += tracker.move(" " + quote);
      value += tracker.move(
        state.safe(node2.title, {
          before: value,
          after: quote,
          ...tracker.current()
        })
      );
      value += tracker.move(quote);
      subexit();
    }
    value += tracker.move(")");
    exit3();
    return value;
  }
  function imagePeek() {
    return "!";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/image-reference.js
  imageReference.peek = imageReferencePeek;
  function imageReference(node2, _, state, info) {
    const type = node2.referenceType;
    const exit3 = state.enter("imageReference");
    let subexit = state.enter("label");
    const tracker = state.createTracker(info);
    let value = tracker.move("![");
    const alt = state.safe(node2.alt, {
      before: value,
      after: "]",
      ...tracker.current()
    });
    value += tracker.move(alt + "][");
    subexit();
    const stack = state.stack;
    state.stack = [];
    subexit = state.enter("reference");
    const reference = state.safe(state.associationId(node2), {
      before: value,
      after: "]",
      ...tracker.current()
    });
    subexit();
    state.stack = stack;
    exit3();
    if (type === "full" || !alt || alt !== reference) {
      value += tracker.move(reference + "]");
    } else if (type === "shortcut") {
      value = value.slice(0, -1);
    } else {
      value += tracker.move("]");
    }
    return value;
  }
  function imageReferencePeek() {
    return "!";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/inline-code.js
  inlineCode.peek = inlineCodePeek;
  function inlineCode(node2, _, state) {
    let value = node2.value || "";
    let sequence = "`";
    let index2 = -1;
    while (new RegExp("(^|[^`])" + sequence + "([^`]|$)").test(value)) {
      sequence += "`";
    }
    if (/[^ \r\n]/.test(value) && (/^[ \r\n]/.test(value) && /[ \r\n]$/.test(value) || /^`|`$/.test(value))) {
      value = " " + value + " ";
    }
    while (++index2 < state.unsafe.length) {
      const pattern = state.unsafe[index2];
      const expression = state.compilePattern(pattern);
      let match;
      if (!pattern.atBreak) continue;
      while (match = expression.exec(value)) {
        let position2 = match.index;
        if (value.charCodeAt(position2) === 10 && value.charCodeAt(position2 - 1) === 13) {
          position2--;
        }
        value = value.slice(0, position2) + " " + value.slice(match.index + 1);
      }
    }
    return sequence + value + sequence;
  }
  function inlineCodePeek() {
    return "`";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/format-link-as-autolink.js
  function formatLinkAsAutolink(node2, state) {
    const raw = toString2(node2);
    return Boolean(
      !state.options.resourceLink && // If there’s a url…
      node2.url && // And there’s a no title…
      !node2.title && // And the content of `node` is a single text node…
      node2.children && node2.children.length === 1 && node2.children[0].type === "text" && // And if the url is the same as the content…
      (raw === node2.url || "mailto:" + raw === node2.url) && // And that starts w/ a protocol…
      /^[a-z][a-z+.-]+:/i.test(node2.url) && // And that doesn’t contain ASCII control codes (character escapes and
      // references don’t work), space, or angle brackets…
      !/[\0- <>\u007F]/.test(node2.url)
    );
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/link.js
  link.peek = linkPeek;
  function link(node2, _, state, info) {
    const quote = checkQuote(state);
    const suffix = quote === '"' ? "Quote" : "Apostrophe";
    const tracker = state.createTracker(info);
    let exit3;
    let subexit;
    if (formatLinkAsAutolink(node2, state)) {
      const stack = state.stack;
      state.stack = [];
      exit3 = state.enter("autolink");
      let value2 = tracker.move("<");
      value2 += tracker.move(
        state.containerPhrasing(node2, {
          before: value2,
          after: ">",
          ...tracker.current()
        })
      );
      value2 += tracker.move(">");
      exit3();
      state.stack = stack;
      return value2;
    }
    exit3 = state.enter("link");
    subexit = state.enter("label");
    let value = tracker.move("[");
    value += tracker.move(
      state.containerPhrasing(node2, {
        before: value,
        after: "](",
        ...tracker.current()
      })
    );
    value += tracker.move("](");
    subexit();
    if (
      // If there’s no url but there is a title…
      !node2.url && node2.title || // If there are control characters or whitespace.
      /[\0- \u007F]/.test(node2.url)
    ) {
      subexit = state.enter("destinationLiteral");
      value += tracker.move("<");
      value += tracker.move(
        state.safe(node2.url, { before: value, after: ">", ...tracker.current() })
      );
      value += tracker.move(">");
    } else {
      subexit = state.enter("destinationRaw");
      value += tracker.move(
        state.safe(node2.url, {
          before: value,
          after: node2.title ? " " : ")",
          ...tracker.current()
        })
      );
    }
    subexit();
    if (node2.title) {
      subexit = state.enter(`title${suffix}`);
      value += tracker.move(" " + quote);
      value += tracker.move(
        state.safe(node2.title, {
          before: value,
          after: quote,
          ...tracker.current()
        })
      );
      value += tracker.move(quote);
      subexit();
    }
    value += tracker.move(")");
    exit3();
    return value;
  }
  function linkPeek(node2, _, state) {
    return formatLinkAsAutolink(node2, state) ? "<" : "[";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/link-reference.js
  linkReference.peek = linkReferencePeek;
  function linkReference(node2, _, state, info) {
    const type = node2.referenceType;
    const exit3 = state.enter("linkReference");
    let subexit = state.enter("label");
    const tracker = state.createTracker(info);
    let value = tracker.move("[");
    const text5 = state.containerPhrasing(node2, {
      before: value,
      after: "]",
      ...tracker.current()
    });
    value += tracker.move(text5 + "][");
    subexit();
    const stack = state.stack;
    state.stack = [];
    subexit = state.enter("reference");
    const reference = state.safe(state.associationId(node2), {
      before: value,
      after: "]",
      ...tracker.current()
    });
    subexit();
    state.stack = stack;
    exit3();
    if (type === "full" || !text5 || text5 !== reference) {
      value += tracker.move(reference + "]");
    } else if (type === "shortcut") {
      value = value.slice(0, -1);
    } else {
      value += tracker.move("]");
    }
    return value;
  }
  function linkReferencePeek() {
    return "[";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-bullet.js
  function checkBullet(state) {
    const marker = state.options.bullet || "*";
    if (marker !== "*" && marker !== "+" && marker !== "-") {
      throw new Error(
        "Cannot serialize items with `" + marker + "` for `options.bullet`, expected `*`, `+`, or `-`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-bullet-other.js
  function checkBulletOther(state) {
    const bullet = checkBullet(state);
    const bulletOther = state.options.bulletOther;
    if (!bulletOther) {
      return bullet === "*" ? "-" : "*";
    }
    if (bulletOther !== "*" && bulletOther !== "+" && bulletOther !== "-") {
      throw new Error(
        "Cannot serialize items with `" + bulletOther + "` for `options.bulletOther`, expected `*`, `+`, or `-`"
      );
    }
    if (bulletOther === bullet) {
      throw new Error(
        "Expected `bullet` (`" + bullet + "`) and `bulletOther` (`" + bulletOther + "`) to be different"
      );
    }
    return bulletOther;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-bullet-ordered.js
  function checkBulletOrdered(state) {
    const marker = state.options.bulletOrdered || ".";
    if (marker !== "." && marker !== ")") {
      throw new Error(
        "Cannot serialize items with `" + marker + "` for `options.bulletOrdered`, expected `.` or `)`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-rule.js
  function checkRule(state) {
    const marker = state.options.rule || "*";
    if (marker !== "*" && marker !== "-" && marker !== "_") {
      throw new Error(
        "Cannot serialize rules with `" + marker + "` for `options.rule`, expected `*`, `-`, or `_`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/list.js
  function list3(node2, parent, state, info) {
    const exit3 = state.enter("list");
    const bulletCurrent = state.bulletCurrent;
    let bullet = node2.ordered ? checkBulletOrdered(state) : checkBullet(state);
    const bulletOther = node2.ordered ? bullet === "." ? ")" : "." : checkBulletOther(state);
    let useDifferentMarker = parent && state.bulletLastUsed ? bullet === state.bulletLastUsed : false;
    if (!node2.ordered) {
      const firstListItem = node2.children ? node2.children[0] : void 0;
      if (
        // Bullet could be used as a thematic break marker:
        (bullet === "*" || bullet === "-") && // Empty first list item:
        firstListItem && (!firstListItem.children || !firstListItem.children[0]) && // Directly in two other list items:
        state.stack[state.stack.length - 1] === "list" && state.stack[state.stack.length - 2] === "listItem" && state.stack[state.stack.length - 3] === "list" && state.stack[state.stack.length - 4] === "listItem" && // That are each the first child.
        state.indexStack[state.indexStack.length - 1] === 0 && state.indexStack[state.indexStack.length - 2] === 0 && state.indexStack[state.indexStack.length - 3] === 0
      ) {
        useDifferentMarker = true;
      }
      if (checkRule(state) === bullet && firstListItem) {
        let index2 = -1;
        while (++index2 < node2.children.length) {
          const item = node2.children[index2];
          if (item && item.type === "listItem" && item.children && item.children[0] && item.children[0].type === "thematicBreak") {
            useDifferentMarker = true;
            break;
          }
        }
      }
    }
    if (useDifferentMarker) {
      bullet = bulletOther;
    }
    state.bulletCurrent = bullet;
    const value = state.containerFlow(node2, info);
    state.bulletLastUsed = bullet;
    state.bulletCurrent = bulletCurrent;
    exit3();
    return value;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-list-item-indent.js
  function checkListItemIndent(state) {
    const style = state.options.listItemIndent || "one";
    if (style !== "tab" && style !== "one" && style !== "mixed") {
      throw new Error(
        "Cannot serialize items with `" + style + "` for `options.listItemIndent`, expected `tab`, `one`, or `mixed`"
      );
    }
    return style;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/list-item.js
  function listItem(node2, parent, state, info) {
    const listItemIndent = checkListItemIndent(state);
    let bullet = state.bulletCurrent || checkBullet(state);
    if (parent && parent.type === "list" && parent.ordered) {
      bullet = (typeof parent.start === "number" && parent.start > -1 ? parent.start : 1) + (state.options.incrementListMarker === false ? 0 : parent.children.indexOf(node2)) + bullet;
    }
    let size = bullet.length + 1;
    if (listItemIndent === "tab" || listItemIndent === "mixed" && (parent && parent.type === "list" && parent.spread || node2.spread)) {
      size = Math.ceil(size / 4) * 4;
    }
    const tracker = state.createTracker(info);
    tracker.move(bullet + " ".repeat(size - bullet.length));
    tracker.shift(size);
    const exit3 = state.enter("listItem");
    const value = state.indentLines(
      state.containerFlow(node2, tracker.current()),
      map4
    );
    exit3();
    return value;
    function map4(line, index2, blank) {
      if (index2) {
        return (blank ? "" : " ".repeat(size)) + line;
      }
      return (blank ? bullet : bullet + " ".repeat(size - bullet.length)) + line;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/paragraph.js
  function paragraph(node2, _, state, info) {
    const exit3 = state.enter("paragraph");
    const subexit = state.enter("phrasing");
    const value = state.containerPhrasing(node2, info);
    subexit();
    exit3();
    return value;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-phrasing@4.1.0/node_modules/mdast-util-phrasing/lib/index.js
  var phrasing = (
    /** @type {(node?: unknown) => node is Exclude<PhrasingContent, Html>} */
    convert([
      "break",
      "delete",
      "emphasis",
      // To do: next major: removed since footnotes were added to GFM.
      "footnote",
      "footnoteReference",
      "image",
      "imageReference",
      "inlineCode",
      // Enabled by `mdast-util-math`:
      "inlineMath",
      "link",
      "linkReference",
      // Enabled by `mdast-util-mdx`:
      "mdxJsxTextElement",
      // Enabled by `mdast-util-mdx`:
      "mdxTextExpression",
      "strong",
      "text",
      // Enabled by `mdast-util-directive`:
      "textDirective"
    ])
  );

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/root.js
  function root(node2, _, state, info) {
    const hasPhrasing = node2.children.some(function(d) {
      return phrasing(d);
    });
    const container = hasPhrasing ? state.containerPhrasing : state.containerFlow;
    return container.call(state, node2, info);
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-strong.js
  function checkStrong(state) {
    const marker = state.options.strong || "*";
    if (marker !== "*" && marker !== "_") {
      throw new Error(
        "Cannot serialize strong with `" + marker + "` for `options.strong`, expected `*`, or `_`"
      );
    }
    return marker;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/strong.js
  strong.peek = strongPeek;
  function strong(node2, _, state, info) {
    const marker = checkStrong(state);
    const exit3 = state.enter("strong");
    const tracker = state.createTracker(info);
    const before = tracker.move(marker + marker);
    let between2 = tracker.move(
      state.containerPhrasing(node2, {
        after: marker,
        before,
        ...tracker.current()
      })
    );
    const betweenHead = between2.charCodeAt(0);
    const open = encodeInfo(
      info.before.charCodeAt(info.before.length - 1),
      betweenHead,
      marker
    );
    if (open.inside) {
      between2 = encodeCharacterReference(betweenHead) + between2.slice(1);
    }
    const betweenTail = between2.charCodeAt(between2.length - 1);
    const close = encodeInfo(info.after.charCodeAt(0), betweenTail, marker);
    if (close.inside) {
      between2 = between2.slice(0, -1) + encodeCharacterReference(betweenTail);
    }
    const after = tracker.move(marker + marker);
    exit3();
    state.attentionEncodeSurroundingInfo = {
      after: close.outside,
      before: open.outside
    };
    return before + between2 + after;
  }
  function strongPeek(_, _1, state) {
    return state.options.strong || "*";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/text.js
  function text(node2, _, state, info) {
    return state.safe(node2.value, info);
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/check-rule-repetition.js
  function checkRuleRepetition(state) {
    const repetition = state.options.ruleRepetition || 3;
    if (repetition < 3) {
      throw new Error(
        "Cannot serialize rules with repetition `" + repetition + "` for `options.ruleRepetition`, expected `3` or more"
      );
    }
    return repetition;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/thematic-break.js
  function thematicBreak2(_, _1, state) {
    const value = (checkRule(state) + (state.options.ruleSpaces ? " " : "")).repeat(checkRuleRepetition(state));
    return state.options.ruleSpaces ? value.slice(0, -1) : value;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/handle/index.js
  var handle = {
    blockquote,
    break: hardBreak,
    code,
    definition: definition2,
    emphasis,
    hardBreak,
    heading,
    html,
    image,
    imageReference,
    inlineCode,
    link,
    linkReference,
    list: list3,
    listItem,
    paragraph,
    root,
    strong,
    text,
    thematicBreak: thematicBreak2
  };

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/join.js
  var join2 = [joinDefaults];
  function joinDefaults(left, right, parent, state) {
    if (right.type === "code" && formatCodeAsIndented(right, state) && (left.type === "list" || left.type === right.type && formatCodeAsIndented(left, state))) {
      return false;
    }
    if ("spread" in parent && typeof parent.spread === "boolean") {
      if (left.type === "paragraph" && // Two paragraphs.
      (left.type === right.type || right.type === "definition" || // Paragraph followed by a setext heading.
      right.type === "heading" && formatHeadingAsSetext(right, state))) {
        return;
      }
      return parent.spread ? 1 : 0;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/unsafe.js
  var fullPhrasingSpans = [
    "autolink",
    "destinationLiteral",
    "destinationRaw",
    "reference",
    "titleQuote",
    "titleApostrophe"
  ];
  var unsafe = [
    { character: "	", after: "[\\r\\n]", inConstruct: "phrasing" },
    { character: "	", before: "[\\r\\n]", inConstruct: "phrasing" },
    {
      character: "	",
      inConstruct: ["codeFencedLangGraveAccent", "codeFencedLangTilde"]
    },
    {
      character: "\r",
      inConstruct: [
        "codeFencedLangGraveAccent",
        "codeFencedLangTilde",
        "codeFencedMetaGraveAccent",
        "codeFencedMetaTilde",
        "destinationLiteral",
        "headingAtx"
      ]
    },
    {
      character: "\n",
      inConstruct: [
        "codeFencedLangGraveAccent",
        "codeFencedLangTilde",
        "codeFencedMetaGraveAccent",
        "codeFencedMetaTilde",
        "destinationLiteral",
        "headingAtx"
      ]
    },
    { character: " ", after: "[\\r\\n]", inConstruct: "phrasing" },
    { character: " ", before: "[\\r\\n]", inConstruct: "phrasing" },
    {
      character: " ",
      inConstruct: ["codeFencedLangGraveAccent", "codeFencedLangTilde"]
    },
    // An exclamation mark can start an image, if it is followed by a link or
    // a link reference.
    {
      character: "!",
      after: "\\[",
      inConstruct: "phrasing",
      notInConstruct: fullPhrasingSpans
    },
    // A quote can break out of a title.
    { character: '"', inConstruct: "titleQuote" },
    // A number sign could start an ATX heading if it starts a line.
    { atBreak: true, character: "#" },
    { character: "#", inConstruct: "headingAtx", after: "(?:[\r\n]|$)" },
    // Dollar sign and percentage are not used in markdown.
    // An ampersand could start a character reference.
    { character: "&", after: "[#A-Za-z]", inConstruct: "phrasing" },
    // An apostrophe can break out of a title.
    { character: "'", inConstruct: "titleApostrophe" },
    // A left paren could break out of a destination raw.
    { character: "(", inConstruct: "destinationRaw" },
    // A left paren followed by `]` could make something into a link or image.
    {
      before: "\\]",
      character: "(",
      inConstruct: "phrasing",
      notInConstruct: fullPhrasingSpans
    },
    // A right paren could start a list item or break out of a destination
    // raw.
    { atBreak: true, before: "\\d+", character: ")" },
    { character: ")", inConstruct: "destinationRaw" },
    // An asterisk can start thematic breaks, list items, emphasis, strong.
    { atBreak: true, character: "*", after: "(?:[ 	\r\n*])" },
    { character: "*", inConstruct: "phrasing", notInConstruct: fullPhrasingSpans },
    // A plus sign could start a list item.
    { atBreak: true, character: "+", after: "(?:[ 	\r\n])" },
    // A dash can start thematic breaks, list items, and setext heading
    // underlines.
    { atBreak: true, character: "-", after: "(?:[ 	\r\n-])" },
    // A dot could start a list item.
    { atBreak: true, before: "\\d+", character: ".", after: "(?:[ 	\r\n]|$)" },
    // Slash, colon, and semicolon are not used in markdown for constructs.
    // A less than can start html (flow or text) or an autolink.
    // HTML could start with an exclamation mark (declaration, cdata, comment),
    // slash (closing tag), question mark (instruction), or a letter (tag).
    // An autolink also starts with a letter.
    // Finally, it could break out of a destination literal.
    { atBreak: true, character: "<", after: "[!/?A-Za-z]" },
    {
      character: "<",
      after: "[!/?A-Za-z]",
      inConstruct: "phrasing",
      notInConstruct: fullPhrasingSpans
    },
    { character: "<", inConstruct: "destinationLiteral" },
    // An equals to can start setext heading underlines.
    { atBreak: true, character: "=" },
    // A greater than can start block quotes and it can break out of a
    // destination literal.
    { atBreak: true, character: ">" },
    { character: ">", inConstruct: "destinationLiteral" },
    // Question mark and at sign are not used in markdown for constructs.
    // A left bracket can start definitions, references, labels,
    { atBreak: true, character: "[" },
    { character: "[", inConstruct: "phrasing", notInConstruct: fullPhrasingSpans },
    { character: "[", inConstruct: ["label", "reference"] },
    // A backslash can start an escape (when followed by punctuation) or a
    // hard break (when followed by an eol).
    // Note: typical escapes are handled in `safe`!
    { character: "\\", after: "[\\r\\n]", inConstruct: "phrasing" },
    // A right bracket can exit labels.
    { character: "]", inConstruct: ["label", "reference"] },
    // Caret is not used in markdown for constructs.
    // An underscore can start emphasis, strong, or a thematic break.
    { atBreak: true, character: "_" },
    { character: "_", inConstruct: "phrasing", notInConstruct: fullPhrasingSpans },
    // A grave accent can start code (fenced or text), or it can break out of
    // a grave accent code fence.
    { atBreak: true, character: "`" },
    {
      character: "`",
      inConstruct: ["codeFencedLangGraveAccent", "codeFencedMetaGraveAccent"]
    },
    { character: "`", inConstruct: "phrasing", notInConstruct: fullPhrasingSpans },
    // Left brace, vertical bar, right brace are not used in markdown for
    // constructs.
    // A tilde can start code (fenced).
    { atBreak: true, character: "~" }
  ];

  // ../../../ui/node_modules/.pnpm/micromark-util-decode-numeric-character-reference@2.0.2/node_modules/micromark-util-decode-numeric-character-reference/index.js
  function decodeNumericCharacterReference(value, base) {
    const code3 = Number.parseInt(value, base);
    if (
      // C0 except for HT, LF, FF, CR, space.
      code3 < 9 || code3 === 11 || code3 > 13 && code3 < 32 || // Control character (DEL) of C0, and C1 controls.
      code3 > 126 && code3 < 160 || // Lone high surrogates and low surrogates.
      code3 > 55295 && code3 < 57344 || // Noncharacters.
      code3 > 64975 && code3 < 65008 || /* eslint-disable no-bitwise */
      (code3 & 65535) === 65535 || (code3 & 65535) === 65534 || /* eslint-enable no-bitwise */
      // Out of range
      code3 > 1114111
    ) {
      return "�";
    }
    return String.fromCodePoint(code3);
  }

  // ../../../ui/node_modules/.pnpm/micromark-util-decode-string@2.0.1/node_modules/micromark-util-decode-string/index.js
  var characterEscapeOrReference = /\\([!-/:-@[-`{-~])|&(#(?:\d{1,7}|x[\da-f]{1,6})|[\da-z]{1,31});/gi;
  function decodeString(value) {
    return value.replace(characterEscapeOrReference, decode);
  }
  function decode($0, $1, $2) {
    if ($1) {
      return $1;
    }
    const head = $2.charCodeAt(0);
    if (head === 35) {
      const head2 = $2.charCodeAt(1);
      const hex = head2 === 120 || head2 === 88;
      return decodeNumericCharacterReference($2.slice(hex ? 2 : 1), hex ? 16 : 10);
    }
    return decodeNamedCharacterReference($2) || $0;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/association.js
  function association(node2) {
    if (node2.label || !node2.identifier) {
      return node2.label || "";
    }
    return decodeString(node2.identifier);
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/compile-pattern.js
  function compilePattern(pattern) {
    if (!pattern._compiled) {
      const before = (pattern.atBreak ? "[\\r\\n][\\t ]*" : "") + (pattern.before ? "(?:" + pattern.before + ")" : "");
      pattern._compiled = new RegExp(
        (before ? "(" + before + ")" : "") + (/[|\\{}()[\]^$+*?.-]/.test(pattern.character) ? "\\" : "") + pattern.character + (pattern.after ? "(?:" + pattern.after + ")" : ""),
        "g"
      );
    }
    return pattern._compiled;
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/container-phrasing.js
  function containerPhrasing(parent, state, info) {
    const indexStack = state.indexStack;
    const children = parent.children || [];
    const results = [];
    let index2 = -1;
    let before = info.before;
    let encodeAfter;
    indexStack.push(-1);
    let tracker = state.createTracker(info);
    while (++index2 < children.length) {
      const child = children[index2];
      let after;
      indexStack[indexStack.length - 1] = index2;
      if (index2 + 1 < children.length) {
        let handle2 = state.handle.handlers[children[index2 + 1].type];
        if (handle2 && handle2.peek) handle2 = handle2.peek;
        after = handle2 ? handle2(children[index2 + 1], parent, state, {
          before: "",
          after: "",
          ...tracker.current()
        }).charAt(0) : "";
      } else {
        after = info.after;
      }
      if (results.length > 0 && (before === "\r" || before === "\n") && child.type === "html") {
        results[results.length - 1] = results[results.length - 1].replace(
          /(\r?\n|\r)$/,
          " "
        );
        before = " ";
        tracker = state.createTracker(info);
        tracker.move(results.join(""));
      }
      let value = state.handle(child, parent, state, {
        ...tracker.current(),
        after,
        before
      });
      if (encodeAfter && encodeAfter === value.slice(0, 1)) {
        value = encodeCharacterReference(encodeAfter.charCodeAt(0)) + value.slice(1);
      }
      const encodingInfo = state.attentionEncodeSurroundingInfo;
      state.attentionEncodeSurroundingInfo = void 0;
      encodeAfter = void 0;
      if (encodingInfo) {
        if (results.length > 0 && encodingInfo.before && before === results[results.length - 1].slice(-1)) {
          results[results.length - 1] = results[results.length - 1].slice(0, -1) + encodeCharacterReference(before.charCodeAt(0));
        }
        if (encodingInfo.after) encodeAfter = after;
      }
      tracker.move(value);
      results.push(value);
      before = value.slice(-1);
    }
    indexStack.pop();
    return results.join("");
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/container-flow.js
  function containerFlow2(parent, state, info) {
    const indexStack = state.indexStack;
    const children = parent.children || [];
    const tracker = state.createTracker(info);
    const results = [];
    let index2 = -1;
    indexStack.push(-1);
    while (++index2 < children.length) {
      const child = children[index2];
      indexStack[indexStack.length - 1] = index2;
      results.push(
        tracker.move(
          state.handle(child, parent, state, {
            before: "\n",
            after: "\n",
            ...tracker.current()
          })
        )
      );
      if (child.type !== "list") {
        state.bulletLastUsed = void 0;
      }
      if (index2 < children.length - 1) {
        results.push(
          tracker.move(between(child, children[index2 + 1], parent, state))
        );
      }
    }
    indexStack.pop();
    return results.join("");
  }
  function between(left, right, parent, state) {
    let index2 = state.join.length;
    while (index2--) {
      const result = state.join[index2](left, right, parent, state);
      if (result === true || result === 1) {
        break;
      }
      if (typeof result === "number") {
        return "\n".repeat(1 + result);
      }
      if (result === false) {
        return "\n\n<!---->\n\n";
      }
    }
    return "\n\n";
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/indent-lines.js
  var eol = /\r?\n|\r/g;
  function indentLines(value, map4) {
    const result = [];
    let start2 = 0;
    let line = 0;
    let match;
    while (match = eol.exec(value)) {
      one2(value.slice(start2, match.index));
      result.push(match[0]);
      start2 = match.index + match[0].length;
      line++;
    }
    one2(value.slice(start2));
    return result.join("");
    function one2(value2) {
      result.push(map4(value2, line, !value2));
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/safe.js
  function safe(state, input, config) {
    const value = (config.before || "") + (input || "") + (config.after || "");
    const positions = [];
    const result = [];
    const infos = {};
    let index2 = -1;
    while (++index2 < state.unsafe.length) {
      const pattern = state.unsafe[index2];
      if (!patternInScope(state.stack, pattern)) {
        continue;
      }
      const expression = state.compilePattern(pattern);
      let match;
      while (match = expression.exec(value)) {
        const before = "before" in pattern || Boolean(pattern.atBreak);
        const after = "after" in pattern;
        const position2 = match.index + (before ? match[1].length : 0);
        if (positions.includes(position2)) {
          if (infos[position2].before && !before) {
            infos[position2].before = false;
          }
          if (infos[position2].after && !after) {
            infos[position2].after = false;
          }
        } else {
          positions.push(position2);
          infos[position2] = { before, after };
        }
      }
    }
    positions.sort(numerical);
    let start2 = config.before ? config.before.length : 0;
    const end = value.length - (config.after ? config.after.length : 0);
    index2 = -1;
    while (++index2 < positions.length) {
      const position2 = positions[index2];
      if (position2 < start2 || position2 >= end) {
        continue;
      }
      if (position2 + 1 < end && positions[index2 + 1] === position2 + 1 && infos[position2].after && !infos[position2 + 1].before && !infos[position2 + 1].after || positions[index2 - 1] === position2 - 1 && infos[position2].before && !infos[position2 - 1].before && !infos[position2 - 1].after) {
        continue;
      }
      if (start2 !== position2) {
        result.push(escapeBackslashes(value.slice(start2, position2), "\\"));
      }
      start2 = position2;
      if (/[!-/:-@[-`{-~]/.test(value.charAt(position2)) && (!config.encode || !config.encode.includes(value.charAt(position2)))) {
        result.push("\\");
      } else {
        result.push(encodeCharacterReference(value.charCodeAt(position2)));
        start2++;
      }
    }
    result.push(escapeBackslashes(value.slice(start2, end), config.after));
    return result.join("");
  }
  function numerical(a, b) {
    return a - b;
  }
  function escapeBackslashes(value, after) {
    const expression = /\\(?=[!-/:-@[-`{-~])/g;
    const positions = [];
    const results = [];
    const whole = value + after;
    let index2 = -1;
    let start2 = 0;
    let match;
    while (match = expression.exec(whole)) {
      positions.push(match.index);
    }
    while (++index2 < positions.length) {
      if (start2 !== positions[index2]) {
        results.push(value.slice(start2, positions[index2]));
      }
      results.push("\\");
      start2 = positions[index2];
    }
    results.push(value.slice(start2));
    return results.join("");
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/util/track.js
  function track(config) {
    const options = config || {};
    const now = options.now || {};
    let lineShift = options.lineShift || 0;
    let line = now.line || 1;
    let column = now.column || 1;
    return { move, current: current2, shift };
    function current2() {
      return { now: { line, column }, lineShift };
    }
    function shift(value) {
      lineShift += value;
    }
    function move(input) {
      const value = input || "";
      const chunks = value.split(/\r?\n|\r/g);
      const tail = chunks[chunks.length - 1];
      line += chunks.length - 1;
      column = chunks.length === 1 ? column + tail.length : 1 + tail.length + lineShift;
      return value;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-to-markdown@2.1.2/node_modules/mdast-util-to-markdown/lib/index.js
  function toMarkdown(tree, options) {
    const settings = options || {};
    const state = {
      associationId: association,
      containerPhrasing: containerPhrasingBound,
      containerFlow: containerFlowBound,
      createTracker: track,
      compilePattern,
      enter,
      // @ts-expect-error: GFM / frontmatter are typed in `mdast` but not defined
      // here.
      handlers: { ...handle },
      // @ts-expect-error: add `handle` in a second.
      handle: void 0,
      indentLines,
      indexStack: [],
      join: [...join2],
      options: {},
      safe: safeBound,
      stack: [],
      unsafe: [...unsafe]
    };
    configure(state, settings);
    if (state.options.tightDefinitions) {
      state.join.push(joinDefinition);
    }
    state.handle = zwitch("type", {
      invalid,
      unknown,
      handlers: state.handlers
    });
    let result = state.handle(tree, void 0, state, {
      before: "\n",
      after: "\n",
      now: { line: 1, column: 1 },
      lineShift: 0
    });
    if (result && result.charCodeAt(result.length - 1) !== 10 && result.charCodeAt(result.length - 1) !== 13) {
      result += "\n";
    }
    return result;
    function enter(name2) {
      state.stack.push(name2);
      return exit3;
      function exit3() {
        state.stack.pop();
      }
    }
  }
  function invalid(value) {
    throw new Error("Cannot handle value `" + value + "`, expected node");
  }
  function unknown(value) {
    const node2 = (
      /** @type {Nodes} */
      value
    );
    throw new Error("Cannot handle unknown node `" + node2.type + "`");
  }
  function joinDefinition(left, right) {
    if (left.type === "definition" && left.type === right.type) {
      return 0;
    }
  }
  function containerPhrasingBound(parent, info) {
    return containerPhrasing(parent, this, info);
  }
  function containerFlowBound(parent, info) {
    return containerFlow2(parent, this, info);
  }
  function safeBound(value, config) {
    return safe(this, value, config);
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm-table@2.0.0/node_modules/mdast-util-gfm-table/lib/index.js
  function gfmTableFromMarkdown() {
    return {
      enter: {
        table: enterTable,
        tableData: enterCell,
        tableHeader: enterCell,
        tableRow: enterRow
      },
      exit: {
        codeText: exitCodeText,
        table: exitTable,
        tableData: exit2,
        tableHeader: exit2,
        tableRow: exit2
      }
    };
  }
  function enterTable(token) {
    const align = token._align;
    ok(align, "expected `_align` on table");
    this.enter(
      {
        type: "table",
        align: align.map(function(d) {
          return d === "none" ? null : d;
        }),
        children: []
      },
      token
    );
    this.data.inTable = true;
  }
  function exitTable(token) {
    this.exit(token);
    this.data.inTable = void 0;
  }
  function enterRow(token) {
    this.enter({ type: "tableRow", children: [] }, token);
  }
  function exit2(token) {
    this.exit(token);
  }
  function enterCell(token) {
    this.enter({ type: "tableCell", children: [] }, token);
  }
  function exitCodeText(token) {
    let value = this.resume();
    if (this.data.inTable) {
      value = value.replace(/\\([\\|])/g, replace);
    }
    const node2 = this.stack[this.stack.length - 1];
    ok(node2.type === "inlineCode");
    node2.value = value;
    this.exit(token);
  }
  function replace($0, $1) {
    return $1 === "|" ? $1 : $0;
  }
  function gfmTableToMarkdown(options) {
    const settings = options || {};
    const padding = settings.tableCellPadding;
    const alignDelimiters = settings.tablePipeAlign;
    const stringLength = settings.stringLength;
    const around = padding ? " " : "|";
    return {
      unsafe: [
        { character: "\r", inConstruct: "tableCell" },
        { character: "\n", inConstruct: "tableCell" },
        // A pipe, when followed by a tab or space (padding), or a dash or colon
        // (unpadded delimiter row), could result in a table.
        { atBreak: true, character: "|", after: "[	 :-]" },
        // A pipe in a cell must be encoded.
        { character: "|", inConstruct: "tableCell" },
        // A colon must be followed by a dash, in which case it could start a
        // delimiter row.
        { atBreak: true, character: ":", after: "-" },
        // A delimiter row can also start with a dash, when followed by more
        // dashes, a colon, or a pipe.
        // This is a stricter version than the built in check for lists, thematic
        // breaks, and setex heading underlines though:
        // <https://github.com/syntax-tree/mdast-util-to-markdown/blob/51a2038/lib/unsafe.js#L57>
        { atBreak: true, character: "-", after: "[:|-]" }
      ],
      handlers: {
        inlineCode: inlineCodeWithTable,
        table: handleTable,
        tableCell: handleTableCell,
        tableRow: handleTableRow
      }
    };
    function handleTable(node2, _, state, info) {
      return serializeData(handleTableAsData(node2, state, info), node2.align);
    }
    function handleTableRow(node2, _, state, info) {
      const row = handleTableRowAsData(node2, state, info);
      const value = serializeData([row]);
      return value.slice(0, value.indexOf("\n"));
    }
    function handleTableCell(node2, _, state, info) {
      const exit3 = state.enter("tableCell");
      const subexit = state.enter("phrasing");
      const value = state.containerPhrasing(node2, {
        ...info,
        before: around,
        after: around
      });
      subexit();
      exit3();
      return value;
    }
    function serializeData(matrix, align) {
      return markdownTable(matrix, {
        align,
        // @ts-expect-error: `markdown-table` types should support `null`.
        alignDelimiters,
        // @ts-expect-error: `markdown-table` types should support `null`.
        padding,
        // @ts-expect-error: `markdown-table` types should support `null`.
        stringLength
      });
    }
    function handleTableAsData(node2, state, info) {
      const children = node2.children;
      let index2 = -1;
      const result = [];
      const subexit = state.enter("table");
      while (++index2 < children.length) {
        result[index2] = handleTableRowAsData(children[index2], state, info);
      }
      subexit();
      return result;
    }
    function handleTableRowAsData(node2, state, info) {
      const children = node2.children;
      let index2 = -1;
      const result = [];
      const subexit = state.enter("tableRow");
      while (++index2 < children.length) {
        result[index2] = handleTableCell(children[index2], node2, state, info);
      }
      subexit();
      return result;
    }
    function inlineCodeWithTable(node2, parent, state) {
      let value = handle.inlineCode(node2, parent, state);
      if (state.stack.includes("tableCell")) {
        value = value.replace(/\|/g, "\\$&");
      }
      return value;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm-task-list-item@2.0.0/node_modules/mdast-util-gfm-task-list-item/lib/index.js
  function gfmTaskListItemFromMarkdown() {
    return {
      exit: {
        taskListCheckValueChecked: exitCheck,
        taskListCheckValueUnchecked: exitCheck,
        paragraph: exitParagraphWithTaskListItem
      }
    };
  }
  function gfmTaskListItemToMarkdown() {
    return {
      unsafe: [{ atBreak: true, character: "-", after: "[:|-]" }],
      handlers: { listItem: listItemWithTaskListItem }
    };
  }
  function exitCheck(token) {
    const node2 = this.stack[this.stack.length - 2];
    ok(node2.type === "listItem");
    node2.checked = token.type === "taskListCheckValueChecked";
  }
  function exitParagraphWithTaskListItem(token) {
    const parent = this.stack[this.stack.length - 2];
    if (parent && parent.type === "listItem" && typeof parent.checked === "boolean") {
      const node2 = this.stack[this.stack.length - 1];
      ok(node2.type === "paragraph");
      const head = node2.children[0];
      if (head && head.type === "text") {
        const siblings = parent.children;
        let index2 = -1;
        let firstParaghraph;
        while (++index2 < siblings.length) {
          const sibling2 = siblings[index2];
          if (sibling2.type === "paragraph") {
            firstParaghraph = sibling2;
            break;
          }
        }
        if (firstParaghraph === node2) {
          head.value = head.value.slice(1);
          if (head.value.length === 0) {
            node2.children.shift();
          } else if (node2.position && head.position && typeof head.position.start.offset === "number") {
            head.position.start.column++;
            head.position.start.offset++;
            node2.position.start = Object.assign({}, head.position.start);
          }
        }
      }
    }
    this.exit(token);
  }
  function listItemWithTaskListItem(node2, parent, state, info) {
    const head = node2.children[0];
    const checkable = typeof node2.checked === "boolean" && head && head.type === "paragraph";
    const checkbox = "[" + (node2.checked ? "x" : " ") + "] ";
    const tracker = state.createTracker(info);
    if (checkable) {
      tracker.move(checkbox);
    }
    let value = handle.listItem(node2, parent, state, {
      ...info,
      ...tracker.current()
    });
    if (checkable) {
      value = value.replace(/^(?:[*+-]|\d+\.)([\r\n]| {1,3})/, check);
    }
    return value;
    function check($0) {
      return $0 + checkbox;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-gfm@3.1.0/node_modules/mdast-util-gfm/lib/index.js
  function gfmFromMarkdown() {
    return [
      gfmAutolinkLiteralFromMarkdown(),
      gfmFootnoteFromMarkdown(),
      gfmStrikethroughFromMarkdown(),
      gfmTableFromMarkdown(),
      gfmTaskListItemFromMarkdown()
    ];
  }
  function gfmToMarkdown(options) {
    return {
      extensions: [
        gfmAutolinkLiteralToMarkdown(),
        gfmFootnoteToMarkdown(options),
        gfmStrikethroughToMarkdown(),
        gfmTableToMarkdown(options),
        gfmTaskListItemToMarkdown()
      ]
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-autolink-literal@2.1.0/node_modules/micromark-extension-gfm-autolink-literal/lib/syntax.js
  var wwwPrefix = {
    tokenize: tokenizeWwwPrefix,
    partial: true
  };
  var domain = {
    tokenize: tokenizeDomain,
    partial: true
  };
  var path = {
    tokenize: tokenizePath,
    partial: true
  };
  var trail = {
    tokenize: tokenizeTrail,
    partial: true
  };
  var emailDomainDotTrail = {
    tokenize: tokenizeEmailDomainDotTrail,
    partial: true
  };
  var wwwAutolink = {
    name: "wwwAutolink",
    tokenize: tokenizeWwwAutolink,
    previous: previousWww
  };
  var protocolAutolink = {
    name: "protocolAutolink",
    tokenize: tokenizeProtocolAutolink,
    previous: previousProtocol
  };
  var emailAutolink = {
    name: "emailAutolink",
    tokenize: tokenizeEmailAutolink,
    previous: previousEmail
  };
  var text2 = {};
  function gfmAutolinkLiteral() {
    return {
      text: text2
    };
  }
  var code2 = 48;
  while (code2 < 123) {
    text2[code2] = emailAutolink;
    code2++;
    if (code2 === 58) code2 = 65;
    else if (code2 === 91) code2 = 97;
  }
  text2[43] = emailAutolink;
  text2[45] = emailAutolink;
  text2[46] = emailAutolink;
  text2[95] = emailAutolink;
  text2[72] = [emailAutolink, protocolAutolink];
  text2[104] = [emailAutolink, protocolAutolink];
  text2[87] = [emailAutolink, wwwAutolink];
  text2[119] = [emailAutolink, wwwAutolink];
  function tokenizeEmailAutolink(effects, ok3, nok) {
    const self2 = this;
    let dot;
    let data2;
    return start2;
    function start2(code3) {
      if (!gfmAtext(code3) || !previousEmail.call(self2, self2.previous) || previousUnbalanced(self2.events)) {
        return nok(code3);
      }
      effects.enter("literalAutolink");
      effects.enter("literalAutolinkEmail");
      return atext(code3);
    }
    function atext(code3) {
      if (gfmAtext(code3)) {
        effects.consume(code3);
        return atext;
      }
      if (code3 === 64) {
        effects.consume(code3);
        return emailDomain;
      }
      return nok(code3);
    }
    function emailDomain(code3) {
      if (code3 === 46) {
        return effects.check(emailDomainDotTrail, emailDomainAfter, emailDomainDot)(code3);
      }
      if (code3 === 45 || code3 === 95 || asciiAlphanumeric(code3)) {
        data2 = true;
        effects.consume(code3);
        return emailDomain;
      }
      return emailDomainAfter(code3);
    }
    function emailDomainDot(code3) {
      effects.consume(code3);
      dot = true;
      return emailDomain;
    }
    function emailDomainAfter(code3) {
      if (data2 && dot && asciiAlpha(self2.previous)) {
        effects.exit("literalAutolinkEmail");
        effects.exit("literalAutolink");
        return ok3(code3);
      }
      return nok(code3);
    }
  }
  function tokenizeWwwAutolink(effects, ok3, nok) {
    const self2 = this;
    return wwwStart;
    function wwwStart(code3) {
      if (code3 !== 87 && code3 !== 119 || !previousWww.call(self2, self2.previous) || previousUnbalanced(self2.events)) {
        return nok(code3);
      }
      effects.enter("literalAutolink");
      effects.enter("literalAutolinkWww");
      return effects.check(wwwPrefix, effects.attempt(domain, effects.attempt(path, wwwAfter), nok), nok)(code3);
    }
    function wwwAfter(code3) {
      effects.exit("literalAutolinkWww");
      effects.exit("literalAutolink");
      return ok3(code3);
    }
  }
  function tokenizeProtocolAutolink(effects, ok3, nok) {
    const self2 = this;
    let buffer = "";
    let seen = false;
    return protocolStart;
    function protocolStart(code3) {
      if ((code3 === 72 || code3 === 104) && previousProtocol.call(self2, self2.previous) && !previousUnbalanced(self2.events)) {
        effects.enter("literalAutolink");
        effects.enter("literalAutolinkHttp");
        buffer += String.fromCodePoint(code3);
        effects.consume(code3);
        return protocolPrefixInside;
      }
      return nok(code3);
    }
    function protocolPrefixInside(code3) {
      if (asciiAlpha(code3) && buffer.length < 5) {
        buffer += String.fromCodePoint(code3);
        effects.consume(code3);
        return protocolPrefixInside;
      }
      if (code3 === 58) {
        const protocol = buffer.toLowerCase();
        if (protocol === "http" || protocol === "https") {
          effects.consume(code3);
          return protocolSlashesInside;
        }
      }
      return nok(code3);
    }
    function protocolSlashesInside(code3) {
      if (code3 === 47) {
        effects.consume(code3);
        if (seen) {
          return afterProtocol;
        }
        seen = true;
        return protocolSlashesInside;
      }
      return nok(code3);
    }
    function afterProtocol(code3) {
      return code3 === null || asciiControl(code3) || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3) || unicodePunctuation(code3) ? nok(code3) : effects.attempt(domain, effects.attempt(path, protocolAfter), nok)(code3);
    }
    function protocolAfter(code3) {
      effects.exit("literalAutolinkHttp");
      effects.exit("literalAutolink");
      return ok3(code3);
    }
  }
  function tokenizeWwwPrefix(effects, ok3, nok) {
    let size = 0;
    return wwwPrefixInside;
    function wwwPrefixInside(code3) {
      if ((code3 === 87 || code3 === 119) && size < 3) {
        size++;
        effects.consume(code3);
        return wwwPrefixInside;
      }
      if (code3 === 46 && size === 3) {
        effects.consume(code3);
        return wwwPrefixAfter;
      }
      return nok(code3);
    }
    function wwwPrefixAfter(code3) {
      return code3 === null ? nok(code3) : ok3(code3);
    }
  }
  function tokenizeDomain(effects, ok3, nok) {
    let underscoreInLastSegment;
    let underscoreInLastLastSegment;
    let seen;
    return domainInside;
    function domainInside(code3) {
      if (code3 === 46 || code3 === 95) {
        return effects.check(trail, domainAfter, domainAtPunctuation)(code3);
      }
      if (code3 === null || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3) || code3 !== 45 && unicodePunctuation(code3)) {
        return domainAfter(code3);
      }
      seen = true;
      effects.consume(code3);
      return domainInside;
    }
    function domainAtPunctuation(code3) {
      if (code3 === 95) {
        underscoreInLastSegment = true;
      } else {
        underscoreInLastLastSegment = underscoreInLastSegment;
        underscoreInLastSegment = void 0;
      }
      effects.consume(code3);
      return domainInside;
    }
    function domainAfter(code3) {
      if (underscoreInLastLastSegment || underscoreInLastSegment || !seen) {
        return nok(code3);
      }
      return ok3(code3);
    }
  }
  function tokenizePath(effects, ok3) {
    let sizeOpen = 0;
    let sizeClose = 0;
    return pathInside;
    function pathInside(code3) {
      if (code3 === 40) {
        sizeOpen++;
        effects.consume(code3);
        return pathInside;
      }
      if (code3 === 41 && sizeClose < sizeOpen) {
        return pathAtPunctuation(code3);
      }
      if (code3 === 33 || code3 === 34 || code3 === 38 || code3 === 39 || code3 === 41 || code3 === 42 || code3 === 44 || code3 === 46 || code3 === 58 || code3 === 59 || code3 === 60 || code3 === 63 || code3 === 93 || code3 === 95 || code3 === 126) {
        return effects.check(trail, ok3, pathAtPunctuation)(code3);
      }
      if (code3 === null || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        return ok3(code3);
      }
      effects.consume(code3);
      return pathInside;
    }
    function pathAtPunctuation(code3) {
      if (code3 === 41) {
        sizeClose++;
      }
      effects.consume(code3);
      return pathInside;
    }
  }
  function tokenizeTrail(effects, ok3, nok) {
    return trail2;
    function trail2(code3) {
      if (code3 === 33 || code3 === 34 || code3 === 39 || code3 === 41 || code3 === 42 || code3 === 44 || code3 === 46 || code3 === 58 || code3 === 59 || code3 === 63 || code3 === 95 || code3 === 126) {
        effects.consume(code3);
        return trail2;
      }
      if (code3 === 38) {
        effects.consume(code3);
        return trailCharacterReferenceStart;
      }
      if (code3 === 93) {
        effects.consume(code3);
        return trailBracketAfter;
      }
      if (
        // `<` is an end.
        code3 === 60 || // So is whitespace.
        code3 === null || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)
      ) {
        return ok3(code3);
      }
      return nok(code3);
    }
    function trailBracketAfter(code3) {
      if (code3 === null || code3 === 40 || code3 === 91 || markdownLineEndingOrSpace(code3) || unicodeWhitespace(code3)) {
        return ok3(code3);
      }
      return trail2(code3);
    }
    function trailCharacterReferenceStart(code3) {
      return asciiAlpha(code3) ? trailCharacterReferenceInside(code3) : nok(code3);
    }
    function trailCharacterReferenceInside(code3) {
      if (code3 === 59) {
        effects.consume(code3);
        return trail2;
      }
      if (asciiAlpha(code3)) {
        effects.consume(code3);
        return trailCharacterReferenceInside;
      }
      return nok(code3);
    }
  }
  function tokenizeEmailDomainDotTrail(effects, ok3, nok) {
    return start2;
    function start2(code3) {
      effects.consume(code3);
      return after;
    }
    function after(code3) {
      return asciiAlphanumeric(code3) ? nok(code3) : ok3(code3);
    }
  }
  function previousWww(code3) {
    return code3 === null || code3 === 40 || code3 === 42 || code3 === 95 || code3 === 91 || code3 === 93 || code3 === 126 || markdownLineEndingOrSpace(code3);
  }
  function previousProtocol(code3) {
    return !asciiAlpha(code3);
  }
  function previousEmail(code3) {
    return !(code3 === 47 || gfmAtext(code3));
  }
  function gfmAtext(code3) {
    return code3 === 43 || code3 === 45 || code3 === 46 || code3 === 95 || asciiAlphanumeric(code3);
  }
  function previousUnbalanced(events) {
    let index2 = events.length;
    let result = false;
    while (index2--) {
      const token = events[index2][1];
      if ((token.type === "labelLink" || token.type === "labelImage") && !token._balanced) {
        result = true;
        break;
      }
      if (token._gfmAutolinkLiteralWalkedInto) {
        result = false;
        break;
      }
    }
    if (events.length > 0 && !result) {
      events[events.length - 1][1]._gfmAutolinkLiteralWalkedInto = true;
    }
    return result;
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-footnote@2.1.0/node_modules/micromark-extension-gfm-footnote/lib/syntax.js
  var indent2 = {
    tokenize: tokenizeIndent2,
    partial: true
  };
  function gfmFootnote() {
    return {
      document: {
        [91]: {
          name: "gfmFootnoteDefinition",
          tokenize: tokenizeDefinitionStart,
          continuation: {
            tokenize: tokenizeDefinitionContinuation
          },
          exit: gfmFootnoteDefinitionEnd
        }
      },
      text: {
        [91]: {
          name: "gfmFootnoteCall",
          tokenize: tokenizeGfmFootnoteCall
        },
        [93]: {
          name: "gfmPotentialFootnoteCall",
          add: "after",
          tokenize: tokenizePotentialGfmFootnoteCall,
          resolveTo: resolveToPotentialGfmFootnoteCall
        }
      }
    };
  }
  function tokenizePotentialGfmFootnoteCall(effects, ok3, nok) {
    const self2 = this;
    let index2 = self2.events.length;
    const defined = self2.parser.gfmFootnotes || (self2.parser.gfmFootnotes = []);
    let labelStart;
    while (index2--) {
      const token = self2.events[index2][1];
      if (token.type === "labelImage") {
        labelStart = token;
        break;
      }
      if (token.type === "gfmFootnoteCall" || token.type === "labelLink" || token.type === "label" || token.type === "image" || token.type === "link") {
        break;
      }
    }
    return start2;
    function start2(code3) {
      if (!labelStart || !labelStart._balanced) {
        return nok(code3);
      }
      const id = normalizeIdentifier(self2.sliceSerialize({
        start: labelStart.end,
        end: self2.now()
      }));
      if (id.codePointAt(0) !== 94 || !defined.includes(id.slice(1))) {
        return nok(code3);
      }
      effects.enter("gfmFootnoteCallLabelMarker");
      effects.consume(code3);
      effects.exit("gfmFootnoteCallLabelMarker");
      return ok3(code3);
    }
  }
  function resolveToPotentialGfmFootnoteCall(events, context) {
    let index2 = events.length;
    let labelStart;
    while (index2--) {
      if (events[index2][1].type === "labelImage" && events[index2][0] === "enter") {
        labelStart = events[index2][1];
        break;
      }
    }
    events[index2 + 1][1].type = "data";
    events[index2 + 3][1].type = "gfmFootnoteCallLabelMarker";
    const call = {
      type: "gfmFootnoteCall",
      start: Object.assign({}, events[index2 + 3][1].start),
      end: Object.assign({}, events[events.length - 1][1].end)
    };
    const marker = {
      type: "gfmFootnoteCallMarker",
      start: Object.assign({}, events[index2 + 3][1].end),
      end: Object.assign({}, events[index2 + 3][1].end)
    };
    marker.end.column++;
    marker.end.offset++;
    marker.end._bufferIndex++;
    const string3 = {
      type: "gfmFootnoteCallString",
      start: Object.assign({}, marker.end),
      end: Object.assign({}, events[events.length - 1][1].start)
    };
    const chunk = {
      type: "chunkString",
      contentType: "string",
      start: Object.assign({}, string3.start),
      end: Object.assign({}, string3.end)
    };
    const replacement = [
      // Take the `labelImageMarker` (now `data`, the `!`)
      events[index2 + 1],
      events[index2 + 2],
      ["enter", call, context],
      // The `[`
      events[index2 + 3],
      events[index2 + 4],
      // The `^`.
      ["enter", marker, context],
      ["exit", marker, context],
      // Everything in between.
      ["enter", string3, context],
      ["enter", chunk, context],
      ["exit", chunk, context],
      ["exit", string3, context],
      // The ending (`]`, properly parsed and labelled).
      events[events.length - 2],
      events[events.length - 1],
      ["exit", call, context]
    ];
    events.splice(index2, events.length - index2 + 1, ...replacement);
    return events;
  }
  function tokenizeGfmFootnoteCall(effects, ok3, nok) {
    const self2 = this;
    const defined = self2.parser.gfmFootnotes || (self2.parser.gfmFootnotes = []);
    let size = 0;
    let data2;
    return start2;
    function start2(code3) {
      effects.enter("gfmFootnoteCall");
      effects.enter("gfmFootnoteCallLabelMarker");
      effects.consume(code3);
      effects.exit("gfmFootnoteCallLabelMarker");
      return callStart;
    }
    function callStart(code3) {
      if (code3 !== 94) return nok(code3);
      effects.enter("gfmFootnoteCallMarker");
      effects.consume(code3);
      effects.exit("gfmFootnoteCallMarker");
      effects.enter("gfmFootnoteCallString");
      effects.enter("chunkString").contentType = "string";
      return callData;
    }
    function callData(code3) {
      if (
        // Too long.
        size > 999 || // Closing brace with nothing.
        code3 === 93 && !data2 || // Space or tab is not supported by GFM for some reason.
        // `\n` and `[` not being supported makes sense.
        code3 === null || code3 === 91 || markdownLineEndingOrSpace(code3)
      ) {
        return nok(code3);
      }
      if (code3 === 93) {
        effects.exit("chunkString");
        const token = effects.exit("gfmFootnoteCallString");
        if (!defined.includes(normalizeIdentifier(self2.sliceSerialize(token)))) {
          return nok(code3);
        }
        effects.enter("gfmFootnoteCallLabelMarker");
        effects.consume(code3);
        effects.exit("gfmFootnoteCallLabelMarker");
        effects.exit("gfmFootnoteCall");
        return ok3;
      }
      if (!markdownLineEndingOrSpace(code3)) {
        data2 = true;
      }
      size++;
      effects.consume(code3);
      return code3 === 92 ? callEscape : callData;
    }
    function callEscape(code3) {
      if (code3 === 91 || code3 === 92 || code3 === 93) {
        effects.consume(code3);
        size++;
        return callData;
      }
      return callData(code3);
    }
  }
  function tokenizeDefinitionStart(effects, ok3, nok) {
    const self2 = this;
    const defined = self2.parser.gfmFootnotes || (self2.parser.gfmFootnotes = []);
    let identifier2;
    let size = 0;
    let data2;
    return start2;
    function start2(code3) {
      effects.enter("gfmFootnoteDefinition")._container = true;
      effects.enter("gfmFootnoteDefinitionLabel");
      effects.enter("gfmFootnoteDefinitionLabelMarker");
      effects.consume(code3);
      effects.exit("gfmFootnoteDefinitionLabelMarker");
      return labelAtMarker;
    }
    function labelAtMarker(code3) {
      if (code3 === 94) {
        effects.enter("gfmFootnoteDefinitionMarker");
        effects.consume(code3);
        effects.exit("gfmFootnoteDefinitionMarker");
        effects.enter("gfmFootnoteDefinitionLabelString");
        effects.enter("chunkString").contentType = "string";
        return labelInside;
      }
      return nok(code3);
    }
    function labelInside(code3) {
      if (
        // Too long.
        size > 999 || // Closing brace with nothing.
        code3 === 93 && !data2 || // Space or tab is not supported by GFM for some reason.
        // `\n` and `[` not being supported makes sense.
        code3 === null || code3 === 91 || markdownLineEndingOrSpace(code3)
      ) {
        return nok(code3);
      }
      if (code3 === 93) {
        effects.exit("chunkString");
        const token = effects.exit("gfmFootnoteDefinitionLabelString");
        identifier2 = normalizeIdentifier(self2.sliceSerialize(token));
        effects.enter("gfmFootnoteDefinitionLabelMarker");
        effects.consume(code3);
        effects.exit("gfmFootnoteDefinitionLabelMarker");
        effects.exit("gfmFootnoteDefinitionLabel");
        return labelAfter;
      }
      if (!markdownLineEndingOrSpace(code3)) {
        data2 = true;
      }
      size++;
      effects.consume(code3);
      return code3 === 92 ? labelEscape : labelInside;
    }
    function labelEscape(code3) {
      if (code3 === 91 || code3 === 92 || code3 === 93) {
        effects.consume(code3);
        size++;
        return labelInside;
      }
      return labelInside(code3);
    }
    function labelAfter(code3) {
      if (code3 === 58) {
        effects.enter("definitionMarker");
        effects.consume(code3);
        effects.exit("definitionMarker");
        if (!defined.includes(identifier2)) {
          defined.push(identifier2);
        }
        return factorySpace(effects, whitespaceAfter, "gfmFootnoteDefinitionWhitespace");
      }
      return nok(code3);
    }
    function whitespaceAfter(code3) {
      return ok3(code3);
    }
  }
  function tokenizeDefinitionContinuation(effects, ok3, nok) {
    return effects.check(blankLine, ok3, effects.attempt(indent2, ok3, nok));
  }
  function gfmFootnoteDefinitionEnd(effects) {
    effects.exit("gfmFootnoteDefinition");
  }
  function tokenizeIndent2(effects, ok3, nok) {
    const self2 = this;
    return factorySpace(effects, afterPrefix, "gfmFootnoteDefinitionIndent", 4 + 1);
    function afterPrefix(code3) {
      const tail = self2.events[self2.events.length - 1];
      return tail && tail[1].type === "gfmFootnoteDefinitionIndent" && tail[2].sliceSerialize(tail[1], true).length === 4 ? ok3(code3) : nok(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-strikethrough@2.1.0/node_modules/micromark-extension-gfm-strikethrough/lib/syntax.js
  function gfmStrikethrough(options) {
    const options_ = options || {};
    let single = options_.singleTilde;
    const tokenizer3 = {
      name: "strikethrough",
      tokenize: tokenizeStrikethrough,
      resolveAll: resolveAllStrikethrough
    };
    if (single === null || single === void 0) {
      single = true;
    }
    return {
      text: {
        [126]: tokenizer3
      },
      insideSpan: {
        null: [tokenizer3]
      },
      attentionMarkers: {
        null: [126]
      }
    };
    function resolveAllStrikethrough(events, context) {
      let index2 = -1;
      while (++index2 < events.length) {
        if (events[index2][0] === "enter" && events[index2][1].type === "strikethroughSequenceTemporary" && events[index2][1]._close) {
          let open = index2;
          while (open--) {
            if (events[open][0] === "exit" && events[open][1].type === "strikethroughSequenceTemporary" && events[open][1]._open && // If the sizes are the same:
            events[index2][1].end.offset - events[index2][1].start.offset === events[open][1].end.offset - events[open][1].start.offset) {
              events[index2][1].type = "strikethroughSequence";
              events[open][1].type = "strikethroughSequence";
              const strikethrough = {
                type: "strikethrough",
                start: Object.assign({}, events[open][1].start),
                end: Object.assign({}, events[index2][1].end)
              };
              const text5 = {
                type: "strikethroughText",
                start: Object.assign({}, events[open][1].end),
                end: Object.assign({}, events[index2][1].start)
              };
              const nextEvents = [["enter", strikethrough, context], ["enter", events[open][1], context], ["exit", events[open][1], context], ["enter", text5, context]];
              const insideSpan2 = context.parser.constructs.insideSpan.null;
              if (insideSpan2) {
                splice(nextEvents, nextEvents.length, 0, resolveAll(insideSpan2, events.slice(open + 1, index2), context));
              }
              splice(nextEvents, nextEvents.length, 0, [["exit", text5, context], ["enter", events[index2][1], context], ["exit", events[index2][1], context], ["exit", strikethrough, context]]);
              splice(events, open - 1, index2 - open + 3, nextEvents);
              index2 = open + nextEvents.length - 2;
              break;
            }
          }
        }
      }
      index2 = -1;
      while (++index2 < events.length) {
        if (events[index2][1].type === "strikethroughSequenceTemporary") {
          events[index2][1].type = "data";
        }
      }
      return events;
    }
    function tokenizeStrikethrough(effects, ok3, nok) {
      const previous3 = this.previous;
      const events = this.events;
      let size = 0;
      return start2;
      function start2(code3) {
        if (previous3 === 126 && events[events.length - 1][1].type !== "characterEscape") {
          return nok(code3);
        }
        effects.enter("strikethroughSequenceTemporary");
        return more(code3);
      }
      function more(code3) {
        const before = classifyCharacter(previous3);
        if (code3 === 126) {
          if (size > 1) return nok(code3);
          effects.consume(code3);
          size++;
          return more;
        }
        if (size < 2 && !single) return nok(code3);
        const token = effects.exit("strikethroughSequenceTemporary");
        const after = classifyCharacter(code3);
        token._open = !after || after === 2 && Boolean(before);
        token._close = !before || before === 2 && Boolean(after);
        return ok3(code3);
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-table@2.1.1/node_modules/micromark-extension-gfm-table/lib/edit-map.js
  var EditMap = class {
    /**
     * Create a new edit map.
     */
    constructor() {
      this.map = [];
    }
    /**
     * Create an edit: a remove and/or add at a certain place.
     *
     * @param {number} index
     * @param {number} remove
     * @param {Array<Event>} add
     * @returns {undefined}
     */
    add(index2, remove5, add) {
      addImplementation(this, index2, remove5, add);
    }
    // To do: add this when moving to `micromark`.
    // /**
    //  * Create an edit: but insert `add` before existing additions.
    //  *
    //  * @param {number} index
    //  * @param {number} remove
    //  * @param {Array<Event>} add
    //  * @returns {undefined}
    //  */
    // addBefore(index, remove, add) {
    //   addImplementation(this, index, remove, add, true)
    // }
    /**
     * Done, change the events.
     *
     * @param {Array<Event>} events
     * @returns {undefined}
     */
    consume(events) {
      this.map.sort(function(a, b) {
        return a[0] - b[0];
      });
      if (this.map.length === 0) {
        return;
      }
      let index2 = this.map.length;
      const vecs = [];
      while (index2 > 0) {
        index2 -= 1;
        vecs.push(events.slice(this.map[index2][0] + this.map[index2][1]), this.map[index2][2]);
        events.length = this.map[index2][0];
      }
      vecs.push(events.slice());
      events.length = 0;
      let slice = vecs.pop();
      while (slice) {
        for (const element of slice) {
          events.push(element);
        }
        slice = vecs.pop();
      }
      this.map.length = 0;
    }
  };
  function addImplementation(editMap, at2, remove5, add) {
    let index2 = 0;
    if (remove5 === 0 && add.length === 0) {
      return;
    }
    while (index2 < editMap.map.length) {
      if (editMap.map[index2][0] === at2) {
        editMap.map[index2][1] += remove5;
        editMap.map[index2][2].push(...add);
        return;
      }
      index2 += 1;
    }
    editMap.map.push([at2, remove5, add]);
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-table@2.1.1/node_modules/micromark-extension-gfm-table/lib/infer.js
  function gfmTableAlign(events, index2) {
    let inDelimiterRow = false;
    const align = [];
    while (index2 < events.length) {
      const event = events[index2];
      if (inDelimiterRow) {
        if (event[0] === "enter") {
          if (event[1].type === "tableContent") {
            align.push(events[index2 + 1][1].type === "tableDelimiterMarker" ? "left" : "none");
          }
        } else if (event[1].type === "tableContent") {
          if (events[index2 - 1][1].type === "tableDelimiterMarker") {
            const alignIndex = align.length - 1;
            align[alignIndex] = align[alignIndex] === "left" ? "center" : "right";
          }
        } else if (event[1].type === "tableDelimiterRow") {
          break;
        }
      } else if (event[0] === "enter" && event[1].type === "tableDelimiterRow") {
        inDelimiterRow = true;
      }
      index2 += 1;
    }
    return align;
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-table@2.1.1/node_modules/micromark-extension-gfm-table/lib/syntax.js
  function gfmTable() {
    return {
      flow: {
        null: {
          name: "table",
          tokenize: tokenizeTable,
          resolveAll: resolveTable
        }
      }
    };
  }
  function tokenizeTable(effects, ok3, nok) {
    const self2 = this;
    let size = 0;
    let sizeB = 0;
    let seen;
    return start2;
    function start2(code3) {
      let index2 = self2.events.length - 1;
      while (index2 > -1) {
        const type = self2.events[index2][1].type;
        if (type === "lineEnding" || // Note: markdown-rs uses `whitespace` instead of `linePrefix`
        type === "linePrefix") index2--;
        else break;
      }
      const tail = index2 > -1 ? self2.events[index2][1].type : null;
      const next = tail === "tableHead" || tail === "tableRow" ? bodyRowStart : headRowBefore;
      if (next === bodyRowStart && self2.parser.lazy[self2.now().line]) {
        return nok(code3);
      }
      return next(code3);
    }
    function headRowBefore(code3) {
      effects.enter("tableHead");
      effects.enter("tableRow");
      return headRowStart(code3);
    }
    function headRowStart(code3) {
      if (code3 === 124) {
        return headRowBreak(code3);
      }
      seen = true;
      sizeB += 1;
      return headRowBreak(code3);
    }
    function headRowBreak(code3) {
      if (code3 === null) {
        return nok(code3);
      }
      if (markdownLineEnding(code3)) {
        if (sizeB > 1) {
          sizeB = 0;
          self2.interrupt = true;
          effects.exit("tableRow");
          effects.enter("lineEnding");
          effects.consume(code3);
          effects.exit("lineEnding");
          return headDelimiterStart;
        }
        return nok(code3);
      }
      if (markdownSpace(code3)) {
        return factorySpace(effects, headRowBreak, "whitespace")(code3);
      }
      sizeB += 1;
      if (seen) {
        seen = false;
        size += 1;
      }
      if (code3 === 124) {
        effects.enter("tableCellDivider");
        effects.consume(code3);
        effects.exit("tableCellDivider");
        seen = true;
        return headRowBreak;
      }
      effects.enter("data");
      return headRowData(code3);
    }
    function headRowData(code3) {
      if (code3 === null || code3 === 124 || markdownLineEndingOrSpace(code3)) {
        effects.exit("data");
        return headRowBreak(code3);
      }
      effects.consume(code3);
      return code3 === 92 ? headRowEscape : headRowData;
    }
    function headRowEscape(code3) {
      if (code3 === 92 || code3 === 124) {
        effects.consume(code3);
        return headRowData;
      }
      return headRowData(code3);
    }
    function headDelimiterStart(code3) {
      self2.interrupt = false;
      if (self2.parser.lazy[self2.now().line]) {
        return nok(code3);
      }
      effects.enter("tableDelimiterRow");
      seen = false;
      if (markdownSpace(code3)) {
        return factorySpace(effects, headDelimiterBefore, "linePrefix", self2.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4)(code3);
      }
      return headDelimiterBefore(code3);
    }
    function headDelimiterBefore(code3) {
      if (code3 === 45 || code3 === 58) {
        return headDelimiterValueBefore(code3);
      }
      if (code3 === 124) {
        seen = true;
        effects.enter("tableCellDivider");
        effects.consume(code3);
        effects.exit("tableCellDivider");
        return headDelimiterCellBefore;
      }
      return headDelimiterNok(code3);
    }
    function headDelimiterCellBefore(code3) {
      if (markdownSpace(code3)) {
        return factorySpace(effects, headDelimiterValueBefore, "whitespace")(code3);
      }
      return headDelimiterValueBefore(code3);
    }
    function headDelimiterValueBefore(code3) {
      if (code3 === 58) {
        sizeB += 1;
        seen = true;
        effects.enter("tableDelimiterMarker");
        effects.consume(code3);
        effects.exit("tableDelimiterMarker");
        return headDelimiterLeftAlignmentAfter;
      }
      if (code3 === 45) {
        sizeB += 1;
        return headDelimiterLeftAlignmentAfter(code3);
      }
      if (code3 === null || markdownLineEnding(code3)) {
        return headDelimiterCellAfter(code3);
      }
      return headDelimiterNok(code3);
    }
    function headDelimiterLeftAlignmentAfter(code3) {
      if (code3 === 45) {
        effects.enter("tableDelimiterFiller");
        return headDelimiterFiller(code3);
      }
      return headDelimiterNok(code3);
    }
    function headDelimiterFiller(code3) {
      if (code3 === 45) {
        effects.consume(code3);
        return headDelimiterFiller;
      }
      if (code3 === 58) {
        seen = true;
        effects.exit("tableDelimiterFiller");
        effects.enter("tableDelimiterMarker");
        effects.consume(code3);
        effects.exit("tableDelimiterMarker");
        return headDelimiterRightAlignmentAfter;
      }
      effects.exit("tableDelimiterFiller");
      return headDelimiterRightAlignmentAfter(code3);
    }
    function headDelimiterRightAlignmentAfter(code3) {
      if (markdownSpace(code3)) {
        return factorySpace(effects, headDelimiterCellAfter, "whitespace")(code3);
      }
      return headDelimiterCellAfter(code3);
    }
    function headDelimiterCellAfter(code3) {
      if (code3 === 124) {
        return headDelimiterBefore(code3);
      }
      if (code3 === null || markdownLineEnding(code3)) {
        if (!seen || size !== sizeB) {
          return headDelimiterNok(code3);
        }
        effects.exit("tableDelimiterRow");
        effects.exit("tableHead");
        return ok3(code3);
      }
      return headDelimiterNok(code3);
    }
    function headDelimiterNok(code3) {
      return nok(code3);
    }
    function bodyRowStart(code3) {
      effects.enter("tableRow");
      return bodyRowBreak(code3);
    }
    function bodyRowBreak(code3) {
      if (code3 === 124) {
        effects.enter("tableCellDivider");
        effects.consume(code3);
        effects.exit("tableCellDivider");
        return bodyRowBreak;
      }
      if (code3 === null || markdownLineEnding(code3)) {
        effects.exit("tableRow");
        return ok3(code3);
      }
      if (markdownSpace(code3)) {
        return factorySpace(effects, bodyRowBreak, "whitespace")(code3);
      }
      effects.enter("data");
      return bodyRowData(code3);
    }
    function bodyRowData(code3) {
      if (code3 === null || code3 === 124 || markdownLineEndingOrSpace(code3)) {
        effects.exit("data");
        return bodyRowBreak(code3);
      }
      effects.consume(code3);
      return code3 === 92 ? bodyRowEscape : bodyRowData;
    }
    function bodyRowEscape(code3) {
      if (code3 === 92 || code3 === 124) {
        effects.consume(code3);
        return bodyRowData;
      }
      return bodyRowData(code3);
    }
  }
  function resolveTable(events, context) {
    let index2 = -1;
    let inFirstCellAwaitingPipe = true;
    let rowKind = 0;
    let lastCell = [0, 0, 0, 0];
    let cell = [0, 0, 0, 0];
    let afterHeadAwaitingFirstBodyRow = false;
    let lastTableEnd = 0;
    let currentTable;
    let currentBody;
    let currentCell;
    const map4 = new EditMap();
    while (++index2 < events.length) {
      const event = events[index2];
      const token = event[1];
      if (event[0] === "enter") {
        if (token.type === "tableHead") {
          afterHeadAwaitingFirstBodyRow = false;
          if (lastTableEnd !== 0) {
            flushTableEnd(map4, context, lastTableEnd, currentTable, currentBody);
            currentBody = void 0;
            lastTableEnd = 0;
          }
          currentTable = {
            type: "table",
            start: Object.assign({}, token.start),
            // Note: correct end is set later.
            end: Object.assign({}, token.end)
          };
          map4.add(index2, 0, [["enter", currentTable, context]]);
        } else if (token.type === "tableRow" || token.type === "tableDelimiterRow") {
          inFirstCellAwaitingPipe = true;
          currentCell = void 0;
          lastCell = [0, 0, 0, 0];
          cell = [0, index2 + 1, 0, 0];
          if (afterHeadAwaitingFirstBodyRow) {
            afterHeadAwaitingFirstBodyRow = false;
            currentBody = {
              type: "tableBody",
              start: Object.assign({}, token.start),
              // Note: correct end is set later.
              end: Object.assign({}, token.end)
            };
            map4.add(index2, 0, [["enter", currentBody, context]]);
          }
          rowKind = token.type === "tableDelimiterRow" ? 2 : currentBody ? 3 : 1;
        } else if (rowKind && (token.type === "data" || token.type === "tableDelimiterMarker" || token.type === "tableDelimiterFiller")) {
          inFirstCellAwaitingPipe = false;
          if (cell[2] === 0) {
            if (lastCell[1] !== 0) {
              cell[0] = cell[1];
              currentCell = flushCell(map4, context, lastCell, rowKind, void 0, currentCell);
              lastCell = [0, 0, 0, 0];
            }
            cell[2] = index2;
          }
        } else if (token.type === "tableCellDivider") {
          if (inFirstCellAwaitingPipe) {
            inFirstCellAwaitingPipe = false;
          } else {
            if (lastCell[1] !== 0) {
              cell[0] = cell[1];
              currentCell = flushCell(map4, context, lastCell, rowKind, void 0, currentCell);
            }
            lastCell = cell;
            cell = [lastCell[1], index2, 0, 0];
          }
        }
      } else if (token.type === "tableHead") {
        afterHeadAwaitingFirstBodyRow = true;
        lastTableEnd = index2;
      } else if (token.type === "tableRow" || token.type === "tableDelimiterRow") {
        lastTableEnd = index2;
        if (lastCell[1] !== 0) {
          cell[0] = cell[1];
          currentCell = flushCell(map4, context, lastCell, rowKind, index2, currentCell);
        } else if (cell[1] !== 0) {
          currentCell = flushCell(map4, context, cell, rowKind, index2, currentCell);
        }
        rowKind = 0;
      } else if (rowKind && (token.type === "data" || token.type === "tableDelimiterMarker" || token.type === "tableDelimiterFiller")) {
        cell[3] = index2;
      }
    }
    if (lastTableEnd !== 0) {
      flushTableEnd(map4, context, lastTableEnd, currentTable, currentBody);
    }
    map4.consume(context.events);
    index2 = -1;
    while (++index2 < context.events.length) {
      const event = context.events[index2];
      if (event[0] === "enter" && event[1].type === "table") {
        event[1]._align = gfmTableAlign(context.events, index2);
      }
    }
    return events;
  }
  function flushCell(map4, context, range, rowKind, rowEnd, previousCell) {
    const groupName = rowKind === 1 ? "tableHeader" : rowKind === 2 ? "tableDelimiter" : "tableData";
    const valueName = "tableContent";
    if (range[0] !== 0) {
      previousCell.end = Object.assign({}, getPoint(context.events, range[0]));
      map4.add(range[0], 0, [["exit", previousCell, context]]);
    }
    const now = getPoint(context.events, range[1]);
    previousCell = {
      type: groupName,
      start: Object.assign({}, now),
      // Note: correct end is set later.
      end: Object.assign({}, now)
    };
    map4.add(range[1], 0, [["enter", previousCell, context]]);
    if (range[2] !== 0) {
      const relatedStart = getPoint(context.events, range[2]);
      const relatedEnd = getPoint(context.events, range[3]);
      const valueToken = {
        type: valueName,
        start: Object.assign({}, relatedStart),
        end: Object.assign({}, relatedEnd)
      };
      map4.add(range[2], 0, [["enter", valueToken, context]]);
      if (rowKind !== 2) {
        const start2 = context.events[range[2]];
        const end = context.events[range[3]];
        start2[1].end = Object.assign({}, end[1].end);
        start2[1].type = "chunkText";
        start2[1].contentType = "text";
        if (range[3] > range[2] + 1) {
          const a = range[2] + 1;
          const b = range[3] - range[2] - 1;
          map4.add(a, b, []);
        }
      }
      map4.add(range[3] + 1, 0, [["exit", valueToken, context]]);
    }
    if (rowEnd !== void 0) {
      previousCell.end = Object.assign({}, getPoint(context.events, rowEnd));
      map4.add(rowEnd, 0, [["exit", previousCell, context]]);
      previousCell = void 0;
    }
    return previousCell;
  }
  function flushTableEnd(map4, context, index2, table, tableBody) {
    const exits = [];
    const related = getPoint(context.events, index2);
    if (tableBody) {
      tableBody.end = Object.assign({}, related);
      exits.push(["exit", tableBody, context]);
    }
    table.end = Object.assign({}, related);
    exits.push(["exit", table, context]);
    map4.add(index2 + 1, 0, exits);
  }
  function getPoint(events, index2) {
    const event = events[index2];
    const side = event[0] === "enter" ? "start" : "end";
    return event[1][side];
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm-task-list-item@2.1.0/node_modules/micromark-extension-gfm-task-list-item/lib/syntax.js
  var tasklistCheck = {
    name: "tasklistCheck",
    tokenize: tokenizeTasklistCheck
  };
  function gfmTaskListItem() {
    return {
      text: {
        [91]: tasklistCheck
      }
    };
  }
  function tokenizeTasklistCheck(effects, ok3, nok) {
    const self2 = this;
    return open;
    function open(code3) {
      if (
        // Exit if there’s stuff before.
        self2.previous !== null || // Exit if not in the first content that is the first child of a list
        // item.
        !self2._gfmTasklistFirstContentOfListItem
      ) {
        return nok(code3);
      }
      effects.enter("taskListCheck");
      effects.enter("taskListCheckMarker");
      effects.consume(code3);
      effects.exit("taskListCheckMarker");
      return inside;
    }
    function inside(code3) {
      if (markdownLineEndingOrSpace(code3)) {
        effects.enter("taskListCheckValueUnchecked");
        effects.consume(code3);
        effects.exit("taskListCheckValueUnchecked");
        return close;
      }
      if (code3 === 88 || code3 === 120) {
        effects.enter("taskListCheckValueChecked");
        effects.consume(code3);
        effects.exit("taskListCheckValueChecked");
        return close;
      }
      return nok(code3);
    }
    function close(code3) {
      if (code3 === 93) {
        effects.enter("taskListCheckMarker");
        effects.consume(code3);
        effects.exit("taskListCheckMarker");
        effects.exit("taskListCheck");
        return after;
      }
      return nok(code3);
    }
    function after(code3) {
      if (markdownLineEnding(code3)) {
        return ok3(code3);
      }
      if (markdownSpace(code3)) {
        return effects.check({
          tokenize: spaceThenNonSpace
        }, ok3, nok)(code3);
      }
      return nok(code3);
    }
  }
  function spaceThenNonSpace(effects, ok3, nok) {
    return factorySpace(effects, after, "whitespace");
    function after(code3) {
      return code3 === null ? nok(code3) : ok3(code3);
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark-extension-gfm@3.0.0/node_modules/micromark-extension-gfm/index.js
  function gfm(options) {
    return combineExtensions([
      gfmAutolinkLiteral(),
      gfmFootnote(),
      gfmStrikethrough(options),
      gfmTable(),
      gfmTaskListItem()
    ]);
  }

  // ../../../ui/node_modules/.pnpm/remark-gfm@4.0.1/node_modules/remark-gfm/lib/index.js
  var emptyOptions4 = {};
  function remarkGfm(options) {
    const self2 = (
      /** @type {Processor<Root>} */
      this
    );
    const settings = options || emptyOptions4;
    const data2 = self2.data();
    const micromarkExtensions = data2.micromarkExtensions || (data2.micromarkExtensions = []);
    const fromMarkdownExtensions = data2.fromMarkdownExtensions || (data2.fromMarkdownExtensions = []);
    const toMarkdownExtensions = data2.toMarkdownExtensions || (data2.toMarkdownExtensions = []);
    micromarkExtensions.push(gfm(settings));
    fromMarkdownExtensions.push(gfmFromMarkdown());
    toMarkdownExtensions.push(gfmToMarkdown(settings));
  }

  // ../../../ui/node_modules/.pnpm/estree-util-value-to-estree@3.5.0/node_modules/estree-util-value-to-estree/dist/estree-util-value-to-estree.js
  function identifier(name2) {
    return { type: "Identifier", name: name2 };
  }
  function literal2(value) {
    return { type: "Literal", value };
  }
  function methodCall(object, name2, args) {
    return {
      type: "CallExpression",
      optional: false,
      callee: {
        type: "MemberExpression",
        computed: false,
        optional: false,
        object,
        property: identifier(name2)
      },
      arguments: args
    };
  }
  function processNumber(number) {
    if (number < 0 || Object.is(number, -0)) {
      return {
        type: "UnaryExpression",
        operator: "-",
        prefix: true,
        argument: processNumber(-number)
      };
    }
    if (typeof number === "bigint") {
      return { type: "Literal", bigint: String(number) };
    }
    if (number === Number.POSITIVE_INFINITY || Number.isNaN(number)) {
      return identifier(String(number));
    }
    return literal2(number);
  }
  function processNumberArray(numbers) {
    return { type: "ArrayExpression", elements: Array.from(numbers, processNumber) };
  }
  function isStringReconstructable(value) {
    return value instanceof URL || value instanceof URLSearchParams;
  }
  function isValueReconstructable(value) {
    return value instanceof Boolean || value instanceof Date || value instanceof Number || value instanceof String;
  }
  var wellKnownSymbols = /* @__PURE__ */ new Map();
  for (const name2 of Reflect.ownKeys(Symbol)) {
    const value = Symbol[name2];
    if (typeof value === "symbol") {
      wellKnownSymbols.set(value, name2);
    }
  }
  function isTemporal(value) {
    return typeof Temporal !== "undefined" && (value instanceof Temporal.Duration || value instanceof Temporal.Instant || value instanceof Temporal.PlainDate || value instanceof Temporal.PlainDateTime || value instanceof Temporal.PlainYearMonth || value instanceof Temporal.PlainMonthDay || value instanceof Temporal.PlainTime || value instanceof Temporal.ZonedDateTime);
  }
  function isTypedArray(value) {
    return value instanceof BigInt64Array || value instanceof BigUint64Array || typeof Float16Array !== "undefined" && value instanceof Float16Array || value instanceof Float32Array || value instanceof Float64Array || value instanceof Int8Array || value instanceof Int16Array || value instanceof Int32Array || value instanceof Uint8Array || value instanceof Uint8ClampedArray || value instanceof Uint16Array || value instanceof Uint32Array;
  }
  function compareContexts(a, b) {
    const aReferencedByB = a.referencedBy.has(b.value);
    const bReferencedByA = b.referencedBy.has(a.value);
    if (aReferencedByB) {
      if (bReferencedByA) {
        return a.count - b.count;
      }
      return -1;
    }
    if (bReferencedByA) {
      return 1;
    }
    return a.count - b.count;
  }
  function replaceAssignment(expression, assignment) {
    if (!assignment || assignment.type !== "AssignmentExpression") {
      return expression;
    }
    let node2 = assignment;
    while (node2.right.type === "AssignmentExpression") {
      node2 = node2.right;
    }
    node2.right = expression;
    return assignment;
  }
  function symbolToEstree(symbol) {
    const name2 = wellKnownSymbols.get(symbol);
    if (name2) {
      return {
        type: "MemberExpression",
        computed: false,
        optional: false,
        object: identifier("Symbol"),
        property: identifier(name2)
      };
    }
    if (symbol.description && symbol === Symbol.for(symbol.description)) {
      return methodCall(identifier("Symbol"), "for", [literal2(symbol.description)]);
    }
    throw new TypeError(`Only global symbols are supported, got: ${String(symbol)}`, {
      cause: symbol
    });
  }
  function property(key, value) {
    const isString = typeof key === "string";
    return {
      type: "Property",
      method: false,
      shorthand: false,
      computed: key === "__proto__" || !isString,
      kind: "init",
      key: isString ? literal2(key) : symbolToEstree(key),
      value
    };
  }
  function valueToEstree(value, options = {}) {
    const stack = [];
    const collectedContexts = /* @__PURE__ */ new Map();
    const namedContexts = [];
    const customTrees = /* @__PURE__ */ new Map();
    function analyze(val) {
      if (typeof val !== "object" && typeof val !== "function") {
        return;
      }
      if (val == null) {
        return;
      }
      const context = collectedContexts.get(val);
      if (context) {
        if (options.preserveReferences) {
          context.count += 1;
        }
        for (const ancestor of stack) {
          context.referencedBy.add(ancestor);
        }
        if (stack.includes(val)) {
          if (!options.preserveReferences) {
            throw new Error(`Found circular reference: ${val}`, { cause: val });
          }
          const parent = stack.at(-1);
          const parentContext = collectedContexts.get(parent);
          parentContext.recursive = true;
          context.recursive = true;
        }
        return;
      }
      collectedContexts.set(val, {
        count: 1,
        recursive: false,
        referencedBy: new Set(stack),
        value: val
      });
      const estree = options?.replacer?.(val);
      if (estree) {
        customTrees.set(val, estree);
        return;
      }
      if (typeof val === "function") {
        throw new TypeError(`Unsupported value: ${val}`, { cause: val });
      }
      if (isTypedArray(val)) {
        return;
      }
      if (isStringReconstructable(val)) {
        return;
      }
      if (isValueReconstructable(val)) {
        return;
      }
      if (value instanceof RegExp) {
        return;
      }
      if (isTemporal(value)) {
        return;
      }
      stack.push(val);
      if (val instanceof Map) {
        for (const pair of val) {
          analyze(pair[0]);
          analyze(pair[1]);
        }
      } else if (Array.isArray(val) || val instanceof Set) {
        for (const entry of val) {
          analyze(entry);
        }
      } else {
        const proto = Object.getPrototypeOf(val);
        if (proto != null && proto !== Object.prototype && !options.instanceAsObject) {
          throw new TypeError(`Unsupported value: ${val}`, { cause: val });
        }
        for (const key of Reflect.ownKeys(val)) {
          analyze(val[key]);
        }
      }
      stack.pop();
    }
    function generate(val, isDeclaration) {
      if (val === void 0) {
        return identifier(String(val));
      }
      if (val == null || typeof val === "string" || typeof val === "boolean") {
        return literal2(val);
      }
      if (typeof val === "bigint" || typeof val === "number") {
        return processNumber(val);
      }
      if (typeof val === "symbol") {
        return symbolToEstree(val);
      }
      const context = collectedContexts.get(val);
      if (!isDeclaration && context?.name) {
        return identifier(context.name);
      }
      const tree = customTrees.get(val);
      if (tree) {
        return tree;
      }
      if (isValueReconstructable(val)) {
        return {
          type: "NewExpression",
          callee: identifier(val.constructor.name),
          arguments: [generate(val.valueOf())]
        };
      }
      if (val instanceof RegExp) {
        return {
          type: "Literal",
          regex: { pattern: val.source, flags: val.flags }
        };
      }
      if (typeof Buffer !== "undefined" && Buffer.isBuffer(val)) {
        return methodCall(identifier("Buffer"), "from", [processNumberArray(val)]);
      }
      if (isTypedArray(val)) {
        return {
          type: "NewExpression",
          callee: identifier(val.constructor.name),
          arguments: [processNumberArray(val)]
        };
      }
      if (isStringReconstructable(val)) {
        return {
          type: "NewExpression",
          callee: identifier(val.constructor.name),
          arguments: [literal2(String(val))]
        };
      }
      if (isTemporal(val)) {
        return methodCall({
          type: "MemberExpression",
          computed: false,
          optional: false,
          object: identifier("Temporal"),
          property: identifier(val.constructor.name)
        }, "from", [literal2(String(val))]);
      }
      if (Array.isArray(val)) {
        const elements = Array.from({ length: val.length });
        let trimmable;
        for (let index2 = 0; index2 < val.length; index2 += 1) {
          if (!(index2 in val)) {
            elements[index2] = null;
            trimmable = void 0;
            continue;
          }
          const child = val[index2];
          const childContext = collectedContexts.get(child);
          if (context && childContext && namedContexts.indexOf(childContext) >= namedContexts.indexOf(context)) {
            elements[index2] = null;
            trimmable || (trimmable = index2);
            childContext.assignment = {
              type: "AssignmentExpression",
              operator: "=",
              left: {
                type: "MemberExpression",
                computed: true,
                optional: false,
                object: identifier(context.name),
                property: literal2(index2)
              },
              right: childContext.assignment || identifier(childContext.name)
            };
          } else {
            elements[index2] = generate(child);
            trimmable = void 0;
          }
        }
        if (trimmable != null) {
          elements.splice(trimmable);
        }
        return {
          type: "ArrayExpression",
          elements
        };
      }
      if (val instanceof Set) {
        const elements = [];
        let finalizer;
        for (const child of val) {
          if (finalizer) {
            finalizer = methodCall(finalizer, "add", [generate(child)]);
          } else {
            const childContext = collectedContexts.get(child);
            if (context && childContext && namedContexts.indexOf(childContext) >= namedContexts.indexOf(context)) {
              finalizer = methodCall(identifier(context.name), "add", [generate(child)]);
            } else {
              elements.push(generate(child));
            }
          }
        }
        if (context && finalizer) {
          context.assignment = replaceAssignment(finalizer, context.assignment);
        }
        return {
          type: "NewExpression",
          callee: identifier("Set"),
          arguments: elements.length ? [{ type: "ArrayExpression", elements }] : []
        };
      }
      if (val instanceof Map) {
        const elements = [];
        let finalizer;
        for (const [key, item] of val) {
          if (finalizer) {
            finalizer = methodCall(finalizer, "set", [generate(key), generate(item)]);
          } else {
            const keyContext = collectedContexts.get(key);
            const itemContext = collectedContexts.get(item);
            if (context && (keyContext && namedContexts.indexOf(keyContext) >= namedContexts.indexOf(context) || itemContext && namedContexts.indexOf(itemContext) >= namedContexts.indexOf(context))) {
              finalizer = methodCall(identifier(context.name), "set", [
                generate(key),
                generate(item)
              ]);
            } else {
              elements.push({
                type: "ArrayExpression",
                elements: [generate(key), generate(item)]
              });
            }
          }
        }
        if (context && finalizer) {
          context.assignment = replaceAssignment(finalizer, context.assignment);
        }
        return {
          type: "NewExpression",
          callee: identifier("Map"),
          arguments: elements.length ? [{ type: "ArrayExpression", elements }] : []
        };
      }
      const properties = [];
      if (Object.getPrototypeOf(val) == null) {
        properties.push({
          type: "Property",
          method: false,
          shorthand: false,
          computed: false,
          kind: "init",
          key: identifier("__proto__"),
          value: literal2(null)
        });
      }
      const object = val;
      const propertyDescriptors = [];
      for (const key of Reflect.ownKeys(val)) {
        const child = object[key];
        const { configurable, enumerable, writable } = Object.getOwnPropertyDescriptor(val, key);
        const childContext = collectedContexts.get(child);
        if (!configurable || !enumerable || !writable) {
          const propertyDescriptor = [property("value", generate(child))];
          if (configurable) {
            propertyDescriptor.push(property("configurable", literal2(true)));
          }
          if (enumerable) {
            propertyDescriptor.push(property("enumerable", literal2(true)));
          }
          if (writable) {
            propertyDescriptor.push(property("writable", literal2(true)));
          }
          propertyDescriptors.push([
            key,
            { type: "ObjectExpression", properties: propertyDescriptor }
          ]);
        } else if (context && childContext && namedContexts.indexOf(childContext) >= namedContexts.indexOf(context)) {
          if (key === "__proto__") {
            propertyDescriptors.push([
              key,
              {
                type: "ObjectExpression",
                properties: [
                  property("value", generate(child)),
                  property("configurable", literal2(true)),
                  property("enumerable", literal2(true)),
                  property("writable", literal2(true))
                ]
              }
            ]);
          } else {
            childContext.assignment = {
              type: "AssignmentExpression",
              operator: "=",
              left: {
                type: "MemberExpression",
                computed: true,
                optional: false,
                object: identifier(context.name),
                property: generate(key)
              },
              right: childContext.assignment || generate(child)
            };
          }
        } else {
          properties.push(property(key, generate(child)));
        }
      }
      const objectExpression = {
        type: "ObjectExpression",
        properties
      };
      if (propertyDescriptors.length) {
        let name2;
        let args;
        if (propertyDescriptors.length === 1) {
          const [[key, expression]] = propertyDescriptors;
          name2 = "defineProperty";
          args = [typeof key === "string" ? literal2(key) : symbolToEstree(key), expression];
        } else {
          name2 = "defineProperties";
          args = [
            {
              type: "ObjectExpression",
              properties: propertyDescriptors.map(([key, expression]) => property(key, expression))
            }
          ];
        }
        if (!context) {
          return methodCall(identifier("Object"), name2, [objectExpression, ...args]);
        }
        context.assignment = replaceAssignment(methodCall(identifier("Object"), name2, [identifier(context.name), ...args]), context.assignment);
      }
      return objectExpression;
    }
    analyze(value);
    for (const [val, context] of collectedContexts) {
      if (context.recursive || context.count > 1) {
        context.name = `$${namedContexts.length}`;
        namedContexts.push(context);
      } else {
        collectedContexts.delete(val);
      }
    }
    if (!namedContexts.length) {
      return generate(value);
    }
    const params = namedContexts.sort(compareContexts).map((context) => ({
      type: "AssignmentPattern",
      left: identifier(context.name),
      right: generate(context.value, true)
    }));
    const rootContext = collectedContexts.get(value);
    const finalizers = [];
    for (const context of collectedContexts.values()) {
      if (context !== rootContext && context.assignment) {
        finalizers.push(context.assignment);
      }
    }
    finalizers.push(rootContext ? rootContext.assignment || identifier(rootContext.name) : generate(value));
    return {
      type: "CallExpression",
      optional: false,
      arguments: [],
      callee: {
        type: "ArrowFunctionExpression",
        expression: false,
        params,
        body: {
          type: "SequenceExpression",
          expressions: finalizers
        }
      }
    };
  }

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/utils-Yx0IL9eW.js
  function toMdxExport(name2, value) {
    return toMdxExportRaw(name2, valueToEstree(value));
  }
  function toMdxExportRaw(name2, expression) {
    return {
      type: "mdxjsEsm",
      value: "",
      data: { estree: {
        type: "Program",
        sourceType: "module",
        body: [{
          type: "ExportNamedDeclaration",
          attributes: [],
          specifiers: [],
          declaration: {
            type: "VariableDeclaration",
            kind: "let",
            declarations: [{
              type: "VariableDeclarator",
              id: {
                type: "Identifier",
                name: name2
              },
              init: expression
            }]
          }
        }]
      } }
    };
  }

  // ../../../ui/node_modules/.pnpm/github-slugger@2.0.0/node_modules/github-slugger/index.js
  var own6 = Object.hasOwnProperty;

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/mdx-plugins/stringifier.js
  function defaultStringifier(config = {}) {
    const { filterMdxAttributes, filterElement = (node2) => {
      switch (node2.type) {
        case "mdxJsxFlowElement":
        case "mdxJsxTextElement":
          switch (node2.name) {
            case "File":
            case "TypeTable":
            case "Callout":
            case "Card":
              return true;
          }
          return "children-only";
      }
      return true;
    }, stringify, ...customExtension } = config;
    function modHandler(handler, ctx) {
      return function(node2, parent, state, info) {
        let visibility = filterElement(node2);
        if (visibility === false) return "";
        if (stringify) {
          const v2 = stringify(node2, parent, state, info, ctx);
          if (v2) return v2;
        }
        const extraInfo = node2.data?._stringify;
        if (extraInfo) if (extraInfo === "children-only") visibility = "children-only";
        else if ("text" in extraInfo) return extraInfo.text;
        else node2 = extraInfo.node;
        if (visibility === "children-only") {
          if (!("children" in node2)) return "";
          switch (node2.type) {
            case "mdxJsxTextElement":
            case "paragraph":
              return state.containerPhrasing(node2, info);
            case "mdxJsxFlowElement":
              return state.containerFlow(node2, info);
            default:
              return state.containerFlow({
                type: "root",
                children: node2.children
              }, info);
          }
        }
        switch (node2.type) {
          case "mdxJsxFlowElement":
          case "mdxJsxTextElement": {
            const stringifiedAttributes = [];
            for (const attr of node2.attributes) {
              if (attr.type === "mdxJsxExpressionAttribute") continue;
              if (filterMdxAttributes && !filterMdxAttributes(node2, attr)) continue;
              const str = typeof attr.value === "string" ? attr.value : attr.value?.value;
              if (!str) continue;
              stringifiedAttributes.push({
                type: "mdxJsxAttribute",
                name: attr.name,
                value: str
              });
            }
            return handler({
              ...node2,
              attributes: stringifiedAttributes
            }, parent, state, info);
          }
          default:
            return handler(node2, parent, state, info);
        }
      };
    }
    const customToMarkdown = { handlers: { _custom(node2, _, state, info) {
      const handlers = state.handlers;
      for (const k in handlers) handlers[k] = modHandler(handlers[k], node2.ctx);
      return state.handle(node2.root, void 0, state, info);
    } } };
    return function(root2, ctx) {
      return toMarkdown({
        type: "_custom",
        root: root2,
        ctx
      }, {
        ...this?.data("settings"),
        extensions: [
          mdxToMarkdown(),
          ...this?.data("toMarkdownExtensions") ?? [],
          customToMarkdown,
          customExtension
        ]
      });
    };
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/initialize/content.js
  var content2 = {
    tokenize: initializeContent
  };
  function initializeContent(effects) {
    const contentStart = effects.attempt(this.parser.constructs.contentInitial, afterContentStartConstruct, paragraphInitial);
    let previous3;
    return contentStart;
    function afterContentStartConstruct(code3) {
      if (code3 === null) {
        effects.consume(code3);
        return;
      }
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      return factorySpace(effects, contentStart, "linePrefix");
    }
    function paragraphInitial(code3) {
      effects.enter("paragraph");
      return lineStart(code3);
    }
    function lineStart(code3) {
      const token = effects.enter("chunkText", {
        contentType: "text",
        previous: previous3
      });
      if (previous3) {
        previous3.next = token;
      }
      previous3 = token;
      return data2(code3);
    }
    function data2(code3) {
      if (code3 === null) {
        effects.exit("chunkText");
        effects.exit("paragraph");
        effects.consume(code3);
        return;
      }
      if (markdownLineEnding(code3)) {
        effects.consume(code3);
        effects.exit("chunkText");
        return lineStart;
      }
      effects.consume(code3);
      return data2;
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/initialize/document.js
  var document = {
    tokenize: initializeDocument
  };
  var containerConstruct = {
    tokenize: tokenizeContainer
  };
  function initializeDocument(effects) {
    const self2 = this;
    const stack = [];
    let continued = 0;
    let childFlow;
    let childToken;
    let lineStartOffset;
    return start2;
    function start2(code3) {
      if (continued < stack.length) {
        const item = stack[continued];
        self2.containerState = item[1];
        return effects.attempt(item[0].continuation, documentContinue, checkNewContainers)(code3);
      }
      return checkNewContainers(code3);
    }
    function documentContinue(code3) {
      continued++;
      if (self2.containerState._closeFlow) {
        self2.containerState._closeFlow = void 0;
        if (childFlow) {
          closeFlow();
        }
        const indexBeforeExits = self2.events.length;
        let indexBeforeFlow = indexBeforeExits;
        let point3;
        while (indexBeforeFlow--) {
          if (self2.events[indexBeforeFlow][0] === "exit" && self2.events[indexBeforeFlow][1].type === "chunkFlow") {
            point3 = self2.events[indexBeforeFlow][1].end;
            break;
          }
        }
        exitContainers(continued);
        let index2 = indexBeforeExits;
        while (index2 < self2.events.length) {
          self2.events[index2][1].end = {
            ...point3
          };
          index2++;
        }
        splice(self2.events, indexBeforeFlow + 1, 0, self2.events.slice(indexBeforeExits));
        self2.events.length = index2;
        return checkNewContainers(code3);
      }
      return start2(code3);
    }
    function checkNewContainers(code3) {
      if (continued === stack.length) {
        if (!childFlow) {
          return documentContinued(code3);
        }
        if (childFlow.currentConstruct && childFlow.currentConstruct.concrete) {
          return flowStart(code3);
        }
        self2.interrupt = Boolean(childFlow.currentConstruct && !childFlow._gfmTableDynamicInterruptHack);
      }
      self2.containerState = {};
      return effects.check(containerConstruct, thereIsANewContainer, thereIsNoNewContainer)(code3);
    }
    function thereIsANewContainer(code3) {
      if (childFlow) closeFlow();
      exitContainers(continued);
      return documentContinued(code3);
    }
    function thereIsNoNewContainer(code3) {
      self2.parser.lazy[self2.now().line] = continued !== stack.length;
      lineStartOffset = self2.now().offset;
      return flowStart(code3);
    }
    function documentContinued(code3) {
      self2.containerState = {};
      return effects.attempt(containerConstruct, containerContinue, flowStart)(code3);
    }
    function containerContinue(code3) {
      continued++;
      stack.push([self2.currentConstruct, self2.containerState]);
      return documentContinued(code3);
    }
    function flowStart(code3) {
      if (code3 === null) {
        if (childFlow) closeFlow();
        exitContainers(0);
        effects.consume(code3);
        return;
      }
      childFlow = childFlow || self2.parser.flow(self2.now());
      effects.enter("chunkFlow", {
        _tokenizer: childFlow,
        contentType: "flow",
        previous: childToken
      });
      return flowContinue(code3);
    }
    function flowContinue(code3) {
      if (code3 === null) {
        writeToChild(effects.exit("chunkFlow"), true);
        exitContainers(0);
        effects.consume(code3);
        return;
      }
      if (markdownLineEnding(code3)) {
        effects.consume(code3);
        writeToChild(effects.exit("chunkFlow"));
        continued = 0;
        self2.interrupt = void 0;
        return start2;
      }
      effects.consume(code3);
      return flowContinue;
    }
    function writeToChild(token, endOfFile) {
      const stream = self2.sliceStream(token);
      if (endOfFile) stream.push(null);
      token.previous = childToken;
      if (childToken) childToken.next = token;
      childToken = token;
      childFlow.defineSkip(token.start);
      childFlow.write(stream);
      if (self2.parser.lazy[token.start.line]) {
        let index2 = childFlow.events.length;
        while (index2--) {
          if (
            // The token starts before the line ending…
            childFlow.events[index2][1].start.offset < lineStartOffset && // …and either is not ended yet…
            (!childFlow.events[index2][1].end || // …or ends after it.
            childFlow.events[index2][1].end.offset > lineStartOffset)
          ) {
            return;
          }
        }
        const indexBeforeExits = self2.events.length;
        let indexBeforeFlow = indexBeforeExits;
        let seen;
        let point3;
        while (indexBeforeFlow--) {
          if (self2.events[indexBeforeFlow][0] === "exit" && self2.events[indexBeforeFlow][1].type === "chunkFlow") {
            if (seen) {
              point3 = self2.events[indexBeforeFlow][1].end;
              break;
            }
            seen = true;
          }
        }
        exitContainers(continued);
        index2 = indexBeforeExits;
        while (index2 < self2.events.length) {
          self2.events[index2][1].end = {
            ...point3
          };
          index2++;
        }
        splice(self2.events, indexBeforeFlow + 1, 0, self2.events.slice(indexBeforeExits));
        self2.events.length = index2;
      }
    }
    function exitContainers(size) {
      let index2 = stack.length;
      while (index2-- > size) {
        const entry = stack[index2];
        self2.containerState = entry[1];
        entry[0].exit.call(self2, effects);
      }
      stack.length = size;
    }
    function closeFlow() {
      childFlow.write([null]);
      childToken = void 0;
      childFlow = void 0;
      self2.containerState._closeFlow = void 0;
    }
  }
  function tokenizeContainer(effects, ok3, nok) {
    return factorySpace(effects, effects.attempt(this.parser.constructs.document, ok3, nok), "linePrefix", this.parser.constructs.disable.null.includes("codeIndented") ? void 0 : 4);
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/initialize/flow.js
  var flow = {
    tokenize: initializeFlow
  };
  function initializeFlow(effects) {
    const self2 = this;
    const initial = effects.attempt(
      // Try to parse a blank line.
      blankLine,
      atBlankEnding,
      // Try to parse initial flow (essentially, only code).
      effects.attempt(this.parser.constructs.flowInitial, afterConstruct, factorySpace(effects, effects.attempt(this.parser.constructs.flow, afterConstruct, effects.attempt(content, afterConstruct)), "linePrefix"))
    );
    return initial;
    function atBlankEnding(code3) {
      if (code3 === null) {
        effects.consume(code3);
        return;
      }
      effects.enter("lineEndingBlank");
      effects.consume(code3);
      effects.exit("lineEndingBlank");
      self2.currentConstruct = void 0;
      return initial;
    }
    function afterConstruct(code3) {
      if (code3 === null) {
        effects.consume(code3);
        return;
      }
      effects.enter("lineEnding");
      effects.consume(code3);
      effects.exit("lineEnding");
      self2.currentConstruct = void 0;
      return initial;
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/initialize/text.js
  var resolver = {
    resolveAll: createResolver()
  };
  var string = initializeFactory("string");
  var text3 = initializeFactory("text");
  function initializeFactory(field) {
    return {
      resolveAll: createResolver(field === "text" ? resolveAllLineSuffixes : void 0),
      tokenize: initializeText
    };
    function initializeText(effects) {
      const self2 = this;
      const constructs2 = this.parser.constructs[field];
      const text5 = effects.attempt(constructs2, start2, notText);
      return start2;
      function start2(code3) {
        return atBreak(code3) ? text5(code3) : notText(code3);
      }
      function notText(code3) {
        if (code3 === null) {
          effects.consume(code3);
          return;
        }
        effects.enter("data");
        effects.consume(code3);
        return data2;
      }
      function data2(code3) {
        if (atBreak(code3)) {
          effects.exit("data");
          return text5(code3);
        }
        effects.consume(code3);
        return data2;
      }
      function atBreak(code3) {
        if (code3 === null) {
          return true;
        }
        const list4 = constructs2[code3];
        let index2 = -1;
        if (list4) {
          while (++index2 < list4.length) {
            const item = list4[index2];
            if (!item.previous || item.previous.call(self2, self2.previous)) {
              return true;
            }
          }
        }
        return false;
      }
    }
  }
  function createResolver(extraResolver) {
    return resolveAllText;
    function resolveAllText(events, context) {
      let index2 = -1;
      let enter;
      while (++index2 <= events.length) {
        if (enter === void 0) {
          if (events[index2] && events[index2][1].type === "data") {
            enter = index2;
            index2++;
          }
        } else if (!events[index2] || events[index2][1].type !== "data") {
          if (index2 !== enter + 2) {
            events[enter][1].end = events[index2 - 1][1].end;
            events.splice(enter + 2, index2 - enter - 2);
            index2 = enter + 2;
          }
          enter = void 0;
        }
      }
      return extraResolver ? extraResolver(events, context) : events;
    }
  }
  function resolveAllLineSuffixes(events, context) {
    let eventIndex = 0;
    while (++eventIndex <= events.length) {
      if ((eventIndex === events.length || events[eventIndex][1].type === "lineEnding") && events[eventIndex - 1][1].type === "data") {
        const data2 = events[eventIndex - 1][1];
        const chunks = context.sliceStream(data2);
        let index2 = chunks.length;
        let bufferIndex = -1;
        let size = 0;
        let tabs;
        while (index2--) {
          const chunk = chunks[index2];
          if (typeof chunk === "string") {
            bufferIndex = chunk.length;
            while (chunk.charCodeAt(bufferIndex - 1) === 32) {
              size++;
              bufferIndex--;
            }
            if (bufferIndex) break;
            bufferIndex = -1;
          } else if (chunk === -2) {
            tabs = true;
            size++;
          } else if (chunk === -1) {
          } else {
            index2++;
            break;
          }
        }
        if (context._contentTypeTextTrailing && eventIndex === events.length) {
          size = 0;
        }
        if (size) {
          const token = {
            type: eventIndex === events.length || tabs || size < 2 ? "lineSuffix" : "hardBreakTrailing",
            start: {
              _bufferIndex: index2 ? bufferIndex : data2.start._bufferIndex + bufferIndex,
              _index: data2.start._index + index2,
              line: data2.end.line,
              column: data2.end.column - size,
              offset: data2.end.offset - size
            },
            end: {
              ...data2.end
            }
          };
          data2.end = {
            ...token.start
          };
          if (data2.start.offset === data2.end.offset) {
            Object.assign(data2, token);
          } else {
            events.splice(eventIndex, 0, ["enter", token, context], ["exit", token, context]);
            eventIndex += 2;
          }
        }
        eventIndex++;
      }
    }
    return events;
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/constructs.js
  var constructs_exports = {};
  __export(constructs_exports, {
    attentionMarkers: () => attentionMarkers,
    contentInitial: () => contentInitial,
    disable: () => disable,
    document: () => document2,
    flow: () => flow2,
    flowInitial: () => flowInitial,
    insideSpan: () => insideSpan,
    string: () => string2,
    text: () => text4
  });
  var document2 = {
    [42]: list,
    [43]: list,
    [45]: list,
    [48]: list,
    [49]: list,
    [50]: list,
    [51]: list,
    [52]: list,
    [53]: list,
    [54]: list,
    [55]: list,
    [56]: list,
    [57]: list,
    [62]: blockQuote
  };
  var contentInitial = {
    [91]: definition
  };
  var flowInitial = {
    [-2]: codeIndented,
    [-1]: codeIndented,
    [32]: codeIndented
  };
  var flow2 = {
    [35]: headingAtx,
    [42]: thematicBreak,
    [45]: [setextUnderline, thematicBreak],
    [60]: htmlFlow,
    [61]: setextUnderline,
    [95]: thematicBreak,
    [96]: codeFenced,
    [126]: codeFenced
  };
  var string2 = {
    [38]: characterReference,
    [92]: characterEscape
  };
  var text4 = {
    [-5]: lineEnding,
    [-4]: lineEnding,
    [-3]: lineEnding,
    [33]: labelStartImage,
    [38]: characterReference,
    [42]: attention,
    [60]: [autolink, htmlText],
    [91]: labelStartLink,
    [92]: [hardBreakEscape, characterEscape],
    [93]: labelEnd,
    [95]: attention,
    [96]: codeText
  };
  var insideSpan = {
    null: [attention, resolver]
  };
  var attentionMarkers = {
    null: [42, 95]
  };
  var disable = {
    null: []
  };

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/create-tokenizer.js
  function createTokenizer(parser, initialize, from) {
    let point3 = {
      _bufferIndex: -1,
      _index: 0,
      line: from && from.line || 1,
      column: from && from.column || 1,
      offset: from && from.offset || 0
    };
    const columnStart = {};
    const resolveAllConstructs = [];
    let chunks = [];
    let stack = [];
    let consumed = true;
    const effects = {
      attempt: constructFactory(onsuccessfulconstruct),
      check: constructFactory(onsuccessfulcheck),
      consume,
      enter,
      exit: exit3,
      interrupt: constructFactory(onsuccessfulcheck, {
        interrupt: true
      })
    };
    const context = {
      code: null,
      containerState: {},
      defineSkip,
      events: [],
      now,
      parser,
      previous: null,
      sliceSerialize,
      sliceStream,
      write
    };
    let state = initialize.tokenize.call(context, effects);
    let expectedCode;
    if (initialize.resolveAll) {
      resolveAllConstructs.push(initialize);
    }
    return context;
    function write(slice) {
      chunks = push(chunks, slice);
      main();
      if (chunks[chunks.length - 1] !== null) {
        return [];
      }
      addResult(initialize, 0);
      context.events = resolveAll(resolveAllConstructs, context.events, context);
      return context.events;
    }
    function sliceSerialize(token, expandTabs) {
      return serializeChunks2(sliceStream(token), expandTabs);
    }
    function sliceStream(token) {
      return sliceChunks(chunks, token);
    }
    function now() {
      const {
        _bufferIndex,
        _index,
        line,
        column,
        offset: offset2
      } = point3;
      return {
        _bufferIndex,
        _index,
        line,
        column,
        offset: offset2
      };
    }
    function defineSkip(value) {
      columnStart[value.line] = value.column;
      accountForPotentialSkip();
    }
    function main() {
      let chunkIndex;
      while (point3._index < chunks.length) {
        const chunk = chunks[point3._index];
        if (typeof chunk === "string") {
          chunkIndex = point3._index;
          if (point3._bufferIndex < 0) {
            point3._bufferIndex = 0;
          }
          while (point3._index === chunkIndex && point3._bufferIndex < chunk.length) {
            go(chunk.charCodeAt(point3._bufferIndex));
          }
        } else {
          go(chunk);
        }
      }
    }
    function go(code3) {
      consumed = void 0;
      expectedCode = code3;
      state = state(code3);
    }
    function consume(code3) {
      if (markdownLineEnding(code3)) {
        point3.line++;
        point3.column = 1;
        point3.offset += code3 === -3 ? 2 : 1;
        accountForPotentialSkip();
      } else if (code3 !== -1) {
        point3.column++;
        point3.offset++;
      }
      if (point3._bufferIndex < 0) {
        point3._index++;
      } else {
        point3._bufferIndex++;
        if (point3._bufferIndex === // Points w/ non-negative `_bufferIndex` reference
        // strings.
        /** @type {string} */
        chunks[point3._index].length) {
          point3._bufferIndex = -1;
          point3._index++;
        }
      }
      context.previous = code3;
      consumed = true;
    }
    function enter(type, fields) {
      const token = fields || {};
      token.type = type;
      token.start = now();
      context.events.push(["enter", token, context]);
      stack.push(token);
      return token;
    }
    function exit3(type) {
      const token = stack.pop();
      token.end = now();
      context.events.push(["exit", token, context]);
      return token;
    }
    function onsuccessfulconstruct(construct, info) {
      addResult(construct, info.from);
    }
    function onsuccessfulcheck(_, info) {
      info.restore();
    }
    function constructFactory(onreturn, fields) {
      return hook;
      function hook(constructs2, returnState, bogusState) {
        let listOfConstructs;
        let constructIndex;
        let currentConstruct;
        let info;
        return Array.isArray(constructs2) ? (
          /* c8 ignore next 1 */
          handleListOfConstructs(constructs2)
        ) : "tokenize" in constructs2 ? (
          // Looks like a construct.
          handleListOfConstructs([
            /** @type {Construct} */
            constructs2
          ])
        ) : handleMapOfConstructs(constructs2);
        function handleMapOfConstructs(map4) {
          return start2;
          function start2(code3) {
            const left = code3 !== null && map4[code3];
            const all2 = code3 !== null && map4.null;
            const list4 = [
              // To do: add more extension tests.
              /* c8 ignore next 2 */
              ...Array.isArray(left) ? left : left ? [left] : [],
              ...Array.isArray(all2) ? all2 : all2 ? [all2] : []
            ];
            return handleListOfConstructs(list4)(code3);
          }
        }
        function handleListOfConstructs(list4) {
          listOfConstructs = list4;
          constructIndex = 0;
          if (list4.length === 0) {
            return bogusState;
          }
          return handleConstruct(list4[constructIndex]);
        }
        function handleConstruct(construct) {
          return start2;
          function start2(code3) {
            info = store2();
            currentConstruct = construct;
            if (!construct.partial) {
              context.currentConstruct = construct;
            }
            if (construct.name && context.parser.constructs.disable.null.includes(construct.name)) {
              return nok(code3);
            }
            return construct.tokenize.call(
              // If we do have fields, create an object w/ `context` as its
              // prototype.
              // This allows a “live binding”, which is needed for `interrupt`.
              fields ? Object.assign(Object.create(context), fields) : context,
              effects,
              ok3,
              nok
            )(code3);
          }
        }
        function ok3(code3) {
          consumed = true;
          onreturn(currentConstruct, info);
          return returnState;
        }
        function nok(code3) {
          consumed = true;
          info.restore();
          if (++constructIndex < listOfConstructs.length) {
            return handleConstruct(listOfConstructs[constructIndex]);
          }
          return bogusState;
        }
      }
    }
    function addResult(construct, from2) {
      if (construct.resolveAll && !resolveAllConstructs.includes(construct)) {
        resolveAllConstructs.push(construct);
      }
      if (construct.resolve) {
        splice(context.events, from2, context.events.length - from2, construct.resolve(context.events.slice(from2), context));
      }
      if (construct.resolveTo) {
        context.events = construct.resolveTo(context.events, context);
      }
    }
    function store2() {
      const startPoint = now();
      const startPrevious = context.previous;
      const startCurrentConstruct = context.currentConstruct;
      const startEventsIndex = context.events.length;
      const startStack = Array.from(stack);
      return {
        from: startEventsIndex,
        restore
      };
      function restore() {
        point3 = startPoint;
        context.previous = startPrevious;
        context.currentConstruct = startCurrentConstruct;
        context.events.length = startEventsIndex;
        stack = startStack;
        accountForPotentialSkip();
      }
    }
    function accountForPotentialSkip() {
      if (point3.line in columnStart && point3.column < 2) {
        point3.column = columnStart[point3.line];
        point3.offset += columnStart[point3.line] - 1;
      }
    }
  }
  function sliceChunks(chunks, token) {
    const startIndex = token.start._index;
    const startBufferIndex = token.start._bufferIndex;
    const endIndex = token.end._index;
    const endBufferIndex = token.end._bufferIndex;
    let view;
    if (startIndex === endIndex) {
      view = [chunks[startIndex].slice(startBufferIndex, endBufferIndex)];
    } else {
      view = chunks.slice(startIndex, endIndex);
      if (startBufferIndex > -1) {
        const head = view[0];
        if (typeof head === "string") {
          view[0] = head.slice(startBufferIndex);
        } else {
          view.shift();
        }
      }
      if (endBufferIndex > 0) {
        view.push(chunks[endIndex].slice(0, endBufferIndex));
      }
    }
    return view;
  }
  function serializeChunks2(chunks, expandTabs) {
    let index2 = -1;
    const result = [];
    let atTab;
    while (++index2 < chunks.length) {
      const chunk = chunks[index2];
      let value;
      if (typeof chunk === "string") {
        value = chunk;
      } else switch (chunk) {
        case -5: {
          value = "\r";
          break;
        }
        case -4: {
          value = "\n";
          break;
        }
        case -3: {
          value = "\r\n";
          break;
        }
        case -2: {
          value = expandTabs ? " " : "	";
          break;
        }
        case -1: {
          if (!expandTabs && atTab) continue;
          value = " ";
          break;
        }
        default: {
          value = String.fromCharCode(chunk);
        }
      }
      atTab = chunk === -2;
      result.push(value);
    }
    return result.join("");
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/parse.js
  function parse4(options) {
    const settings = options || {};
    const constructs2 = (
      /** @type {FullNormalizedExtension} */
      combineExtensions([constructs_exports, ...settings.extensions || []])
    );
    const parser = {
      constructs: constructs2,
      content: create6(content2),
      defined: [],
      document: create6(document),
      flow: create6(flow),
      lazy: {},
      string: create6(string),
      text: create6(text3)
    };
    return parser;
    function create6(initial) {
      return creator;
      function creator(from) {
        return createTokenizer(parser, initial, from);
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/postprocess.js
  function postprocess(events) {
    while (!subtokenize(events)) {
    }
    return events;
  }

  // ../../../ui/node_modules/.pnpm/micromark@4.0.2/node_modules/micromark/lib/preprocess.js
  var search = /[\0\t\n\r]/g;
  function preprocess() {
    let column = 1;
    let buffer = "";
    let start2 = true;
    let atCarriageReturn;
    return preprocessor;
    function preprocessor(value, encoding, end) {
      const chunks = [];
      let match;
      let next;
      let startPosition;
      let endPosition;
      let code3;
      value = buffer + (typeof value === "string" ? value.toString() : new TextDecoder(encoding || void 0).decode(value));
      startPosition = 0;
      buffer = "";
      if (start2) {
        if (value.charCodeAt(0) === 65279) {
          startPosition++;
        }
        start2 = void 0;
      }
      while (startPosition < value.length) {
        search.lastIndex = startPosition;
        match = search.exec(value);
        endPosition = match && match.index !== void 0 ? match.index : value.length;
        code3 = value.charCodeAt(endPosition);
        if (!match) {
          buffer = value.slice(startPosition);
          break;
        }
        if (code3 === 10 && startPosition === endPosition && atCarriageReturn) {
          chunks.push(-3);
          atCarriageReturn = void 0;
        } else {
          if (atCarriageReturn) {
            chunks.push(-5);
            atCarriageReturn = void 0;
          }
          if (startPosition < endPosition) {
            chunks.push(value.slice(startPosition, endPosition));
            column += endPosition - startPosition;
          }
          switch (code3) {
            case 0: {
              chunks.push(65533);
              column++;
              break;
            }
            case 9: {
              next = Math.ceil(column / 4) * 4;
              chunks.push(-2);
              while (column++ < next) chunks.push(-1);
              break;
            }
            case 10: {
              chunks.push(-4);
              column = 1;
              break;
            }
            default: {
              atCarriageReturn = true;
              column = 1;
            }
          }
        }
        startPosition = endPosition + 1;
      }
      if (end) {
        if (atCarriageReturn) chunks.push(-5);
        if (buffer) chunks.push(buffer);
        chunks.push(null);
      }
      return chunks;
    }
  }

  // ../../../ui/node_modules/.pnpm/mdast-util-from-markdown@2.0.2/node_modules/mdast-util-from-markdown/lib/index.js
  var own7 = {}.hasOwnProperty;
  function fromMarkdown(value, encoding, options) {
    if (typeof encoding !== "string") {
      options = encoding;
      encoding = void 0;
    }
    return compiler(options)(postprocess(parse4(options).document().write(preprocess()(value, encoding, true))));
  }
  function compiler(options) {
    const config = {
      transforms: [],
      canContainEols: ["emphasis", "fragment", "heading", "paragraph", "strong"],
      enter: {
        autolink: opener(link2),
        autolinkProtocol: onenterdata,
        autolinkEmail: onenterdata,
        atxHeading: opener(heading2),
        blockQuote: opener(blockQuote2),
        characterEscape: onenterdata,
        characterReference: onenterdata,
        codeFenced: opener(codeFlow),
        codeFencedFenceInfo: buffer,
        codeFencedFenceMeta: buffer,
        codeIndented: opener(codeFlow, buffer),
        codeText: opener(codeText2, buffer),
        codeTextData: onenterdata,
        data: onenterdata,
        codeFlowValue: onenterdata,
        definition: opener(definition3),
        definitionDestinationString: buffer,
        definitionLabelString: buffer,
        definitionTitleString: buffer,
        emphasis: opener(emphasis2),
        hardBreakEscape: opener(hardBreak2),
        hardBreakTrailing: opener(hardBreak2),
        htmlFlow: opener(html2, buffer),
        htmlFlowData: onenterdata,
        htmlText: opener(html2, buffer),
        htmlTextData: onenterdata,
        image: opener(image2),
        label: buffer,
        link: opener(link2),
        listItem: opener(listItem2),
        listItemValue: onenterlistitemvalue,
        listOrdered: opener(list4, onenterlistordered),
        listUnordered: opener(list4),
        paragraph: opener(paragraph2),
        reference: onenterreference,
        referenceString: buffer,
        resourceDestinationString: buffer,
        resourceTitleString: buffer,
        setextHeading: opener(heading2),
        strong: opener(strong2),
        thematicBreak: opener(thematicBreak3)
      },
      exit: {
        atxHeading: closer(),
        atxHeadingSequence: onexitatxheadingsequence,
        autolink: closer(),
        autolinkEmail: onexitautolinkemail,
        autolinkProtocol: onexitautolinkprotocol,
        blockQuote: closer(),
        characterEscapeValue: onexitdata,
        characterReferenceMarkerHexadecimal: onexitcharacterreferencemarker,
        characterReferenceMarkerNumeric: onexitcharacterreferencemarker,
        characterReferenceValue: onexitcharacterreferencevalue,
        characterReference: onexitcharacterreference,
        codeFenced: closer(onexitcodefenced),
        codeFencedFence: onexitcodefencedfence,
        codeFencedFenceInfo: onexitcodefencedfenceinfo,
        codeFencedFenceMeta: onexitcodefencedfencemeta,
        codeFlowValue: onexitdata,
        codeIndented: closer(onexitcodeindented),
        codeText: closer(onexitcodetext),
        codeTextData: onexitdata,
        data: onexitdata,
        definition: closer(),
        definitionDestinationString: onexitdefinitiondestinationstring,
        definitionLabelString: onexitdefinitionlabelstring,
        definitionTitleString: onexitdefinitiontitlestring,
        emphasis: closer(),
        hardBreakEscape: closer(onexithardbreak),
        hardBreakTrailing: closer(onexithardbreak),
        htmlFlow: closer(onexithtmlflow),
        htmlFlowData: onexitdata,
        htmlText: closer(onexithtmltext),
        htmlTextData: onexitdata,
        image: closer(onexitimage),
        label: onexitlabel,
        labelText: onexitlabeltext,
        lineEnding: onexitlineending,
        link: closer(onexitlink),
        listItem: closer(),
        listOrdered: closer(),
        listUnordered: closer(),
        paragraph: closer(),
        referenceString: onexitreferencestring,
        resourceDestinationString: onexitresourcedestinationstring,
        resourceTitleString: onexitresourcetitlestring,
        resource: onexitresource,
        setextHeading: closer(onexitsetextheading),
        setextHeadingLineSequence: onexitsetextheadinglinesequence,
        setextHeadingText: onexitsetextheadingtext,
        strong: closer(),
        thematicBreak: closer()
      }
    };
    configure2(config, (options || {}).mdastExtensions || []);
    const data2 = {};
    return compile;
    function compile(events) {
      let tree = {
        type: "root",
        children: []
      };
      const context = {
        stack: [tree],
        tokenStack: [],
        config,
        enter,
        exit: exit3,
        buffer,
        resume,
        data: data2
      };
      const listStack = [];
      let index2 = -1;
      while (++index2 < events.length) {
        if (events[index2][1].type === "listOrdered" || events[index2][1].type === "listUnordered") {
          if (events[index2][0] === "enter") {
            listStack.push(index2);
          } else {
            const tail = listStack.pop();
            index2 = prepareList(events, tail, index2);
          }
        }
      }
      index2 = -1;
      while (++index2 < events.length) {
        const handler = config[events[index2][0]];
        if (own7.call(handler, events[index2][1].type)) {
          handler[events[index2][1].type].call(Object.assign({
            sliceSerialize: events[index2][2].sliceSerialize
          }, context), events[index2][1]);
        }
      }
      if (context.tokenStack.length > 0) {
        const tail = context.tokenStack[context.tokenStack.length - 1];
        const handler = tail[1] || defaultOnError;
        handler.call(context, void 0, tail[0]);
      }
      tree.position = {
        start: point2(events.length > 0 ? events[0][1].start : {
          line: 1,
          column: 1,
          offset: 0
        }),
        end: point2(events.length > 0 ? events[events.length - 2][1].end : {
          line: 1,
          column: 1,
          offset: 0
        })
      };
      index2 = -1;
      while (++index2 < config.transforms.length) {
        tree = config.transforms[index2](tree) || tree;
      }
      return tree;
    }
    function prepareList(events, start2, length) {
      let index2 = start2 - 1;
      let containerBalance = -1;
      let listSpread = false;
      let listItem3;
      let lineIndex;
      let firstBlankLineIndex;
      let atMarker;
      while (++index2 <= length) {
        const event = events[index2];
        switch (event[1].type) {
          case "listUnordered":
          case "listOrdered":
          case "blockQuote": {
            if (event[0] === "enter") {
              containerBalance++;
            } else {
              containerBalance--;
            }
            atMarker = void 0;
            break;
          }
          case "lineEndingBlank": {
            if (event[0] === "enter") {
              if (listItem3 && !atMarker && !containerBalance && !firstBlankLineIndex) {
                firstBlankLineIndex = index2;
              }
              atMarker = void 0;
            }
            break;
          }
          case "linePrefix":
          case "listItemValue":
          case "listItemMarker":
          case "listItemPrefix":
          case "listItemPrefixWhitespace": {
            break;
          }
          default: {
            atMarker = void 0;
          }
        }
        if (!containerBalance && event[0] === "enter" && event[1].type === "listItemPrefix" || containerBalance === -1 && event[0] === "exit" && (event[1].type === "listUnordered" || event[1].type === "listOrdered")) {
          if (listItem3) {
            let tailIndex = index2;
            lineIndex = void 0;
            while (tailIndex--) {
              const tailEvent = events[tailIndex];
              if (tailEvent[1].type === "lineEnding" || tailEvent[1].type === "lineEndingBlank") {
                if (tailEvent[0] === "exit") continue;
                if (lineIndex) {
                  events[lineIndex][1].type = "lineEndingBlank";
                  listSpread = true;
                }
                tailEvent[1].type = "lineEnding";
                lineIndex = tailIndex;
              } else if (tailEvent[1].type === "linePrefix" || tailEvent[1].type === "blockQuotePrefix" || tailEvent[1].type === "blockQuotePrefixWhitespace" || tailEvent[1].type === "blockQuoteMarker" || tailEvent[1].type === "listItemIndent") {
              } else {
                break;
              }
            }
            if (firstBlankLineIndex && (!lineIndex || firstBlankLineIndex < lineIndex)) {
              listItem3._spread = true;
            }
            listItem3.end = Object.assign({}, lineIndex ? events[lineIndex][1].start : event[1].end);
            events.splice(lineIndex || index2, 0, ["exit", listItem3, event[2]]);
            index2++;
            length++;
          }
          if (event[1].type === "listItemPrefix") {
            const item = {
              type: "listItem",
              _spread: false,
              start: Object.assign({}, event[1].start),
              // @ts-expect-error: we’ll add `end` in a second.
              end: void 0
            };
            listItem3 = item;
            events.splice(index2, 0, ["enter", item, event[2]]);
            index2++;
            length++;
            firstBlankLineIndex = void 0;
            atMarker = true;
          }
        }
      }
      events[start2][1]._spread = listSpread;
      return length;
    }
    function opener(create6, and) {
      return open;
      function open(token) {
        enter.call(this, create6(token), token);
        if (and) and.call(this, token);
      }
    }
    function buffer() {
      this.stack.push({
        type: "fragment",
        children: []
      });
    }
    function enter(node2, token, errorHandler) {
      const parent = this.stack[this.stack.length - 1];
      const siblings = parent.children;
      siblings.push(node2);
      this.stack.push(node2);
      this.tokenStack.push([token, errorHandler || void 0]);
      node2.position = {
        start: point2(token.start),
        // @ts-expect-error: `end` will be patched later.
        end: void 0
      };
    }
    function closer(and) {
      return close;
      function close(token) {
        if (and) and.call(this, token);
        exit3.call(this, token);
      }
    }
    function exit3(token, onExitError) {
      const node2 = this.stack.pop();
      const open = this.tokenStack.pop();
      if (!open) {
        throw new Error("Cannot close `" + token.type + "` (" + stringifyPosition({
          start: token.start,
          end: token.end
        }) + "): it’s not open");
      } else if (open[0].type !== token.type) {
        if (onExitError) {
          onExitError.call(this, token, open[0]);
        } else {
          const handler = open[1] || defaultOnError;
          handler.call(this, token, open[0]);
        }
      }
      node2.position.end = point2(token.end);
    }
    function resume() {
      return toString2(this.stack.pop());
    }
    function onenterlistordered() {
      this.data.expectingFirstListItemValue = true;
    }
    function onenterlistitemvalue(token) {
      if (this.data.expectingFirstListItemValue) {
        const ancestor = this.stack[this.stack.length - 2];
        ancestor.start = Number.parseInt(this.sliceSerialize(token), 10);
        this.data.expectingFirstListItemValue = void 0;
      }
    }
    function onexitcodefencedfenceinfo() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.lang = data3;
    }
    function onexitcodefencedfencemeta() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.meta = data3;
    }
    function onexitcodefencedfence() {
      if (this.data.flowCodeInside) return;
      this.buffer();
      this.data.flowCodeInside = true;
    }
    function onexitcodefenced() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.value = data3.replace(/^(\r?\n|\r)|(\r?\n|\r)$/g, "");
      this.data.flowCodeInside = void 0;
    }
    function onexitcodeindented() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.value = data3.replace(/(\r?\n|\r)$/g, "");
    }
    function onexitdefinitionlabelstring(token) {
      const label = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.label = label;
      node2.identifier = normalizeIdentifier(this.sliceSerialize(token)).toLowerCase();
    }
    function onexitdefinitiontitlestring() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.title = data3;
    }
    function onexitdefinitiondestinationstring() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.url = data3;
    }
    function onexitatxheadingsequence(token) {
      const node2 = this.stack[this.stack.length - 1];
      if (!node2.depth) {
        const depth = this.sliceSerialize(token).length;
        node2.depth = depth;
      }
    }
    function onexitsetextheadingtext() {
      this.data.setextHeadingSlurpLineEnding = true;
    }
    function onexitsetextheadinglinesequence(token) {
      const node2 = this.stack[this.stack.length - 1];
      node2.depth = this.sliceSerialize(token).codePointAt(0) === 61 ? 1 : 2;
    }
    function onexitsetextheading() {
      this.data.setextHeadingSlurpLineEnding = void 0;
    }
    function onenterdata(token) {
      const node2 = this.stack[this.stack.length - 1];
      const siblings = node2.children;
      let tail = siblings[siblings.length - 1];
      if (!tail || tail.type !== "text") {
        tail = text5();
        tail.position = {
          start: point2(token.start),
          // @ts-expect-error: we’ll add `end` later.
          end: void 0
        };
        siblings.push(tail);
      }
      this.stack.push(tail);
    }
    function onexitdata(token) {
      const tail = this.stack.pop();
      tail.value += this.sliceSerialize(token);
      tail.position.end = point2(token.end);
    }
    function onexitlineending(token) {
      const context = this.stack[this.stack.length - 1];
      if (this.data.atHardBreak) {
        const tail = context.children[context.children.length - 1];
        tail.position.end = point2(token.end);
        this.data.atHardBreak = void 0;
        return;
      }
      if (!this.data.setextHeadingSlurpLineEnding && config.canContainEols.includes(context.type)) {
        onenterdata.call(this, token);
        onexitdata.call(this, token);
      }
    }
    function onexithardbreak() {
      this.data.atHardBreak = true;
    }
    function onexithtmlflow() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.value = data3;
    }
    function onexithtmltext() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.value = data3;
    }
    function onexitcodetext() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.value = data3;
    }
    function onexitlink() {
      const node2 = this.stack[this.stack.length - 1];
      if (this.data.inReference) {
        const referenceType = this.data.referenceType || "shortcut";
        node2.type += "Reference";
        node2.referenceType = referenceType;
        delete node2.url;
        delete node2.title;
      } else {
        delete node2.identifier;
        delete node2.label;
      }
      this.data.referenceType = void 0;
    }
    function onexitimage() {
      const node2 = this.stack[this.stack.length - 1];
      if (this.data.inReference) {
        const referenceType = this.data.referenceType || "shortcut";
        node2.type += "Reference";
        node2.referenceType = referenceType;
        delete node2.url;
        delete node2.title;
      } else {
        delete node2.identifier;
        delete node2.label;
      }
      this.data.referenceType = void 0;
    }
    function onexitlabeltext(token) {
      const string3 = this.sliceSerialize(token);
      const ancestor = this.stack[this.stack.length - 2];
      ancestor.label = decodeString(string3);
      ancestor.identifier = normalizeIdentifier(string3).toLowerCase();
    }
    function onexitlabel() {
      const fragment = this.stack[this.stack.length - 1];
      const value = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      this.data.inReference = true;
      if (node2.type === "link") {
        const children = fragment.children;
        node2.children = children;
      } else {
        node2.alt = value;
      }
    }
    function onexitresourcedestinationstring() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.url = data3;
    }
    function onexitresourcetitlestring() {
      const data3 = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.title = data3;
    }
    function onexitresource() {
      this.data.inReference = void 0;
    }
    function onenterreference() {
      this.data.referenceType = "collapsed";
    }
    function onexitreferencestring(token) {
      const label = this.resume();
      const node2 = this.stack[this.stack.length - 1];
      node2.label = label;
      node2.identifier = normalizeIdentifier(this.sliceSerialize(token)).toLowerCase();
      this.data.referenceType = "full";
    }
    function onexitcharacterreferencemarker(token) {
      this.data.characterReferenceType = token.type;
    }
    function onexitcharacterreferencevalue(token) {
      const data3 = this.sliceSerialize(token);
      const type = this.data.characterReferenceType;
      let value;
      if (type) {
        value = decodeNumericCharacterReference(data3, type === "characterReferenceMarkerNumeric" ? 10 : 16);
        this.data.characterReferenceType = void 0;
      } else {
        const result = decodeNamedCharacterReference(data3);
        value = result;
      }
      const tail = this.stack[this.stack.length - 1];
      tail.value += value;
    }
    function onexitcharacterreference(token) {
      const tail = this.stack.pop();
      tail.position.end = point2(token.end);
    }
    function onexitautolinkprotocol(token) {
      onexitdata.call(this, token);
      const node2 = this.stack[this.stack.length - 1];
      node2.url = this.sliceSerialize(token);
    }
    function onexitautolinkemail(token) {
      onexitdata.call(this, token);
      const node2 = this.stack[this.stack.length - 1];
      node2.url = "mailto:" + this.sliceSerialize(token);
    }
    function blockQuote2() {
      return {
        type: "blockquote",
        children: []
      };
    }
    function codeFlow() {
      return {
        type: "code",
        lang: null,
        meta: null,
        value: ""
      };
    }
    function codeText2() {
      return {
        type: "inlineCode",
        value: ""
      };
    }
    function definition3() {
      return {
        type: "definition",
        identifier: "",
        label: null,
        title: null,
        url: ""
      };
    }
    function emphasis2() {
      return {
        type: "emphasis",
        children: []
      };
    }
    function heading2() {
      return {
        type: "heading",
        // @ts-expect-error `depth` will be set later.
        depth: 0,
        children: []
      };
    }
    function hardBreak2() {
      return {
        type: "break"
      };
    }
    function html2() {
      return {
        type: "html",
        value: ""
      };
    }
    function image2() {
      return {
        type: "image",
        title: null,
        url: "",
        alt: null
      };
    }
    function link2() {
      return {
        type: "link",
        title: null,
        url: "",
        children: []
      };
    }
    function list4(token) {
      return {
        type: "list",
        ordered: token.type === "listOrdered",
        start: null,
        spread: token._spread,
        children: []
      };
    }
    function listItem2(token) {
      return {
        type: "listItem",
        spread: token._spread,
        checked: null,
        children: []
      };
    }
    function paragraph2() {
      return {
        type: "paragraph",
        children: []
      };
    }
    function strong2() {
      return {
        type: "strong",
        children: []
      };
    }
    function text5() {
      return {
        type: "text",
        value: ""
      };
    }
    function thematicBreak3() {
      return {
        type: "thematicBreak"
      };
    }
  }
  function point2(d) {
    return {
      line: d.line,
      column: d.column,
      offset: d.offset
    };
  }
  function configure2(combined, extensions) {
    let index2 = -1;
    while (++index2 < extensions.length) {
      const value = extensions[index2];
      if (Array.isArray(value)) {
        configure2(combined, value);
      } else {
        extension(combined, value);
      }
    }
  }
  function extension(combined, extension2) {
    let key;
    for (key in extension2) {
      if (own7.call(extension2, key)) {
        switch (key) {
          case "canContainEols": {
            const right = extension2[key];
            if (right) {
              combined[key].push(...right);
            }
            break;
          }
          case "transforms": {
            const right = extension2[key];
            if (right) {
              combined[key].push(...right);
            }
            break;
          }
          case "enter":
          case "exit": {
            const right = extension2[key];
            if (right) {
              Object.assign(combined[key], right);
            }
            break;
          }
        }
      }
    }
  }
  function defaultOnError(left, right) {
    if (left) {
      throw new Error("Cannot close `" + left.type + "` (" + stringifyPosition({
        start: left.start,
        end: left.end
      }) + "): a different token (`" + right.type + "`, " + stringifyPosition({
        start: right.start,
        end: right.end
      }) + ") is open");
    } else {
      throw new Error("Cannot close document, a token (`" + right.type + "`, " + stringifyPosition({
        start: right.start,
        end: right.end
      }) + ") is still open");
    }
  }

  // ../../../ui/node_modules/.pnpm/remark-parse@11.0.0/node_modules/remark-parse/lib/index.js
  function remarkParse(options) {
    const self2 = this;
    self2.parser = parser;
    function parser(doc) {
      return fromMarkdown(doc, {
        ...self2.data("settings"),
        ...options,
        // Note: these options are not in the readme.
        // The goal is for them to be set by plugins on `data` instead of being
        // passed by users.
        extensions: self2.data("micromarkExtensions") || [],
        mdastExtensions: self2.data("fromMarkdownExtensions") || []
      });
    }
  }

  // ../../../ui/node_modules/.pnpm/remark-stringify@11.0.0/node_modules/remark-stringify/lib/index.js
  function remarkStringify(options) {
    const self2 = this;
    self2.compiler = compiler2;
    function compiler2(tree) {
      return toMarkdown(tree, {
        ...self2.data("settings"),
        ...options,
        // Note: this option is not in the readme.
        // The goal is for it to be set by plugins on `data` instead of being
        // passed by users.
        extensions: self2.data("toMarkdownExtensions") || []
      });
    }
  }

  // ../../../ui/node_modules/.pnpm/remark@15.0.1/node_modules/remark/index.js
  var remark = unified().use(remarkParse).use(remarkStringify).freeze();

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/mdx-plugins/remark-structure.js
  var remarkStructureDefaultOptions = {
    types: [
      "heading",
      "paragraph",
      "blockquote",
      "tableCell",
      "mdxJsxFlowElement"
    ],
    mdxTypes(node2) {
      return !node2.children || node2.children.length === 0;
    },
    exportAs: false
  };
  function remarkStructure({ types: types2 = remarkStructureDefaultOptions.types, mdxTypes = remarkStructureDefaultOptions.mdxTypes, stringify: stringifyOptions, allowedMdxAttributes, exportAs = remarkStructureDefaultOptions.exportAs } = {}) {
    if (Array.isArray(types2)) {
      const arr = types2;
      types2 = (node2) => arr.includes(node2.type);
    }
    const stringify = typeof stringifyOptions === "function" ? stringifyOptions : defaultStringifier2({
      filterMdxAttributes: Array.isArray(allowedMdxAttributes) ? (_node, attribute) => attribute.type === "mdxJsxAttribute" && allowedMdxAttributes.includes(attribute.name) : allowedMdxAttributes,
      ...stringifyOptions
    });
    return (tree, file) => {
      const data2 = {
        contents: [],
        headings: []
      };
      let lastHeading;
      if (file.data.frontmatter) {
        const openapiData = file.data.frontmatter._openapi?.structuredData;
        if (openapiData) {
          data2.headings.push(...openapiData.headings);
          data2.contents.push(...openapiData.contents);
        }
      }
      const stringifierCtx = { addContent(...content3) {
        for (const item of content3) data2.contents.push({
          ...item,
          heading: item.heading ?? lastHeading
        });
      } };
      visit2(tree, (element) => {
        var _a;
        if (!types2(element)) return;
        switch (element.type) {
          case "root":
            return;
          case "mdxJsxFlowElement":
          case "mdxJsxTextElement":
            if (!mdxTypes(element)) return;
            break;
          case "heading": {
            element.data || (element.data = {});
            (_a = element.data).hProperties || (_a.hProperties = {});
            const id = element.data.hProperties.id;
            if (typeof id !== "string") {
              console.warn("[remark-structure] hProperties.id is missing in heading node, it is required to generate heading data. You can add remark-heading prior to remark-structure to generate heading IDs.");
              return "skip";
            }
            const content4 = stringify.call(this, element, stringifierCtx).trim();
            if (content4.length > 0) data2.headings.push({
              id,
              content: content4
            });
            lastHeading = id;
            return "skip";
          }
        }
        const content3 = stringify.call(this, element, stringifierCtx).trim();
        if (content3.length > 0) data2.contents.push({
          heading: lastHeading,
          content: content3
        });
        return "skip";
      });
      file.data.structuredData = data2;
      if (exportAs) tree.children.unshift(toMdxExport(typeof exportAs === "string" ? exportAs : "structuredData", data2));
    };
  }
  function defaultStringifier2(config) {
    return defaultStringifier({
      ...config,
      handlers: {
        link(node2, _, state, info) {
          return state.containerPhrasing(node2, info);
        },
        heading(node2, _, state, info) {
          return state.containerPhrasing(node2, info);
        },
        image() {
          return "";
        },
        ...config.handlers
      },
      stringify(node2, parent, state, info, ctx) {
        if (node2.data?.structuredData) ctx.addContent(...node2.data.structuredData.contents);
        return config.stringify?.(node2, parent, state, info, ctx);
      }
    });
  }

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/utils-Dn9VIXRN.js
  var VisitBreak = Symbol("VisitBreak");

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/build-doc-CTQdPLc0.js
  function buildDocuments(indexes) {
    const docs = [];
    for (const page of indexes) {
      const pageTag = page.tag ?? [];
      const tags = Array.isArray(pageTag) ? pageTag : [pageTag];
      const data2 = page.structuredData;
      let id = 0;
      docs.push({
        id: page.id,
        page_id: page.id,
        type: "page",
        content: page.title,
        breadcrumbs: page.breadcrumbs,
        tags,
        url: page.url
      });
      const nextId = () => `${page.id}-${id++}`;
      if (page.description) docs.push({
        id: nextId(),
        page_id: page.id,
        tags,
        type: "text",
        url: page.url,
        content: page.description
      });
      for (const heading2 of data2.headings) docs.push({
        id: nextId(),
        page_id: page.id,
        type: "heading",
        tags,
        url: `${page.url}#${heading2.id}`,
        content: heading2.content
      });
      for (const content3 of data2.contents) docs.push({
        id: nextId(),
        page_id: page.id,
        tags,
        type: "text",
        url: content3.heading ? `${page.url}#${content3.heading}` : page.url,
        content: content3.content
      });
    }
    return docs;
  }

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/search/index.js
  function escapeRegExp(input) {
    return input.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  }
  function buildRegexFromQuery(q) {
    const trimmed = q.trim();
    if (trimmed.length === 0) return null;
    const terms = Array.from(new Set(trimmed.split(/\s+/).filter(Boolean)));
    if (terms.length === 0) return null;
    const escaped = terms.map(escapeRegExp).join("|");
    return new RegExp(`(${escaped})`, "gi");
  }
  var processor = remark();
  function createContentHighlighter(query) {
    const regex2 = typeof query === "string" ? buildRegexFromQuery(query) : query;
    return {
      /**
      * @deprecated use `highlightMarkdown()` instead.
      */
      highlight(content3) {
        if (!regex2) return [{
          type: "text",
          content: content3
        }];
        const out = [];
        let i = 0;
        for (const match of content3.matchAll(regex2)) {
          if (i < match.index) out.push({
            type: "text",
            content: content3.substring(i, match.index)
          });
          out.push({
            type: "text",
            content: match[0],
            styles: { highlight: true }
          });
          i = match.index + match[0].length;
        }
        if (i < content3.length) out.push({
          type: "text",
          content: content3.substring(i)
        });
        return out;
      },
      /**
      * @param content - Markdown, it assumes the content is already sanitized & safe, no escape is performed.
      */
      highlightMarkdown(content3) {
        if (!regex2) return content3;
        const tree = processor.parse(content3);
        highlightInTree(tree, regex2);
        return processor.stringify(tree).trim();
      }
    };
  }
  function highlightInTree(tree, regex2) {
    visit2(tree, "text", (node2) => {
      let out = "";
      const content3 = node2.value;
      let i = 0;
      for (const match of content3.matchAll(regex2)) {
        if (i < match.index) out += content3.substring(i, match.index);
        out += `<mark>${match[0]}</mark>`;
        i = match.index + match[0].length;
      }
      if (i < content3.length) out += content3.substring(i);
      node2.type = "html";
      node2.value = out;
    });
  }

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/remove-undefined-CzMSKybq.js
  function removeUndefined(value, deep = false) {
    const obj = value;
    for (const key in obj) {
      if (obj[key] === void 0) delete obj[key];
      if (!deep) continue;
      const entry = obj[key];
      if (isPlainObject2(entry)) {
        removeUndefined(entry, deep);
        continue;
      }
      if (Array.isArray(entry)) {
        for (const item of entry) if (isPlainObject2(item)) removeUndefined(item, deep);
      }
    }
    return value;
  }
  function isPlainObject2(value) {
    if (typeof value !== "object" || value === null) return false;
    const prototype = Object.getPrototypeOf(value);
    return prototype === null || prototype === Object.prototype || Object.getPrototypeOf(prototype) === null;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/tokenizer/languages.js
  var STEMMERS = {
    arabic: "ar",
    armenian: "am",
    bulgarian: "bg",
    czech: "cz",
    danish: "dk",
    dutch: "nl",
    english: "en",
    finnish: "fi",
    french: "fr",
    german: "de",
    greek: "gr",
    hungarian: "hu",
    indian: "in",
    indonesian: "id",
    irish: "ie",
    italian: "it",
    lithuanian: "lt",
    nepali: "np",
    norwegian: "no",
    portuguese: "pt",
    romanian: "ro",
    russian: "ru",
    serbian: "rs",
    slovenian: "ru",
    spanish: "es",
    swedish: "se",
    tamil: "ta",
    turkish: "tr",
    ukrainian: "uk",
    sanskrit: "sk"
  };
  var SPLITTERS = {
    dutch: /[^A-Za-zàèéìòóù0-9_'-]+/gim,
    english: /[^A-Za-zàèéìòóù0-9_'-]+/gim,
    french: /[^a-z0-9äâàéèëêïîöôùüûœç-]+/gim,
    italian: /[^A-Za-zàèéìòóù0-9_'-]+/gim,
    norwegian: /[^a-z0-9_æøåÆØÅäÄöÖüÜ]+/gim,
    portuguese: /[^a-z0-9à-úÀ-Ú]/gim,
    russian: /[^a-z0-9а-яА-ЯёЁ]+/gim,
    spanish: /[^a-z0-9A-Zá-úÁ-ÚñÑüÜ]+/gim,
    swedish: /[^a-z0-9_åÅäÄöÖüÜ-]+/gim,
    german: /[^a-z0-9A-ZäöüÄÖÜß]+/gim,
    finnish: /[^a-z0-9äöÄÖ]+/gim,
    danish: /[^a-z0-9æøåÆØÅ]+/gim,
    hungarian: /[^a-z0-9áéíóöőúüűÁÉÍÓÖŐÚÜŰ]+/gim,
    romanian: /[^a-z0-9ăâîșțĂÂÎȘȚ]+/gim,
    serbian: /[^a-z0-9čćžšđČĆŽŠĐ]+/gim,
    turkish: /[^a-z0-9çÇğĞıİöÖşŞüÜ]+/gim,
    lithuanian: /[^a-z0-9ąčęėįšųūžĄČĘĖĮŠŲŪŽ]+/gim,
    arabic: /[^a-z0-9أ-ي]+/gim,
    nepali: /[^a-z0-9अ-ह]+/gim,
    irish: /[^a-z0-9áéíóúÁÉÍÓÚ]+/gim,
    indian: /[^a-z0-9अ-ह]+/gim,
    armenian: /[^a-z0-9ա-ֆ]+/gim,
    greek: /[^a-z0-9α-ωά-ώ]+/gim,
    indonesian: /[^a-z0-9]+/gim,
    ukrainian: /[^a-z0-9а-яА-ЯіїєІЇЄ]+/gim,
    slovenian: /[^a-z0-9čžšČŽŠ]+/gim,
    bulgarian: /[^a-z0-9а-яА-Я]+/gim,
    tamil: /[^a-z0-9அ-ஹ]+/gim,
    sanskrit: /[^a-z0-9A-Zāīūṛḷṃṁḥśṣṭḍṇṅñḻḹṝ]+/gim,
    czech: /[^A-Z0-9a-zěščřžýáíéúůóťďĚŠČŘŽÝÁÍÉÓÚŮŤĎ-]+/gim
  };
  var SUPPORTED_LANGUAGES = Object.keys(STEMMERS);
  function getLocale(language) {
    return language !== void 0 && SUPPORTED_LANGUAGES.includes(language) ? STEMMERS[language] : void 0;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/utils.js
  var baseId = Date.now().toString().slice(5);
  var lastId = 0;
  var nano = BigInt(1e3);
  var milli = BigInt(1e6);
  var second = BigInt(1e9);
  var MAX_ARGUMENT_FOR_STACK = 65535;
  function safeArrayPush(arr, newArr) {
    if (newArr.length < MAX_ARGUMENT_FOR_STACK) {
      Array.prototype.push.apply(arr, newArr);
    } else {
      const newArrLength = newArr.length;
      for (let i = 0; i < newArrLength; i += MAX_ARGUMENT_FOR_STACK) {
        Array.prototype.push.apply(arr, newArr.slice(i, i + MAX_ARGUMENT_FOR_STACK));
      }
    }
  }
  function sprintf(template, ...args) {
    return template.replace(/%(?:(?<position>\d+)\$)?(?<width>-?\d*\.?\d*)(?<type>[dfs])/g, function(...replaceArgs) {
      const groups = replaceArgs[replaceArgs.length - 1];
      const { width: rawWidth, type, position: position2 } = groups;
      const replacement = position2 ? args[Number.parseInt(position2) - 1] : args.shift();
      const width = rawWidth === "" ? 0 : Number.parseInt(rawWidth);
      switch (type) {
        case "d":
          return replacement.toString().padStart(width, "0");
        case "f": {
          let value = replacement;
          const [padding, precision] = rawWidth.split(".").map((w) => Number.parseFloat(w));
          if (typeof precision === "number" && precision >= 0) {
            value = value.toFixed(precision);
          }
          return typeof padding === "number" && padding >= 0 ? value.toString().padStart(width, "0") : value.toString();
        }
        case "s":
          return width < 0 ? replacement.toString().padEnd(-width, " ") : replacement.toString().padStart(width, " ");
        default:
          return replacement;
      }
    });
  }
  function isInsideWebWorker() {
    return typeof WorkerGlobalScope !== "undefined" && self instanceof WorkerGlobalScope;
  }
  function isInsideNode() {
    return typeof process !== "undefined" && process.release && process.release.name === "node";
  }
  function getNanosecondTimeViaPerformance() {
    return BigInt(Math.floor(performance.now() * 1e6));
  }
  function formatNanoseconds(value) {
    if (typeof value === "number") {
      value = BigInt(value);
    }
    if (value < nano) {
      return `${value}ns`;
    } else if (value < milli) {
      return `${value / nano}μs`;
    } else if (value < second) {
      return `${value / milli}ms`;
    }
    return `${value / second}s`;
  }
  function getNanosecondsTime() {
    if (isInsideWebWorker()) {
      return getNanosecondTimeViaPerformance();
    }
    if (isInsideNode()) {
      return process.hrtime.bigint();
    }
    if (typeof process !== "undefined" && typeof process?.hrtime?.bigint === "function") {
      return process.hrtime.bigint();
    }
    if (typeof performance !== "undefined") {
      return getNanosecondTimeViaPerformance();
    }
    return BigInt(0);
  }
  function uniqueId() {
    return `${baseId}-${lastId++}`;
  }
  function getOwnProperty(object, property2) {
    if (Object.hasOwn === void 0) {
      return Object.prototype.hasOwnProperty.call(object, property2) ? object[property2] : void 0;
    }
    return Object.hasOwn(object, property2) ? object[property2] : void 0;
  }
  function sortTokenScorePredicate(a, b) {
    if (b[1] === a[1]) {
      return a[0] - b[0];
    }
    return b[1] - a[1];
  }
  function intersect(arrays) {
    if (arrays.length === 0) {
      return [];
    } else if (arrays.length === 1) {
      return arrays[0];
    }
    for (let i = 1; i < arrays.length; i++) {
      if (arrays[i].length < arrays[0].length) {
        const tmp = arrays[0];
        arrays[0] = arrays[i];
        arrays[i] = tmp;
      }
    }
    const set = /* @__PURE__ */ new Map();
    for (const elem of arrays[0]) {
      set.set(elem, 1);
    }
    for (let i = 1; i < arrays.length; i++) {
      let found = 0;
      for (const elem of arrays[i]) {
        const count3 = set.get(elem);
        if (count3 === i) {
          set.set(elem, count3 + 1);
          found++;
        }
      }
      if (found === 0)
        return [];
    }
    return arrays[0].filter((e) => {
      const count3 = set.get(e);
      if (count3 !== void 0)
        set.set(e, 0);
      return count3 === arrays.length;
    });
  }
  function getDocumentProperties(doc, paths) {
    const properties = {};
    const pathsLength = paths.length;
    for (let i = 0; i < pathsLength; i++) {
      const path2 = paths[i];
      const pathTokens = path2.split(".");
      let current2 = doc;
      const pathTokensLength = pathTokens.length;
      for (let j = 0; j < pathTokensLength; j++) {
        current2 = current2[pathTokens[j]];
        if (typeof current2 === "object") {
          if (current2 !== null && "lat" in current2 && "lon" in current2 && typeof current2.lat === "number" && typeof current2.lon === "number") {
            current2 = properties[path2] = current2;
            break;
          } else if (!Array.isArray(current2) && current2 !== null && j === pathTokensLength - 1) {
            current2 = void 0;
            break;
          }
        } else if ((current2 === null || typeof current2 !== "object") && j < pathTokensLength - 1) {
          current2 = void 0;
          break;
        }
      }
      if (typeof current2 !== "undefined") {
        properties[path2] = current2;
      }
    }
    return properties;
  }
  function getNested(obj, path2) {
    const props = getDocumentProperties(obj, [path2]);
    return props[path2];
  }
  var mapDistanceToMeters = {
    cm: 0.01,
    m: 1,
    km: 1e3,
    ft: 0.3048,
    yd: 0.9144,
    mi: 1609.344
  };
  function convertDistanceToMeters(distance, unit) {
    const ratio = mapDistanceToMeters[unit];
    if (ratio === void 0) {
      throw new Error(createError("INVALID_DISTANCE_SUFFIX", distance).message);
    }
    return distance * ratio;
  }
  function removeVectorsFromHits(searchResult, vectorProperties) {
    searchResult.hits = searchResult.hits.map((result) => ({
      ...result,
      document: {
        ...result.document,
        // Remove embeddings from the result
        ...vectorProperties.reduce((acc, prop) => {
          const path2 = prop.split(".");
          const lastKey = path2.pop();
          let obj = acc;
          for (const key of path2) {
            obj[key] = obj[key] ?? {};
            obj = obj[key];
          }
          obj[lastKey] = null;
          return acc;
        }, result.document)
      }
    }));
  }
  function isAsyncFunction(func) {
    if (Array.isArray(func)) {
      return func.some((item) => isAsyncFunction(item));
    }
    return func?.constructor?.name === "AsyncFunction";
  }
  var withIntersection = "intersection" in /* @__PURE__ */ new Set();
  function setIntersection(...sets) {
    if (sets.length === 0) {
      return /* @__PURE__ */ new Set();
    }
    if (sets.length === 1) {
      return sets[0];
    }
    if (sets.length === 2) {
      const set1 = sets[0];
      const set2 = sets[1];
      if (withIntersection) {
        return set1.intersection(set2);
      }
      const result = /* @__PURE__ */ new Set();
      const base2 = set1.size < set2.size ? set1 : set2;
      const other = base2 === set1 ? set2 : set1;
      for (const value of base2) {
        if (other.has(value)) {
          result.add(value);
        }
      }
      return result;
    }
    const min = {
      index: 0,
      size: sets[0].size
    };
    for (let i = 1; i < sets.length; i++) {
      if (sets[i].size < min.size) {
        min.index = i;
        min.size = sets[i].size;
      }
    }
    if (withIntersection) {
      let base2 = sets[min.index];
      for (let i = 0; i < sets.length; i++) {
        if (i === min.index) {
          continue;
        }
        base2 = base2.intersection(sets[i]);
      }
      return base2;
    }
    const base = sets[min.index];
    for (let i = 0; i < sets.length; i++) {
      if (i === min.index) {
        continue;
      }
      const other = sets[i];
      for (const value of base) {
        if (!other.has(value)) {
          base.delete(value);
        }
      }
    }
    return base;
  }
  var withUnion = "union" in /* @__PURE__ */ new Set();
  function setUnion(set1, set2) {
    if (withUnion) {
      if (set1) {
        return set1.union(set2);
      }
      return set2;
    }
    if (!set1) {
      return new Set(set2);
    }
    return /* @__PURE__ */ new Set([...set1, ...set2]);
  }
  function setDifference(set1, set2) {
    const result = /* @__PURE__ */ new Set();
    for (const value of set1) {
      if (!set2.has(value)) {
        result.add(value);
      }
    }
    return result;
  }
  function sleep(ms) {
    if (typeof SharedArrayBuffer !== "undefined" && typeof Atomics !== "undefined") {
      const nil = new Int32Array(new SharedArrayBuffer(4));
      const valid = ms > 0 && ms < Infinity;
      if (valid === false) {
        if (typeof ms !== "number" && typeof ms !== "bigint") {
          throw TypeError("sleep: ms must be a number");
        }
        throw RangeError("sleep: ms must be a number that is greater than 0 but less than Infinity");
      }
      Atomics.wait(nil, 0, 0, Number(ms));
    } else {
      const valid = ms > 0 && ms < Infinity;
      if (valid === false) {
        if (typeof ms !== "number" && typeof ms !== "bigint") {
          throw TypeError("sleep: ms must be a number");
        }
        throw RangeError("sleep: ms must be a number that is greater than 0 but less than Infinity");
      }
      const target = Date.now() + Number(ms);
      while (target > Date.now()) {
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/errors.js
  var allLanguages = SUPPORTED_LANGUAGES.join("\n - ");
  var errors = {
    NO_LANGUAGE_WITH_CUSTOM_TOKENIZER: "Do not pass the language option to create when using a custom tokenizer.",
    LANGUAGE_NOT_SUPPORTED: `Language "%s" is not supported.
Supported languages are:
 - ${allLanguages}`,
    INVALID_STEMMER_FUNCTION_TYPE: `config.stemmer property must be a function.`,
    MISSING_STEMMER: `As of version 1.0.0 @orama/orama does not ship non English stemmers by default. To solve this, please explicitly import and specify the "%s" stemmer from the package @orama/stemmers. See https://docs.orama.com/docs/orama-js/text-analysis/stemming for more information.`,
    CUSTOM_STOP_WORDS_MUST_BE_FUNCTION_OR_ARRAY: "Custom stop words array must only contain strings.",
    UNSUPPORTED_COMPONENT: `Unsupported component "%s".`,
    COMPONENT_MUST_BE_FUNCTION: `The component "%s" must be a function.`,
    COMPONENT_MUST_BE_FUNCTION_OR_ARRAY_FUNCTIONS: `The component "%s" must be a function or an array of functions.`,
    INVALID_SCHEMA_TYPE: `Unsupported schema type "%s" at "%s". Expected "string", "boolean" or "number" or array of them.`,
    DOCUMENT_ID_MUST_BE_STRING: `Document id must be of type "string". Got "%s" instead.`,
    DOCUMENT_ALREADY_EXISTS: `A document with id "%s" already exists.`,
    DOCUMENT_DOES_NOT_EXIST: `A document with id "%s" does not exists.`,
    MISSING_DOCUMENT_PROPERTY: `Missing searchable property "%s".`,
    INVALID_DOCUMENT_PROPERTY: `Invalid document property "%s": expected "%s", got "%s"`,
    UNKNOWN_INDEX: `Invalid property name "%s". Expected a wildcard string ("*") or array containing one of the following properties: %s`,
    INVALID_BOOST_VALUE: `Boost value must be a number greater than, or less than 0.`,
    INVALID_FILTER_OPERATION: `You can only use one operation per filter, you requested %d.`,
    SCHEMA_VALIDATION_FAILURE: `Cannot insert document due schema validation failure on "%s" property.`,
    INVALID_SORT_SCHEMA_TYPE: `Unsupported sort schema type "%s" at "%s". Expected "string" or "number".`,
    CANNOT_SORT_BY_ARRAY: `Cannot configure sort for "%s" because it is an array (%s).`,
    UNABLE_TO_SORT_ON_UNKNOWN_FIELD: `Unable to sort on unknown field "%s". Allowed fields: %s`,
    SORT_DISABLED: `Sort is disabled. Please read the documentation at https://docs.orama.com/docs/orama-js for more information.`,
    UNKNOWN_GROUP_BY_PROPERTY: `Unknown groupBy property "%s".`,
    INVALID_GROUP_BY_PROPERTY: `Invalid groupBy property "%s". Allowed types: "%s", but given "%s".`,
    UNKNOWN_FILTER_PROPERTY: `Unknown filter property "%s".`,
    UNKNOWN_VECTOR_PROPERTY: `Unknown vector property "%s". Make sure the property exists in the schema and is configured as a vector.`,
    INVALID_VECTOR_SIZE: `Vector size must be a number greater than 0. Got "%s" instead.`,
    INVALID_VECTOR_VALUE: `Vector value must be a number greater than 0. Got "%s" instead.`,
    INVALID_INPUT_VECTOR: `Property "%s" was declared as a %s-dimensional vector, but got a %s-dimensional vector instead.
Input vectors must be of the size declared in the schema, as calculating similarity between vectors of different sizes can lead to unexpected results.`,
    WRONG_SEARCH_PROPERTY_TYPE: `Property "%s" is not searchable. Only "string" properties are searchable.`,
    FACET_NOT_SUPPORTED: `Facet doens't support the type "%s".`,
    INVALID_DISTANCE_SUFFIX: `Invalid distance suffix "%s". Valid suffixes are: cm, m, km, mi, yd, ft.`,
    INVALID_SEARCH_MODE: `Invalid search mode "%s". Valid modes are: "fulltext", "vector", "hybrid".`,
    MISSING_VECTOR_AND_SECURE_PROXY: `No vector was provided and no secure proxy was configured. Please provide a vector or configure an Orama Secure Proxy to perform hybrid search.`,
    MISSING_TERM: `"term" is a required parameter when performing hybrid search. Please provide a search term.`,
    INVALID_VECTOR_INPUT: `Invalid "vector" property. Expected an object with "value" and "property" properties, but got "%s" instead.`,
    PLUGIN_CRASHED: `A plugin crashed during initialization. Please check the error message for more information:`,
    PLUGIN_SECURE_PROXY_NOT_FOUND: `Could not find '@orama/secure-proxy-plugin' installed in your Orama instance.
Please install it before proceeding with creating an answer session.
Read more at https://docs.orama.com/docs/orama-js/plugins/plugin-secure-proxy#plugin-secure-proxy
`,
    PLUGIN_SECURE_PROXY_MISSING_CHAT_MODEL: `Could not find a chat model defined in the secure proxy plugin configuration.
Please provide a chat model before proceeding with creating an answer session.
Read more at https://docs.orama.com/docs/orama-js/plugins/plugin-secure-proxy#plugin-secure-proxy
`,
    ANSWER_SESSION_LAST_MESSAGE_IS_NOT_ASSISTANT: `The last message in the session is not an assistant message. Cannot regenerate non-assistant messages.`,
    PLUGIN_COMPONENT_CONFLICT: `The component "%s" is already defined. The plugin "%s" is trying to redefine it.`
  };
  function createError(code3, ...args) {
    const error = new Error(sprintf(errors[code3] ?? `Unsupported Orama Error code: ${code3}`, ...args));
    error.code = code3;
    if ("captureStackTrace" in Error.prototype) {
      Error.captureStackTrace(error);
    }
    return error;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/defaults.js
  function formatElapsedTime(n) {
    return {
      raw: Number(n),
      formatted: formatNanoseconds(n)
    };
  }
  function getDocumentIndexId(doc) {
    if (doc.id) {
      if (typeof doc.id !== "string") {
        throw createError("DOCUMENT_ID_MUST_BE_STRING", typeof doc.id);
      }
      return doc.id;
    }
    return uniqueId();
  }
  function validateSchema(doc, schema) {
    for (const [prop, type] of Object.entries(schema)) {
      const value = doc[prop];
      if (typeof value === "undefined") {
        continue;
      }
      if (type === "geopoint" && typeof value === "object" && typeof value.lon === "number" && typeof value.lat === "number") {
        continue;
      }
      if (type === "enum" && (typeof value === "string" || typeof value === "number")) {
        continue;
      }
      if (type === "enum[]" && Array.isArray(value)) {
        const valueLength = value.length;
        for (let i = 0; i < valueLength; i++) {
          if (typeof value[i] !== "string" && typeof value[i] !== "number") {
            return prop + "." + i;
          }
        }
        continue;
      }
      if (isVectorType(type)) {
        const vectorSize = getVectorSize(type);
        if (!Array.isArray(value) || value.length !== vectorSize) {
          throw createError("INVALID_INPUT_VECTOR", prop, vectorSize, value.length);
        }
        continue;
      }
      if (isArrayType(type)) {
        if (!Array.isArray(value)) {
          return prop;
        }
        const expectedType = getInnerType(type);
        const valueLength = value.length;
        for (let i = 0; i < valueLength; i++) {
          if (typeof value[i] !== expectedType) {
            return prop + "." + i;
          }
        }
        continue;
      }
      if (typeof type === "object") {
        if (!value || typeof value !== "object") {
          return prop;
        }
        const subProp = validateSchema(value, type);
        if (subProp) {
          return prop + "." + subProp;
        }
        continue;
      }
      if (typeof value !== type) {
        return prop;
      }
    }
    return void 0;
  }
  var IS_ARRAY_TYPE = {
    string: false,
    number: false,
    boolean: false,
    enum: false,
    geopoint: false,
    "string[]": true,
    "number[]": true,
    "boolean[]": true,
    "enum[]": true
  };
  var INNER_TYPE = {
    "string[]": "string",
    "number[]": "number",
    "boolean[]": "boolean",
    "enum[]": "enum"
  };
  function isGeoPointType(type) {
    return type === "geopoint";
  }
  function isVectorType(type) {
    return typeof type === "string" && /^vector\[\d+\]$/.test(type);
  }
  function isArrayType(type) {
    return typeof type === "string" && IS_ARRAY_TYPE[type];
  }
  function getInnerType(type) {
    return INNER_TYPE[type];
  }
  function getVectorSize(type) {
    const size = Number(type.slice(7, -1));
    switch (true) {
      case isNaN(size):
        throw createError("INVALID_VECTOR_VALUE", type);
      case size <= 0:
        throw createError("INVALID_VECTOR_SIZE", type);
      default:
        return size;
    }
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/internal-document-id-store.js
  function createInternalDocumentIDStore() {
    return {
      idToInternalId: /* @__PURE__ */ new Map(),
      internalIdToId: [],
      save,
      load
    };
  }
  function save(store2) {
    return {
      internalIdToId: store2.internalIdToId
    };
  }
  function load(orama, raw) {
    const { internalIdToId } = raw;
    orama.internalDocumentIDStore.idToInternalId.clear();
    orama.internalDocumentIDStore.internalIdToId = [];
    const internalIdToIdLength = internalIdToId.length;
    for (let i = 0; i < internalIdToIdLength; i++) {
      const internalIdItem = internalIdToId[i];
      orama.internalDocumentIDStore.idToInternalId.set(internalIdItem, i + 1);
      orama.internalDocumentIDStore.internalIdToId.push(internalIdItem);
    }
  }
  function getInternalDocumentId(store2, id) {
    if (typeof id === "string") {
      const internalId = store2.idToInternalId.get(id);
      if (internalId) {
        return internalId;
      }
      const currentId = store2.idToInternalId.size + 1;
      store2.idToInternalId.set(id, currentId);
      store2.internalIdToId.push(id);
      return currentId;
    }
    if (id > store2.internalIdToId.length) {
      return getInternalDocumentId(store2, id.toString());
    }
    return id;
  }
  function getDocumentIdFromInternalId(store2, internalId) {
    if (store2.internalIdToId.length < internalId) {
      throw new Error(`Invalid internalId ${internalId}`);
    }
    return store2.internalIdToId[internalId - 1];
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/documents-store.js
  function create(_, sharedInternalDocumentStore) {
    return {
      sharedInternalDocumentStore,
      docs: {},
      count: 0
    };
  }
  function get(store2, id) {
    const internalId = getInternalDocumentId(store2.sharedInternalDocumentStore, id);
    return store2.docs[internalId];
  }
  function getMultiple(store2, ids) {
    const idsLength = ids.length;
    const found = Array.from({ length: idsLength });
    for (let i = 0; i < idsLength; i++) {
      const internalId = getInternalDocumentId(store2.sharedInternalDocumentStore, ids[i]);
      found[i] = store2.docs[internalId];
    }
    return found;
  }
  function getAll(store2) {
    return store2.docs;
  }
  function store(store2, id, internalId, doc) {
    if (typeof store2.docs[internalId] !== "undefined") {
      return false;
    }
    store2.docs[internalId] = doc;
    store2.count++;
    return true;
  }
  function remove(store2, id) {
    const internalId = getInternalDocumentId(store2.sharedInternalDocumentStore, id);
    if (typeof store2.docs[internalId] === "undefined") {
      return false;
    }
    delete store2.docs[internalId];
    store2.count--;
    return true;
  }
  function count(store2) {
    return store2.count;
  }
  function load2(sharedInternalDocumentStore, raw) {
    const rawDocument = raw;
    return {
      docs: rawDocument.docs,
      count: rawDocument.count,
      sharedInternalDocumentStore
    };
  }
  function save2(store2) {
    return {
      docs: store2.docs,
      count: store2.count
    };
  }
  function createDocumentsStore() {
    return {
      create,
      get,
      getMultiple,
      getAll,
      store,
      remove,
      count,
      load: load2,
      save: save2
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/plugins.js
  var AVAILABLE_PLUGIN_HOOKS = [
    "beforeInsert",
    "afterInsert",
    "beforeRemove",
    "afterRemove",
    "beforeUpdate",
    "afterUpdate",
    "beforeUpsert",
    "afterUpsert",
    "beforeSearch",
    "afterSearch",
    "beforeInsertMultiple",
    "afterInsertMultiple",
    "beforeRemoveMultiple",
    "afterRemoveMultiple",
    "beforeUpdateMultiple",
    "afterUpdateMultiple",
    "beforeUpsertMultiple",
    "afterUpsertMultiple",
    "beforeLoad",
    "afterLoad",
    "afterCreate"
  ];
  function getAllPluginsByHook(orama, hook) {
    const pluginsToRun = [];
    const pluginsLength = orama.plugins?.length;
    if (!pluginsLength) {
      return pluginsToRun;
    }
    for (let i = 0; i < pluginsLength; i++) {
      try {
        const plugin = orama.plugins[i];
        if (typeof plugin[hook] === "function") {
          pluginsToRun.push(plugin[hook]);
        }
      } catch (error) {
        console.error("Caught error in getAllPluginsByHook:", error);
        throw createError("PLUGIN_CRASHED");
      }
    }
    return pluginsToRun;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/hooks.js
  var OBJECT_COMPONENTS = ["tokenizer", "index", "documentsStore", "sorter", "pinning"];
  var FUNCTION_COMPONENTS = [
    "validateSchema",
    "getDocumentIndexId",
    "getDocumentProperties",
    "formatElapsedTime"
  ];
  function runSingleHook(hooks, orama, id, doc) {
    const needAsync = hooks.some(isAsyncFunction);
    if (needAsync) {
      return (async () => {
        for (const hook of hooks) {
          await hook(orama, id, doc);
        }
      })();
    } else {
      for (const hook of hooks) {
        hook(orama, id, doc);
      }
    }
  }
  function runMultipleHook(hooks, orama, docsOrIds) {
    const needAsync = hooks.some(isAsyncFunction);
    if (needAsync) {
      return (async () => {
        for (const hook of hooks) {
          await hook(orama, docsOrIds);
        }
      })();
    } else {
      for (const hook of hooks) {
        hook(orama, docsOrIds);
      }
    }
  }
  function runAfterSearch(hooks, db, params, language, results) {
    const needAsync = hooks.some(isAsyncFunction);
    if (needAsync) {
      return (async () => {
        for (const hook of hooks) {
          await hook(db, params, language, results);
        }
      })();
    } else {
      for (const hook of hooks) {
        hook(db, params, language, results);
      }
    }
  }
  function runBeforeSearch(hooks, db, params, language) {
    const needAsync = hooks.some(isAsyncFunction);
    if (needAsync) {
      return (async () => {
        for (const hook of hooks) {
          await hook(db, params, language);
        }
      })();
    } else {
      for (const hook of hooks) {
        hook(db, params, language);
      }
    }
  }
  function runAfterCreate(hooks, db) {
    const needAsync = hooks.some(isAsyncFunction);
    if (needAsync) {
      return (async () => {
        for (const hook of hooks) {
          await hook(db);
        }
      })();
    } else {
      for (const hook of hooks) {
        hook(db);
      }
    }
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/avl.js
  var AVLNode = class _AVLNode {
    constructor(key, value) {
      __publicField(this, "k");
      __publicField(this, "v");
      __publicField(this, "l", null);
      __publicField(this, "r", null);
      __publicField(this, "h", 1);
      this.k = key;
      this.v = new Set(value);
    }
    updateHeight() {
      this.h = Math.max(_AVLNode.getHeight(this.l), _AVLNode.getHeight(this.r)) + 1;
    }
    static getHeight(node2) {
      return node2 ? node2.h : 0;
    }
    getBalanceFactor() {
      return _AVLNode.getHeight(this.l) - _AVLNode.getHeight(this.r);
    }
    rotateLeft() {
      const newRoot = this.r;
      this.r = newRoot.l;
      newRoot.l = this;
      this.updateHeight();
      newRoot.updateHeight();
      return newRoot;
    }
    rotateRight() {
      const newRoot = this.l;
      this.l = newRoot.r;
      newRoot.r = this;
      this.updateHeight();
      newRoot.updateHeight();
      return newRoot;
    }
    toJSON() {
      return {
        k: this.k,
        v: Array.from(this.v),
        l: this.l ? this.l.toJSON() : null,
        r: this.r ? this.r.toJSON() : null,
        h: this.h
      };
    }
    static fromJSON(json) {
      const node2 = new _AVLNode(json.k, json.v);
      node2.l = json.l ? _AVLNode.fromJSON(json.l) : null;
      node2.r = json.r ? _AVLNode.fromJSON(json.r) : null;
      node2.h = json.h;
      return node2;
    }
  };
  var AVLTree = class _AVLTree {
    constructor(key, value) {
      __publicField(this, "root", null);
      __publicField(this, "insertCount", 0);
      if (key !== void 0 && value !== void 0) {
        this.root = new AVLNode(key, value);
      }
    }
    insert(key, value, rebalanceThreshold = 1e3) {
      this.root = this.insertNode(this.root, key, value, rebalanceThreshold);
    }
    insertMultiple(key, value, rebalanceThreshold = 1e3) {
      for (const v2 of value) {
        this.insert(key, v2, rebalanceThreshold);
      }
    }
    // Rebalance the tree if the insert count reaches the threshold.
    // This will improve insertion performance since we won't be rebalancing the tree on every insert.
    // When inserting docs using `insertMultiple`, the threshold will be set to the number of docs being inserted.
    // We can force rebalancing the tree by setting the threshold to 1 (default).
    rebalance() {
      if (this.root) {
        this.root = this.rebalanceNode(this.root);
      }
    }
    toJSON() {
      return {
        root: this.root ? this.root.toJSON() : null,
        insertCount: this.insertCount
      };
    }
    static fromJSON(json) {
      const tree = new _AVLTree();
      tree.root = json.root ? AVLNode.fromJSON(json.root) : null;
      tree.insertCount = json.insertCount || 0;
      return tree;
    }
    insertNode(node2, key, value, rebalanceThreshold) {
      if (node2 === null) {
        return new AVLNode(key, [value]);
      }
      const path2 = [];
      let current2 = node2;
      let parent = null;
      while (current2 !== null) {
        path2.push({ parent, node: current2 });
        if (key < current2.k) {
          if (current2.l === null) {
            current2.l = new AVLNode(key, [value]);
            path2.push({ parent: current2, node: current2.l });
            break;
          } else {
            parent = current2;
            current2 = current2.l;
          }
        } else if (key > current2.k) {
          if (current2.r === null) {
            current2.r = new AVLNode(key, [value]);
            path2.push({ parent: current2, node: current2.r });
            break;
          } else {
            parent = current2;
            current2 = current2.r;
          }
        } else {
          current2.v.add(value);
          return node2;
        }
      }
      let needRebalance = false;
      if (this.insertCount++ % rebalanceThreshold === 0) {
        needRebalance = true;
      }
      for (let i = path2.length - 1; i >= 0; i--) {
        const { parent: parent2, node: currentNode } = path2[i];
        currentNode.updateHeight();
        if (needRebalance) {
          const rebalancedNode = this.rebalanceNode(currentNode);
          if (parent2) {
            if (parent2.l === currentNode) {
              parent2.l = rebalancedNode;
            } else if (parent2.r === currentNode) {
              parent2.r = rebalancedNode;
            }
          } else {
            node2 = rebalancedNode;
          }
        }
      }
      return node2;
    }
    rebalanceNode(node2) {
      const balanceFactor = node2.getBalanceFactor();
      if (balanceFactor > 1) {
        if (node2.l && node2.l.getBalanceFactor() >= 0) {
          return node2.rotateRight();
        } else if (node2.l) {
          node2.l = node2.l.rotateLeft();
          return node2.rotateRight();
        }
      }
      if (balanceFactor < -1) {
        if (node2.r && node2.r.getBalanceFactor() <= 0) {
          return node2.rotateLeft();
        } else if (node2.r) {
          node2.r = node2.r.rotateRight();
          return node2.rotateLeft();
        }
      }
      return node2;
    }
    find(key) {
      const node2 = this.findNodeByKey(key);
      return node2 ? node2.v : null;
    }
    contains(key) {
      return this.find(key) !== null;
    }
    getSize() {
      let count3 = 0;
      const stack = [];
      let current2 = this.root;
      while (current2 || stack.length > 0) {
        while (current2) {
          stack.push(current2);
          current2 = current2.l;
        }
        current2 = stack.pop();
        count3++;
        current2 = current2.r;
      }
      return count3;
    }
    isBalanced() {
      if (!this.root)
        return true;
      const stack = [this.root];
      while (stack.length > 0) {
        const node2 = stack.pop();
        const balanceFactor = node2.getBalanceFactor();
        if (Math.abs(balanceFactor) > 1) {
          return false;
        }
        if (node2.l)
          stack.push(node2.l);
        if (node2.r)
          stack.push(node2.r);
      }
      return true;
    }
    remove(key) {
      this.root = this.removeNode(this.root, key);
    }
    removeDocument(key, id) {
      const node2 = this.findNodeByKey(key);
      if (!node2) {
        return;
      }
      if (node2.v.size === 1) {
        this.root = this.removeNode(this.root, key);
      } else {
        node2.v = new Set([...node2.v.values()].filter((v2) => v2 !== id));
      }
    }
    findNodeByKey(key) {
      let node2 = this.root;
      while (node2) {
        if (key < node2.k) {
          node2 = node2.l;
        } else if (key > node2.k) {
          node2 = node2.r;
        } else {
          return node2;
        }
      }
      return null;
    }
    removeNode(node2, key) {
      if (node2 === null)
        return null;
      const path2 = [];
      let current2 = node2;
      while (current2 !== null && current2.k !== key) {
        path2.push(current2);
        if (key < current2.k) {
          current2 = current2.l;
        } else {
          current2 = current2.r;
        }
      }
      if (current2 === null) {
        return node2;
      }
      if (current2.l === null || current2.r === null) {
        const child = current2.l ? current2.l : current2.r;
        if (path2.length === 0) {
          node2 = child;
        } else {
          const parent = path2[path2.length - 1];
          if (parent.l === current2) {
            parent.l = child;
          } else {
            parent.r = child;
          }
        }
      } else {
        let successorParent = current2;
        let successor = current2.r;
        while (successor.l !== null) {
          successorParent = successor;
          successor = successor.l;
        }
        current2.k = successor.k;
        current2.v = successor.v;
        if (successorParent.l === successor) {
          successorParent.l = successor.r;
        } else {
          successorParent.r = successor.r;
        }
        current2 = successorParent;
      }
      path2.push(current2);
      for (let i = path2.length - 1; i >= 0; i--) {
        const currentNode = path2[i];
        currentNode.updateHeight();
        const rebalancedNode = this.rebalanceNode(currentNode);
        if (i > 0) {
          const parent = path2[i - 1];
          if (parent.l === currentNode) {
            parent.l = rebalancedNode;
          } else if (parent.r === currentNode) {
            parent.r = rebalancedNode;
          }
        } else {
          node2 = rebalancedNode;
        }
      }
      return node2;
    }
    rangeSearch(min, max) {
      const result = /* @__PURE__ */ new Set();
      const stack = [];
      let current2 = this.root;
      while (current2 || stack.length > 0) {
        while (current2) {
          stack.push(current2);
          current2 = current2.l;
        }
        current2 = stack.pop();
        if (current2.k >= min && current2.k <= max) {
          for (const value of current2.v) {
            result.add(value);
          }
        }
        if (current2.k > max) {
          break;
        }
        current2 = current2.r;
      }
      return result;
    }
    greaterThan(key, inclusive = false) {
      const result = /* @__PURE__ */ new Set();
      const stack = [];
      let current2 = this.root;
      while (current2 || stack.length > 0) {
        while (current2) {
          stack.push(current2);
          current2 = current2.r;
        }
        current2 = stack.pop();
        if (inclusive && current2.k >= key || !inclusive && current2.k > key) {
          for (const value of current2.v) {
            result.add(value);
          }
        } else if (current2.k <= key) {
          break;
        }
        current2 = current2.l;
      }
      return result;
    }
    lessThan(key, inclusive = false) {
      const result = /* @__PURE__ */ new Set();
      const stack = [];
      let current2 = this.root;
      while (current2 || stack.length > 0) {
        while (current2) {
          stack.push(current2);
          current2 = current2.l;
        }
        current2 = stack.pop();
        if (inclusive && current2.k <= key || !inclusive && current2.k < key) {
          for (const value of current2.v) {
            result.add(value);
          }
        } else if (current2.k > key) {
          break;
        }
        current2 = current2.r;
      }
      return result;
    }
  };

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/flat.js
  var FlatTree = class _FlatTree {
    constructor() {
      __publicField(this, "numberToDocumentId");
      this.numberToDocumentId = /* @__PURE__ */ new Map();
    }
    insert(key, value) {
      if (this.numberToDocumentId.has(key)) {
        this.numberToDocumentId.get(key).add(value);
      } else {
        this.numberToDocumentId.set(key, /* @__PURE__ */ new Set([value]));
      }
    }
    find(key) {
      const idSet = this.numberToDocumentId.get(key);
      return idSet ? Array.from(idSet) : null;
    }
    remove(key) {
      this.numberToDocumentId.delete(key);
    }
    removeDocument(id, key) {
      const idSet = this.numberToDocumentId.get(key);
      if (idSet) {
        idSet.delete(id);
        if (idSet.size === 0) {
          this.numberToDocumentId.delete(key);
        }
      }
    }
    contains(key) {
      return this.numberToDocumentId.has(key);
    }
    getSize() {
      let size = 0;
      for (const idSet of this.numberToDocumentId.values()) {
        size += idSet.size;
      }
      return size;
    }
    filter(operation) {
      const operationKeys = Object.keys(operation);
      if (operationKeys.length !== 1) {
        throw new Error("Invalid operation");
      }
      const operationType = operationKeys[0];
      switch (operationType) {
        case "eq": {
          const value = operation[operationType];
          const idSet = this.numberToDocumentId.get(value);
          return idSet ? Array.from(idSet) : [];
        }
        case "in": {
          const values = operation[operationType];
          const resultSet = /* @__PURE__ */ new Set();
          for (const value of values) {
            const idSet = this.numberToDocumentId.get(value);
            if (idSet) {
              for (const id of idSet) {
                resultSet.add(id);
              }
            }
          }
          return Array.from(resultSet);
        }
        case "nin": {
          const excludeValues = new Set(operation[operationType]);
          const resultSet = /* @__PURE__ */ new Set();
          for (const [key, idSet] of this.numberToDocumentId.entries()) {
            if (!excludeValues.has(key)) {
              for (const id of idSet) {
                resultSet.add(id);
              }
            }
          }
          return Array.from(resultSet);
        }
        default:
          throw new Error("Invalid operation");
      }
    }
    filterArr(operation) {
      const operationKeys = Object.keys(operation);
      if (operationKeys.length !== 1) {
        throw new Error("Invalid operation");
      }
      const operationType = operationKeys[0];
      switch (operationType) {
        case "containsAll": {
          const values = operation[operationType];
          const idSets = values.map((value) => this.numberToDocumentId.get(value) ?? /* @__PURE__ */ new Set());
          if (idSets.length === 0)
            return [];
          const intersection = idSets.reduce((prev, curr) => {
            return new Set([...prev].filter((id) => curr.has(id)));
          });
          return Array.from(intersection);
        }
        case "containsAny": {
          const values = operation[operationType];
          const idSets = values.map((value) => this.numberToDocumentId.get(value) ?? /* @__PURE__ */ new Set());
          if (idSets.length === 0)
            return [];
          const union = idSets.reduce((prev, curr) => {
            return /* @__PURE__ */ new Set([...prev, ...curr]);
          });
          return Array.from(union);
        }
        default:
          throw new Error("Invalid operation");
      }
    }
    static fromJSON(json) {
      if (!json.numberToDocumentId) {
        throw new Error("Invalid Flat Tree JSON");
      }
      const tree = new _FlatTree();
      for (const [key, ids] of json.numberToDocumentId) {
        tree.numberToDocumentId.set(key, new Set(ids));
      }
      return tree;
    }
    toJSON() {
      return {
        numberToDocumentId: Array.from(this.numberToDocumentId.entries()).map(([key, idSet]) => [key, Array.from(idSet)])
      };
    }
  };

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/levenshtein.js
  function _boundedLevenshtein(term, word, tolerance) {
    if (tolerance < 0)
      return -1;
    if (term === word)
      return 0;
    const m = term.length;
    const n = word.length;
    if (m === 0)
      return n <= tolerance ? n : -1;
    if (n === 0)
      return m <= tolerance ? m : -1;
    const diff = Math.abs(m - n);
    if (term.startsWith(word)) {
      return diff <= tolerance ? diff : -1;
    }
    if (word.startsWith(term)) {
      return 0;
    }
    if (diff > tolerance)
      return -1;
    const matrix = [];
    for (let i = 0; i <= m; i++) {
      matrix[i] = [i];
      for (let j = 1; j <= n; j++) {
        matrix[i][j] = i === 0 ? j : 0;
      }
    }
    for (let i = 1; i <= m; i++) {
      let rowMin = Infinity;
      for (let j = 1; j <= n; j++) {
        if (term[i - 1] === word[j - 1]) {
          matrix[i][j] = matrix[i - 1][j - 1];
        } else {
          matrix[i][j] = Math.min(
            matrix[i - 1][j] + 1,
            // deletion
            matrix[i][j - 1] + 1,
            // insertion
            matrix[i - 1][j - 1] + 1
            // substitution
          );
        }
        rowMin = Math.min(rowMin, matrix[i][j]);
      }
      if (rowMin > tolerance) {
        return -1;
      }
    }
    return matrix[m][n] <= tolerance ? matrix[m][n] : -1;
  }
  function syncBoundedLevenshtein(term, w, tolerance) {
    const distance = _boundedLevenshtein(term, w, tolerance);
    return {
      distance,
      isBounded: distance >= 0
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/radix.js
  var RadixNode = class _RadixNode {
    constructor(key, subWord, end) {
      // Node key
      __publicField(this, "k");
      // Node subword
      __publicField(this, "s");
      // Node children
      __publicField(this, "c", /* @__PURE__ */ new Map());
      // Node documents
      __publicField(this, "d", /* @__PURE__ */ new Set());
      // Node end
      __publicField(this, "e");
      // Node word
      __publicField(this, "w", "");
      this.k = key;
      this.s = subWord;
      this.e = end;
    }
    updateParent(parent) {
      this.w = parent.w + this.s;
    }
    addDocument(docID) {
      this.d.add(docID);
    }
    removeDocument(docID) {
      return this.d.delete(docID);
    }
    findAllWords(output, term, exact, tolerance) {
      const stack = [this];
      while (stack.length > 0) {
        const node2 = stack.pop();
        if (node2.e) {
          const { w, d: docIDs } = node2;
          if (exact && w !== term) {
            continue;
          }
          if (getOwnProperty(output, w) !== null) {
            if (tolerance) {
              const difference = Math.abs(term.length - w.length);
              if (difference <= tolerance && syncBoundedLevenshtein(term, w, tolerance).isBounded) {
                output[w] = [];
              } else {
                continue;
              }
            } else {
              output[w] = [];
            }
          }
          if (getOwnProperty(output, w) != null && docIDs.size > 0) {
            const docs = output[w];
            for (const docID of docIDs) {
              if (!docs.includes(docID)) {
                docs.push(docID);
              }
            }
          }
        }
        if (node2.c.size > 0) {
          stack.push(...node2.c.values());
        }
      }
      return output;
    }
    insert(word, docId) {
      let node2 = this;
      let i = 0;
      const wordLength = word.length;
      while (i < wordLength) {
        const currentCharacter = word[i];
        const childNode = node2.c.get(currentCharacter);
        if (childNode) {
          const edgeLabel = childNode.s;
          const edgeLabelLength = edgeLabel.length;
          let j = 0;
          while (j < edgeLabelLength && i + j < wordLength && edgeLabel[j] === word[i + j]) {
            j++;
          }
          if (j === edgeLabelLength) {
            node2 = childNode;
            i += j;
            if (i === wordLength) {
              if (!childNode.e) {
                childNode.e = true;
              }
              childNode.addDocument(docId);
              return;
            }
            continue;
          }
          const commonPrefix = edgeLabel.slice(0, j);
          const newEdgeLabel = edgeLabel.slice(j);
          const newWordLabel = word.slice(i + j);
          const inbetweenNode = new _RadixNode(commonPrefix[0], commonPrefix, false);
          node2.c.set(commonPrefix[0], inbetweenNode);
          inbetweenNode.updateParent(node2);
          childNode.s = newEdgeLabel;
          childNode.k = newEdgeLabel[0];
          inbetweenNode.c.set(newEdgeLabel[0], childNode);
          childNode.updateParent(inbetweenNode);
          if (newWordLabel) {
            const newNode = new _RadixNode(newWordLabel[0], newWordLabel, true);
            newNode.addDocument(docId);
            inbetweenNode.c.set(newWordLabel[0], newNode);
            newNode.updateParent(inbetweenNode);
          } else {
            inbetweenNode.e = true;
            inbetweenNode.addDocument(docId);
          }
          return;
        } else {
          const newNode = new _RadixNode(currentCharacter, word.slice(i), true);
          newNode.addDocument(docId);
          node2.c.set(currentCharacter, newNode);
          newNode.updateParent(node2);
          return;
        }
      }
      if (!node2.e) {
        node2.e = true;
      }
      node2.addDocument(docId);
    }
    _findLevenshtein(term, index2, tolerance, originalTolerance, output) {
      const stack = [{ node: this, index: index2, tolerance }];
      while (stack.length > 0) {
        const { node: node2, index: index3, tolerance: tolerance2 } = stack.pop();
        if (node2.w.startsWith(term)) {
          node2.findAllWords(output, term, false, 0);
          continue;
        }
        if (tolerance2 < 0) {
          continue;
        }
        if (node2.e) {
          const { w, d: docIDs } = node2;
          if (w) {
            if (syncBoundedLevenshtein(term, w, originalTolerance).isBounded) {
              output[w] = [];
            }
            if (getOwnProperty(output, w) !== void 0 && docIDs.size > 0) {
              const docs = new Set(output[w]);
              for (const docID of docIDs) {
                docs.add(docID);
              }
              output[w] = Array.from(docs);
            }
          }
        }
        if (index3 >= term.length) {
          continue;
        }
        const currentChar = term[index3];
        if (node2.c.has(currentChar)) {
          const childNode = node2.c.get(currentChar);
          stack.push({ node: childNode, index: index3 + 1, tolerance: tolerance2 });
        }
        stack.push({ node: node2, index: index3 + 1, tolerance: tolerance2 - 1 });
        for (const [character, childNode] of node2.c) {
          stack.push({ node: childNode, index: index3, tolerance: tolerance2 - 1 });
          if (character !== currentChar) {
            stack.push({ node: childNode, index: index3 + 1, tolerance: tolerance2 - 1 });
          }
        }
      }
    }
    find(params) {
      const { term, exact, tolerance } = params;
      if (tolerance && !exact) {
        const output = {};
        this._findLevenshtein(term, 0, tolerance, tolerance, output);
        return output;
      } else {
        let node2 = this;
        let i = 0;
        const termLength = term.length;
        while (i < termLength) {
          const character = term[i];
          const childNode = node2.c.get(character);
          if (childNode) {
            const edgeLabel = childNode.s;
            const edgeLabelLength = edgeLabel.length;
            let j = 0;
            while (j < edgeLabelLength && i + j < termLength && edgeLabel[j] === term[i + j]) {
              j++;
            }
            if (j === edgeLabelLength) {
              node2 = childNode;
              i += j;
            } else if (i + j === termLength) {
              if (j === termLength - i) {
                if (exact) {
                  return {};
                } else {
                  const output2 = {};
                  childNode.findAllWords(output2, term, exact, tolerance);
                  return output2;
                }
              } else {
                return {};
              }
            } else {
              return {};
            }
          } else {
            return {};
          }
        }
        const output = {};
        node2.findAllWords(output, term, exact, tolerance);
        return output;
      }
    }
    contains(term) {
      let node2 = this;
      let i = 0;
      const termLength = term.length;
      while (i < termLength) {
        const character = term[i];
        const childNode = node2.c.get(character);
        if (childNode) {
          const edgeLabel = childNode.s;
          const edgeLabelLength = edgeLabel.length;
          let j = 0;
          while (j < edgeLabelLength && i + j < termLength && edgeLabel[j] === term[i + j]) {
            j++;
          }
          if (j < edgeLabelLength) {
            return false;
          }
          i += edgeLabelLength;
          node2 = childNode;
        } else {
          return false;
        }
      }
      return true;
    }
    removeWord(term) {
      if (!term) {
        return false;
      }
      let node2 = this;
      const termLength = term.length;
      const stack = [];
      for (let i = 0; i < termLength; i++) {
        const character = term[i];
        if (node2.c.has(character)) {
          const childNode = node2.c.get(character);
          stack.push({ parent: node2, character });
          i += childNode.s.length - 1;
          node2 = childNode;
        } else {
          return false;
        }
      }
      node2.d.clear();
      node2.e = false;
      while (stack.length > 0 && node2.c.size === 0 && !node2.e && node2.d.size === 0) {
        const { parent, character } = stack.pop();
        parent.c.delete(character);
        node2 = parent;
      }
      return true;
    }
    removeDocumentByWord(term, docID, exact = true) {
      if (!term) {
        return true;
      }
      let node2 = this;
      const termLength = term.length;
      for (let i = 0; i < termLength; i++) {
        const character = term[i];
        if (node2.c.has(character)) {
          const childNode = node2.c.get(character);
          i += childNode.s.length - 1;
          node2 = childNode;
          if (exact && node2.w !== term) {
          } else {
            node2.removeDocument(docID);
          }
        } else {
          return false;
        }
      }
      return true;
    }
    static getCommonPrefix(a, b) {
      const len = Math.min(a.length, b.length);
      let i = 0;
      while (i < len && a.charCodeAt(i) === b.charCodeAt(i)) {
        i++;
      }
      return a.slice(0, i);
    }
    toJSON() {
      return {
        w: this.w,
        s: this.s,
        e: this.e,
        k: this.k,
        d: Array.from(this.d),
        c: Array.from(this.c?.entries())?.map(([key, node2]) => [key, node2.toJSON()])
      };
    }
    static fromJSON(json) {
      const node2 = new _RadixNode(json.k, json.s, json.e);
      node2.w = json.w;
      node2.d = new Set(json.d);
      node2.c = new Map(json?.c?.map(([key, nodeJson]) => [key, _RadixNode.fromJSON(nodeJson)]) || []);
      return node2;
    }
  };
  var RadixTree = class _RadixTree extends RadixNode {
    constructor() {
      super("", "", false);
    }
    static fromJSON(json) {
      const tree = new _RadixTree();
      tree.w = json.w;
      tree.s = json.s;
      tree.e = json.e;
      tree.k = json.k;
      tree.d = new Set(json.d);
      tree.c = new Map(json?.c?.map(([key, nodeJson]) => [key, RadixNode.fromJSON(nodeJson)]) || []);
      return tree;
    }
    toJSON() {
      return super.toJSON();
    }
  };

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/bkd.js
  var K = 2;
  var EARTH_RADIUS = 6371e3;
  var BKDNode = class _BKDNode {
    constructor(point3, docIDs) {
      __publicField(this, "point");
      __publicField(this, "docIDs");
      __publicField(this, "left");
      __publicField(this, "right");
      __publicField(this, "parent");
      this.point = point3;
      this.docIDs = new Set(docIDs);
      this.left = null;
      this.right = null;
      this.parent = null;
    }
    toJSON() {
      return {
        point: this.point,
        docIDs: Array.from(this.docIDs),
        left: this.left ? this.left.toJSON() : null,
        right: this.right ? this.right.toJSON() : null
      };
    }
    static fromJSON(json, parent = null) {
      const node2 = new _BKDNode(json.point, json.docIDs);
      node2.parent = parent;
      if (json.left) {
        node2.left = _BKDNode.fromJSON(json.left, node2);
      }
      if (json.right) {
        node2.right = _BKDNode.fromJSON(json.right, node2);
      }
      return node2;
    }
  };
  var BKDTree = class _BKDTree {
    constructor() {
      __publicField(this, "root");
      __publicField(this, "nodeMap");
      this.root = null;
      this.nodeMap = /* @__PURE__ */ new Map();
    }
    getPointKey(point3) {
      return `${point3.lon},${point3.lat}`;
    }
    insert(point3, docIDs) {
      const pointKey = this.getPointKey(point3);
      const existingNode = this.nodeMap.get(pointKey);
      if (existingNode) {
        docIDs.forEach((id) => existingNode.docIDs.add(id));
        return;
      }
      const newNode = new BKDNode(point3, docIDs);
      this.nodeMap.set(pointKey, newNode);
      if (this.root == null) {
        this.root = newNode;
        return;
      }
      let node2 = this.root;
      let depth = 0;
      while (true) {
        const axis = depth % K;
        if (axis === 0) {
          if (point3.lon < node2.point.lon) {
            if (node2.left == null) {
              node2.left = newNode;
              newNode.parent = node2;
              return;
            }
            node2 = node2.left;
          } else {
            if (node2.right == null) {
              node2.right = newNode;
              newNode.parent = node2;
              return;
            }
            node2 = node2.right;
          }
        } else {
          if (point3.lat < node2.point.lat) {
            if (node2.left == null) {
              node2.left = newNode;
              newNode.parent = node2;
              return;
            }
            node2 = node2.left;
          } else {
            if (node2.right == null) {
              node2.right = newNode;
              newNode.parent = node2;
              return;
            }
            node2 = node2.right;
          }
        }
        depth++;
      }
    }
    contains(point3) {
      const pointKey = this.getPointKey(point3);
      return this.nodeMap.has(pointKey);
    }
    getDocIDsByCoordinates(point3) {
      const pointKey = this.getPointKey(point3);
      const node2 = this.nodeMap.get(pointKey);
      if (node2) {
        return Array.from(node2.docIDs);
      }
      return null;
    }
    removeDocByID(point3, docID) {
      const pointKey = this.getPointKey(point3);
      const node2 = this.nodeMap.get(pointKey);
      if (node2) {
        node2.docIDs.delete(docID);
        if (node2.docIDs.size === 0) {
          this.nodeMap.delete(pointKey);
          this.deleteNode(node2);
        }
      }
    }
    deleteNode(node2) {
      const parent = node2.parent;
      const child = node2.left ? node2.left : node2.right;
      if (child) {
        child.parent = parent;
      }
      if (parent) {
        if (parent.left === node2) {
          parent.left = child;
        } else if (parent.right === node2) {
          parent.right = child;
        }
      } else {
        this.root = child;
        if (this.root) {
          this.root.parent = null;
        }
      }
    }
    searchByRadius(center, radius, inclusive = true, sort = "asc", highPrecision = false) {
      const distanceFn = highPrecision ? _BKDTree.vincentyDistance : _BKDTree.haversineDistance;
      const stack = [{ node: this.root, depth: 0 }];
      const result = [];
      while (stack.length > 0) {
        const { node: node2, depth } = stack.pop();
        if (node2 == null)
          continue;
        const dist = distanceFn(center, node2.point);
        if (inclusive ? dist <= radius : dist > radius) {
          result.push({ point: node2.point, docIDs: Array.from(node2.docIDs) });
        }
        if (node2.left != null) {
          stack.push({ node: node2.left, depth: depth + 1 });
        }
        if (node2.right != null) {
          stack.push({ node: node2.right, depth: depth + 1 });
        }
      }
      if (sort) {
        result.sort((a, b) => {
          const distA = distanceFn(center, a.point);
          const distB = distanceFn(center, b.point);
          return sort.toLowerCase() === "asc" ? distA - distB : distB - distA;
        });
      }
      return result;
    }
    searchByPolygon(polygon, inclusive = true, sort = null, highPrecision = false) {
      const stack = [{ node: this.root, depth: 0 }];
      const result = [];
      while (stack.length > 0) {
        const { node: node2, depth } = stack.pop();
        if (node2 == null)
          continue;
        if (node2.left != null) {
          stack.push({ node: node2.left, depth: depth + 1 });
        }
        if (node2.right != null) {
          stack.push({ node: node2.right, depth: depth + 1 });
        }
        const isInsidePolygon = _BKDTree.isPointInPolygon(polygon, node2.point);
        if (isInsidePolygon && inclusive || !isInsidePolygon && !inclusive) {
          result.push({ point: node2.point, docIDs: Array.from(node2.docIDs) });
        }
      }
      const centroid = _BKDTree.calculatePolygonCentroid(polygon);
      if (sort) {
        const distanceFn = highPrecision ? _BKDTree.vincentyDistance : _BKDTree.haversineDistance;
        result.sort((a, b) => {
          const distA = distanceFn(centroid, a.point);
          const distB = distanceFn(centroid, b.point);
          return sort.toLowerCase() === "asc" ? distA - distB : distB - distA;
        });
      }
      return result;
    }
    toJSON() {
      return {
        root: this.root ? this.root.toJSON() : null
      };
    }
    static fromJSON(json) {
      const tree = new _BKDTree();
      if (json.root) {
        tree.root = BKDNode.fromJSON(json.root);
        tree.buildNodeMap(tree.root);
      }
      return tree;
    }
    buildNodeMap(node2) {
      if (node2 == null)
        return;
      const pointKey = this.getPointKey(node2.point);
      this.nodeMap.set(pointKey, node2);
      if (node2.left) {
        this.buildNodeMap(node2.left);
      }
      if (node2.right) {
        this.buildNodeMap(node2.right);
      }
    }
    static calculatePolygonCentroid(polygon) {
      let totalArea = 0;
      let centroidX = 0;
      let centroidY = 0;
      const polygonLength = polygon.length;
      for (let i = 0, j = polygonLength - 1; i < polygonLength; j = i++) {
        const xi = polygon[i].lon;
        const yi = polygon[i].lat;
        const xj = polygon[j].lon;
        const yj = polygon[j].lat;
        const areaSegment = xi * yj - xj * yi;
        totalArea += areaSegment;
        centroidX += (xi + xj) * areaSegment;
        centroidY += (yi + yj) * areaSegment;
      }
      totalArea /= 2;
      const centroidCoordinate = 6 * totalArea;
      centroidX /= centroidCoordinate;
      centroidY /= centroidCoordinate;
      return { lon: centroidX, lat: centroidY };
    }
    static isPointInPolygon(polygon, point3) {
      let isInside = false;
      const x = point3.lon;
      const y = point3.lat;
      const polygonLength = polygon.length;
      for (let i = 0, j = polygonLength - 1; i < polygonLength; j = i++) {
        const xi = polygon[i].lon;
        const yi = polygon[i].lat;
        const xj = polygon[j].lon;
        const yj = polygon[j].lat;
        const intersect2 = yi > y !== yj > y && x < (xj - xi) * (y - yi) / (yj - yi) + xi;
        if (intersect2)
          isInside = !isInside;
      }
      return isInside;
    }
    static haversineDistance(coord1, coord2) {
      const P = Math.PI / 180;
      const lat1 = coord1.lat * P;
      const lat2 = coord2.lat * P;
      const deltaLat = (coord2.lat - coord1.lat) * P;
      const deltaLon = (coord2.lon - coord1.lon) * P;
      const a = Math.sin(deltaLat / 2) * Math.sin(deltaLat / 2) + Math.cos(lat1) * Math.cos(lat2) * Math.sin(deltaLon / 2) * Math.sin(deltaLon / 2);
      const c2 = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));
      return EARTH_RADIUS * c2;
    }
    static vincentyDistance(coord1, coord2) {
      const a = 6378137;
      const f = 1 / 298.257223563;
      const b = (1 - f) * a;
      const P = Math.PI / 180;
      const lat1 = coord1.lat * P;
      const lat2 = coord2.lat * P;
      const deltaLon = (coord2.lon - coord1.lon) * P;
      const U1 = Math.atan((1 - f) * Math.tan(lat1));
      const U2 = Math.atan((1 - f) * Math.tan(lat2));
      const sinU1 = Math.sin(U1);
      const cosU1 = Math.cos(U1);
      const sinU2 = Math.sin(U2);
      const cosU2 = Math.cos(U2);
      let lambda = deltaLon;
      let prevLambda;
      let iterationLimit = 1e3;
      let sinSigma;
      let cosSigma;
      let sigma;
      let sinAlpha;
      let cos2Alpha;
      let cos2SigmaM;
      do {
        const sinLambda = Math.sin(lambda);
        const cosLambda = Math.cos(lambda);
        sinSigma = Math.sqrt(cosU2 * sinLambda * (cosU2 * sinLambda) + (cosU1 * sinU2 - sinU1 * cosU2 * cosLambda) * (cosU1 * sinU2 - sinU1 * cosU2 * cosLambda));
        if (sinSigma === 0)
          return 0;
        cosSigma = sinU1 * sinU2 + cosU1 * cosU2 * cosLambda;
        sigma = Math.atan2(sinSigma, cosSigma);
        sinAlpha = cosU1 * cosU2 * sinLambda / sinSigma;
        cos2Alpha = 1 - sinAlpha * sinAlpha;
        cos2SigmaM = cosSigma - 2 * sinU1 * sinU2 / cos2Alpha;
        if (isNaN(cos2SigmaM))
          cos2SigmaM = 0;
        const C2 = f / 16 * cos2Alpha * (4 + f * (4 - 3 * cos2Alpha));
        prevLambda = lambda;
        lambda = deltaLon + (1 - C2) * f * sinAlpha * (sigma + C2 * sinSigma * (cos2SigmaM + C2 * cosSigma * (-1 + 2 * cos2SigmaM * cos2SigmaM)));
      } while (Math.abs(lambda - prevLambda) > 1e-12 && --iterationLimit > 0);
      if (iterationLimit === 0) {
        return NaN;
      }
      const uSquared = cos2Alpha * (a * a - b * b) / (b * b);
      const A = 1 + uSquared / 16384 * (4096 + uSquared * (-768 + uSquared * (320 - 175 * uSquared)));
      const B = uSquared / 1024 * (256 + uSquared * (-128 + uSquared * (74 - 47 * uSquared)));
      const deltaSigma = B * sinSigma * (cos2SigmaM + B / 4 * (cosSigma * (-1 + 2 * cos2SigmaM * cos2SigmaM) - B / 6 * cos2SigmaM * (-3 + 4 * sinSigma * sinSigma) * (-3 + 4 * cos2SigmaM * cos2SigmaM)));
      const s = b * A * (sigma - deltaSigma);
      return s;
    }
  };

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/bool.js
  var BoolNode = class _BoolNode {
    constructor() {
      __publicField(this, "true");
      __publicField(this, "false");
      this.true = /* @__PURE__ */ new Set();
      this.false = /* @__PURE__ */ new Set();
    }
    insert(value, bool) {
      if (bool) {
        this.true.add(value);
      } else {
        this.false.add(value);
      }
    }
    delete(value, bool) {
      if (bool) {
        this.true.delete(value);
      } else {
        this.false.delete(value);
      }
    }
    getSize() {
      return this.true.size + this.false.size;
    }
    toJSON() {
      return {
        true: Array.from(this.true),
        false: Array.from(this.false)
      };
    }
    static fromJSON(json) {
      const node2 = new _BoolNode();
      node2.true = new Set(json.true);
      node2.false = new Set(json.false);
      return node2;
    }
  };

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/algorithms.js
  function BM25(tf, matchingCount, docsCount, fieldLength, averageFieldLength, { k, b, d }) {
    const idf = Math.log(1 + (docsCount - matchingCount + 0.5) / (matchingCount + 0.5));
    return idf * (d + tf * (k + 1)) / (tf + k * (1 - b + b * fieldLength / averageFieldLength));
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/trees/vector.js
  var DEFAULT_SIMILARITY = 0.8;
  var VectorIndex = class _VectorIndex {
    constructor(size) {
      __publicField(this, "size");
      __publicField(this, "vectors", /* @__PURE__ */ new Map());
      this.size = size;
    }
    add(internalDocumentId, value) {
      if (!(value instanceof Float32Array)) {
        value = new Float32Array(value);
      }
      const magnitude = getMagnitude(value, this.size);
      this.vectors.set(internalDocumentId, [magnitude, value]);
    }
    remove(internalDocumentId) {
      this.vectors.delete(internalDocumentId);
    }
    find(vector, similarity, whereFiltersIDs) {
      if (!(vector instanceof Float32Array)) {
        vector = new Float32Array(vector);
      }
      const results = findSimilarVectors(vector, whereFiltersIDs, this.vectors, this.size, similarity);
      return results;
    }
    toJSON() {
      const vectors = [];
      for (const [id, [magnitude, vector]] of this.vectors) {
        vectors.push([id, [magnitude, Array.from(vector)]]);
      }
      return {
        size: this.size,
        vectors
      };
    }
    static fromJSON(json) {
      const raw = json;
      const index2 = new _VectorIndex(raw.size);
      for (const [id, [magnitude, vector]] of raw.vectors) {
        index2.vectors.set(id, [magnitude, new Float32Array(vector)]);
      }
      return index2;
    }
  };
  function getMagnitude(vector, vectorLength) {
    let magnitude = 0;
    for (let i = 0; i < vectorLength; i++) {
      magnitude += vector[i] * vector[i];
    }
    return Math.sqrt(magnitude);
  }
  function findSimilarVectors(targetVector, keys, vectors, length, threshold) {
    const targetMagnitude = getMagnitude(targetVector, length);
    const similarVectors = [];
    const base = keys ? keys : vectors.keys();
    for (const vectorId of base) {
      const entry = vectors.get(vectorId);
      if (!entry) {
        continue;
      }
      const magnitude = entry[0];
      const vector = entry[1];
      let dotProduct = 0;
      for (let i = 0; i < length; i++) {
        dotProduct += targetVector[i] * vector[i];
      }
      const similarity = dotProduct / (targetMagnitude * magnitude);
      if (similarity >= threshold) {
        similarVectors.push([vectorId, similarity]);
      }
    }
    return similarVectors;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/index.js
  function insertDocumentScoreParameters(index2, prop, id, tokens, docsCount) {
    const internalId = getInternalDocumentId(index2.sharedInternalDocumentStore, id);
    index2.avgFieldLength[prop] = ((index2.avgFieldLength[prop] ?? 0) * (docsCount - 1) + tokens.length) / docsCount;
    index2.fieldLengths[prop][internalId] = tokens.length;
    index2.frequencies[prop][internalId] = {};
  }
  function insertTokenScoreParameters(index2, prop, id, tokens, token) {
    let tokenFrequency = 0;
    for (const t of tokens) {
      if (t === token) {
        tokenFrequency++;
      }
    }
    const internalId = getInternalDocumentId(index2.sharedInternalDocumentStore, id);
    const tf = tokenFrequency / tokens.length;
    index2.frequencies[prop][internalId][token] = tf;
    if (!(token in index2.tokenOccurrences[prop])) {
      index2.tokenOccurrences[prop][token] = 0;
    }
    index2.tokenOccurrences[prop][token] = (index2.tokenOccurrences[prop][token] ?? 0) + 1;
  }
  function removeDocumentScoreParameters(index2, prop, id, docsCount) {
    const internalId = getInternalDocumentId(index2.sharedInternalDocumentStore, id);
    if (docsCount > 1) {
      index2.avgFieldLength[prop] = (index2.avgFieldLength[prop] * docsCount - index2.fieldLengths[prop][internalId]) / (docsCount - 1);
    } else {
      index2.avgFieldLength[prop] = void 0;
    }
    index2.fieldLengths[prop][internalId] = void 0;
    index2.frequencies[prop][internalId] = void 0;
  }
  function removeTokenScoreParameters(index2, prop, token) {
    index2.tokenOccurrences[prop][token]--;
  }
  function create2(orama, sharedInternalDocumentStore, schema, index2, prefix = "") {
    if (!index2) {
      index2 = {
        sharedInternalDocumentStore,
        indexes: {},
        vectorIndexes: {},
        searchableProperties: [],
        searchablePropertiesWithTypes: {},
        frequencies: {},
        tokenOccurrences: {},
        avgFieldLength: {},
        fieldLengths: {}
      };
    }
    for (const [prop, type] of Object.entries(schema)) {
      const path2 = `${prefix}${prefix ? "." : ""}${prop}`;
      if (typeof type === "object" && !Array.isArray(type)) {
        create2(orama, sharedInternalDocumentStore, type, index2, path2);
        continue;
      }
      if (isVectorType(type)) {
        index2.searchableProperties.push(path2);
        index2.searchablePropertiesWithTypes[path2] = type;
        index2.vectorIndexes[path2] = {
          type: "Vector",
          node: new VectorIndex(getVectorSize(type)),
          isArray: false
        };
      } else {
        const isArray2 = /\[/.test(type);
        switch (type) {
          case "boolean":
          case "boolean[]":
            index2.indexes[path2] = { type: "Bool", node: new BoolNode(), isArray: isArray2 };
            break;
          case "number":
          case "number[]":
            index2.indexes[path2] = { type: "AVL", node: new AVLTree(0, []), isArray: isArray2 };
            break;
          case "string":
          case "string[]":
            index2.indexes[path2] = { type: "Radix", node: new RadixTree(), isArray: isArray2 };
            index2.avgFieldLength[path2] = 0;
            index2.frequencies[path2] = {};
            index2.tokenOccurrences[path2] = {};
            index2.fieldLengths[path2] = {};
            break;
          case "enum":
          case "enum[]":
            index2.indexes[path2] = { type: "Flat", node: new FlatTree(), isArray: isArray2 };
            break;
          case "geopoint":
            index2.indexes[path2] = { type: "BKD", node: new BKDTree(), isArray: isArray2 };
            break;
          default:
            throw createError("INVALID_SCHEMA_TYPE", Array.isArray(type) ? "array" : type, path2);
        }
        index2.searchableProperties.push(path2);
        index2.searchablePropertiesWithTypes[path2] = type;
      }
    }
    return index2;
  }
  function insertScalarBuilder(implementation, index2, prop, internalId, language, tokenizer3, docsCount, options) {
    return (value) => {
      const { type, node: node2 } = index2.indexes[prop];
      switch (type) {
        case "Bool": {
          node2[value ? "true" : "false"].add(internalId);
          break;
        }
        case "AVL": {
          const avlRebalanceThreshold = options?.avlRebalanceThreshold ?? 1;
          node2.insert(value, internalId, avlRebalanceThreshold);
          break;
        }
        case "Radix": {
          const tokens = tokenizer3.tokenize(value, language, prop, false);
          implementation.insertDocumentScoreParameters(index2, prop, internalId, tokens, docsCount);
          for (const token of tokens) {
            implementation.insertTokenScoreParameters(index2, prop, internalId, tokens, token);
            node2.insert(token, internalId);
          }
          break;
        }
        case "Flat": {
          node2.insert(value, internalId);
          break;
        }
        case "BKD": {
          node2.insert(value, [internalId]);
          break;
        }
      }
    };
  }
  function insert(implementation, index2, prop, id, internalId, value, schemaType, language, tokenizer3, docsCount, options) {
    if (isVectorType(schemaType)) {
      return insertVector(index2, prop, value, id, internalId);
    }
    const insertScalar = insertScalarBuilder(implementation, index2, prop, internalId, language, tokenizer3, docsCount, options);
    if (!isArrayType(schemaType)) {
      return insertScalar(value);
    }
    const elements = value;
    const elementsLength = elements.length;
    for (let i = 0; i < elementsLength; i++) {
      insertScalar(elements[i]);
    }
  }
  function insertVector(index2, prop, value, id, internalDocumentId) {
    index2.vectorIndexes[prop].node.add(internalDocumentId, value);
  }
  function removeScalar(implementation, index2, prop, id, internalId, value, schemaType, language, tokenizer3, docsCount) {
    if (isVectorType(schemaType)) {
      index2.vectorIndexes[prop].node.remove(internalId);
      return true;
    }
    const { type, node: node2 } = index2.indexes[prop];
    switch (type) {
      case "AVL": {
        node2.removeDocument(value, internalId);
        return true;
      }
      case "Bool": {
        node2[value ? "true" : "false"].delete(internalId);
        return true;
      }
      case "Radix": {
        const tokens = tokenizer3.tokenize(value, language, prop);
        implementation.removeDocumentScoreParameters(index2, prop, id, docsCount);
        for (const token of tokens) {
          implementation.removeTokenScoreParameters(index2, prop, token);
          node2.removeDocumentByWord(token, internalId);
        }
        return true;
      }
      case "Flat": {
        node2.removeDocument(internalId, value);
        return true;
      }
      case "BKD": {
        node2.removeDocByID(value, internalId);
        return false;
      }
    }
  }
  function remove2(implementation, index2, prop, id, internalId, value, schemaType, language, tokenizer3, docsCount) {
    if (!isArrayType(schemaType)) {
      return removeScalar(implementation, index2, prop, id, internalId, value, schemaType, language, tokenizer3, docsCount);
    }
    const innerSchemaType = getInnerType(schemaType);
    const elements = value;
    const elementsLength = elements.length;
    for (let i = 0; i < elementsLength; i++) {
      removeScalar(implementation, index2, prop, id, internalId, elements[i], innerSchemaType, language, tokenizer3, docsCount);
    }
    return true;
  }
  function calculateResultScores(index2, prop, term, ids, docsCount, bm25Relevance, resultsMap, boostPerProperty, whereFiltersIDs, keywordMatchesMap) {
    const documentIDs = Array.from(ids);
    const avgFieldLength = index2.avgFieldLength[prop];
    const fieldLengths = index2.fieldLengths[prop];
    const oramaOccurrences = index2.tokenOccurrences[prop];
    const oramaFrequencies = index2.frequencies[prop];
    const termOccurrences = typeof oramaOccurrences[term] === "number" ? oramaOccurrences[term] ?? 0 : 0;
    const documentIDsLength = documentIDs.length;
    for (let k = 0; k < documentIDsLength; k++) {
      const internalId = documentIDs[k];
      if (whereFiltersIDs && !whereFiltersIDs.has(internalId)) {
        continue;
      }
      if (!keywordMatchesMap.has(internalId)) {
        keywordMatchesMap.set(internalId, /* @__PURE__ */ new Map());
      }
      const propertyMatches = keywordMatchesMap.get(internalId);
      propertyMatches.set(prop, (propertyMatches.get(prop) || 0) + 1);
      const tf = oramaFrequencies?.[internalId]?.[term] ?? 0;
      const bm25 = BM25(tf, termOccurrences, docsCount, fieldLengths[internalId], avgFieldLength, bm25Relevance);
      if (resultsMap.has(internalId)) {
        resultsMap.set(internalId, resultsMap.get(internalId) + bm25 * boostPerProperty);
      } else {
        resultsMap.set(internalId, bm25 * boostPerProperty);
      }
    }
  }
  function search2(index2, term, tokenizer3, language, propertiesToSearch, exact, tolerance, boost, relevance, docsCount, whereFiltersIDs, threshold = 0) {
    const tokens = tokenizer3.tokenize(term, language);
    const keywordsCount = tokens.length || 1;
    const keywordMatchesMap = /* @__PURE__ */ new Map();
    const tokenFoundMap = /* @__PURE__ */ new Map();
    const resultsMap = /* @__PURE__ */ new Map();
    for (const prop of propertiesToSearch) {
      if (!(prop in index2.indexes)) {
        continue;
      }
      const tree = index2.indexes[prop];
      const { type } = tree;
      if (type !== "Radix") {
        throw createError("WRONG_SEARCH_PROPERTY_TYPE", prop);
      }
      const boostPerProperty = boost[prop] ?? 1;
      if (boostPerProperty <= 0) {
        throw createError("INVALID_BOOST_VALUE", boostPerProperty);
      }
      if (tokens.length === 0 && !term) {
        tokens.push("");
      }
      const tokenLength = tokens.length;
      for (let i = 0; i < tokenLength; i++) {
        const token = tokens[i];
        const searchResult = tree.node.find({ term: token, exact, tolerance });
        const termsFound = Object.keys(searchResult);
        if (termsFound.length > 0) {
          tokenFoundMap.set(token, true);
        }
        const termsFoundLength = termsFound.length;
        for (let j = 0; j < termsFoundLength; j++) {
          const word = termsFound[j];
          const ids = searchResult[word];
          calculateResultScores(index2, prop, word, ids, docsCount, relevance, resultsMap, boostPerProperty, whereFiltersIDs, keywordMatchesMap);
        }
      }
    }
    const results = Array.from(resultsMap.entries()).map(([id, score]) => [id, score]).sort((a, b) => b[1] - a[1]);
    if (results.length === 0) {
      return [];
    }
    if (threshold === 1) {
      return results;
    }
    if (threshold === 0) {
      if (keywordsCount === 1) {
        return results;
      }
      for (const token of tokens) {
        if (!tokenFoundMap.get(token)) {
          return [];
        }
      }
      const fullMatches2 = results.filter(([id]) => {
        const propertyMatches = keywordMatchesMap.get(id);
        if (!propertyMatches)
          return false;
        return Array.from(propertyMatches.values()).some((matches) => matches === keywordsCount);
      });
      return fullMatches2;
    }
    const fullMatches = results.filter(([id]) => {
      const propertyMatches = keywordMatchesMap.get(id);
      if (!propertyMatches)
        return false;
      return Array.from(propertyMatches.values()).some((matches) => matches === keywordsCount);
    });
    if (fullMatches.length > 0) {
      const remainingResults = results.filter(([id]) => !fullMatches.some(([fid]) => fid === id));
      const additionalResults = Math.ceil(remainingResults.length * threshold);
      return [...fullMatches, ...remainingResults.slice(0, additionalResults)];
    }
    return results;
  }
  function searchByWhereClause(index2, tokenizer3, filters, language) {
    if ("and" in filters && filters.and && Array.isArray(filters.and)) {
      const andFilters = filters.and;
      if (andFilters.length === 0) {
        return /* @__PURE__ */ new Set();
      }
      const results = andFilters.map((filter) => searchByWhereClause(index2, tokenizer3, filter, language));
      return setIntersection(...results);
    }
    if ("or" in filters && filters.or && Array.isArray(filters.or)) {
      const orFilters = filters.or;
      if (orFilters.length === 0) {
        return /* @__PURE__ */ new Set();
      }
      const results = orFilters.map((filter) => searchByWhereClause(index2, tokenizer3, filter, language));
      return results.reduce((acc, set) => setUnion(acc, set), /* @__PURE__ */ new Set());
    }
    if ("not" in filters && filters.not) {
      const notFilter = filters.not;
      const allDocs = /* @__PURE__ */ new Set();
      const docsStore = index2.sharedInternalDocumentStore;
      for (let i = 1; i <= docsStore.internalIdToId.length; i++) {
        allDocs.add(i);
      }
      const notResult = searchByWhereClause(index2, tokenizer3, notFilter, language);
      return setDifference(allDocs, notResult);
    }
    const filterKeys = Object.keys(filters);
    const filtersMap = filterKeys.reduce((acc, key) => ({
      [key]: /* @__PURE__ */ new Set(),
      ...acc
    }), {});
    for (const param of filterKeys) {
      const operation = filters[param];
      if (typeof index2.indexes[param] === "undefined") {
        throw createError("UNKNOWN_FILTER_PROPERTY", param);
      }
      const { node: node2, type, isArray: isArray2 } = index2.indexes[param];
      if (type === "Bool") {
        const idx = node2;
        const filteredIDs = operation ? idx.true : idx.false;
        filtersMap[param] = setUnion(filtersMap[param], filteredIDs);
        continue;
      }
      if (type === "BKD") {
        let reqOperation;
        if ("radius" in operation) {
          reqOperation = "radius";
        } else if ("polygon" in operation) {
          reqOperation = "polygon";
        } else {
          throw new Error(`Invalid operation ${operation}`);
        }
        if (reqOperation === "radius") {
          const { value, coordinates, unit = "m", inside = true, highPrecision = false } = operation[reqOperation];
          const distanceInMeters = convertDistanceToMeters(value, unit);
          const ids = node2.searchByRadius(coordinates, distanceInMeters, inside, void 0, highPrecision);
          filtersMap[param] = addGeoResult(filtersMap[param], ids);
        } else {
          const { coordinates, inside = true, highPrecision = false } = operation[reqOperation];
          const ids = node2.searchByPolygon(coordinates, inside, void 0, highPrecision);
          filtersMap[param] = addGeoResult(filtersMap[param], ids);
        }
        continue;
      }
      if (type === "Radix" && (typeof operation === "string" || Array.isArray(operation))) {
        for (const raw of [operation].flat()) {
          const term = tokenizer3.tokenize(raw, language, param);
          for (const t of term) {
            const filteredIDsResults = node2.find({ term: t, exact: true });
            filtersMap[param] = addFindResult(filtersMap[param], filteredIDsResults);
          }
        }
        continue;
      }
      const operationKeys = Object.keys(operation);
      if (operationKeys.length > 1) {
        throw createError("INVALID_FILTER_OPERATION", operationKeys.length);
      }
      if (type === "Flat") {
        const results = new Set(isArray2 ? node2.filterArr(operation) : node2.filter(operation));
        filtersMap[param] = setUnion(filtersMap[param], results);
        continue;
      }
      if (type === "AVL") {
        const operationOpt = operationKeys[0];
        const operationValue = operation[operationOpt];
        let filteredIDs;
        switch (operationOpt) {
          case "gt": {
            filteredIDs = node2.greaterThan(operationValue, false);
            break;
          }
          case "gte": {
            filteredIDs = node2.greaterThan(operationValue, true);
            break;
          }
          case "lt": {
            filteredIDs = node2.lessThan(operationValue, false);
            break;
          }
          case "lte": {
            filteredIDs = node2.lessThan(operationValue, true);
            break;
          }
          case "eq": {
            const ret = node2.find(operationValue);
            filteredIDs = ret ?? /* @__PURE__ */ new Set();
            break;
          }
          case "between": {
            const [min, max] = operationValue;
            filteredIDs = node2.rangeSearch(min, max);
            break;
          }
          default:
            throw createError("INVALID_FILTER_OPERATION", operationOpt);
        }
        filtersMap[param] = setUnion(filtersMap[param], filteredIDs);
      }
    }
    return setIntersection(...Object.values(filtersMap));
  }
  function getSearchableProperties(index2) {
    return index2.searchableProperties;
  }
  function getSearchablePropertiesWithTypes(index2) {
    return index2.searchablePropertiesWithTypes;
  }
  function load3(sharedInternalDocumentStore, raw) {
    const { indexes: rawIndexes, vectorIndexes: rawVectorIndexes, searchableProperties, searchablePropertiesWithTypes, frequencies, tokenOccurrences, avgFieldLength, fieldLengths } = raw;
    const indexes = {};
    const vectorIndexes = {};
    for (const prop of Object.keys(rawIndexes)) {
      const { node: node2, type, isArray: isArray2 } = rawIndexes[prop];
      switch (type) {
        case "Radix":
          indexes[prop] = {
            type: "Radix",
            node: RadixTree.fromJSON(node2),
            isArray: isArray2
          };
          break;
        case "Flat":
          indexes[prop] = {
            type: "Flat",
            node: FlatTree.fromJSON(node2),
            isArray: isArray2
          };
          break;
        case "AVL":
          indexes[prop] = {
            type: "AVL",
            node: AVLTree.fromJSON(node2),
            isArray: isArray2
          };
          break;
        case "BKD":
          indexes[prop] = {
            type: "BKD",
            node: BKDTree.fromJSON(node2),
            isArray: isArray2
          };
          break;
        case "Bool":
          indexes[prop] = {
            type: "Bool",
            node: BoolNode.fromJSON(node2),
            isArray: isArray2
          };
          break;
        default:
          indexes[prop] = rawIndexes[prop];
      }
    }
    for (const idx of Object.keys(rawVectorIndexes)) {
      vectorIndexes[idx] = {
        type: "Vector",
        isArray: false,
        node: VectorIndex.fromJSON(rawVectorIndexes[idx])
      };
    }
    return {
      sharedInternalDocumentStore,
      indexes,
      vectorIndexes,
      searchableProperties,
      searchablePropertiesWithTypes,
      frequencies,
      tokenOccurrences,
      avgFieldLength,
      fieldLengths
    };
  }
  function save3(index2) {
    const { indexes, vectorIndexes, searchableProperties, searchablePropertiesWithTypes, frequencies, tokenOccurrences, avgFieldLength, fieldLengths } = index2;
    const dumpVectorIndexes = {};
    for (const idx of Object.keys(vectorIndexes)) {
      dumpVectorIndexes[idx] = vectorIndexes[idx].node.toJSON();
    }
    const savedIndexes = {};
    for (const name2 of Object.keys(indexes)) {
      const { type, node: node2, isArray: isArray2 } = indexes[name2];
      if (type === "Flat" || type === "Radix" || type === "AVL" || type === "BKD" || type === "Bool") {
        savedIndexes[name2] = {
          type,
          node: node2.toJSON(),
          isArray: isArray2
        };
      } else {
        savedIndexes[name2] = indexes[name2];
        savedIndexes[name2].node = savedIndexes[name2].node.toJSON();
      }
    }
    return {
      indexes: savedIndexes,
      vectorIndexes: dumpVectorIndexes,
      searchableProperties,
      searchablePropertiesWithTypes,
      frequencies,
      tokenOccurrences,
      avgFieldLength,
      fieldLengths
    };
  }
  function createIndex() {
    return {
      create: create2,
      insert,
      remove: remove2,
      insertDocumentScoreParameters,
      insertTokenScoreParameters,
      removeDocumentScoreParameters,
      removeTokenScoreParameters,
      calculateResultScores,
      search: search2,
      searchByWhereClause,
      getSearchableProperties,
      getSearchablePropertiesWithTypes,
      load: load3,
      save: save3
    };
  }
  function addGeoResult(set, ids) {
    if (!set) {
      set = /* @__PURE__ */ new Set();
    }
    const idsLength = ids.length;
    for (let i = 0; i < idsLength; i++) {
      const entry = ids[i].docIDs;
      const idsLength2 = entry.length;
      for (let j = 0; j < idsLength2; j++) {
        set.add(entry[j]);
      }
    }
    return set;
  }
  function createGeoTokenScores(ids, centerPoint, highPrecision = false) {
    const distanceFn = highPrecision ? BKDTree.vincentyDistance : BKDTree.haversineDistance;
    const results = [];
    const distances = [];
    for (const { point: point3 } of ids) {
      distances.push(distanceFn(centerPoint, point3));
    }
    const maxDistance = Math.max(...distances);
    let index2 = 0;
    for (const { docIDs } of ids) {
      const distance = distances[index2];
      const score = maxDistance - distance + 1;
      for (const docID of docIDs) {
        results.push([docID, score]);
      }
      index2++;
    }
    results.sort((a, b) => b[1] - a[1]);
    return results;
  }
  function isGeosearchOnlyQuery(filters, index2) {
    const filterKeys = Object.keys(filters);
    if (filterKeys.length !== 1) {
      return { isGeoOnly: false };
    }
    const param = filterKeys[0];
    const operation = filters[param];
    if (typeof index2.indexes[param] === "undefined") {
      return { isGeoOnly: false };
    }
    const { type } = index2.indexes[param];
    if (type === "BKD" && operation && ("radius" in operation || "polygon" in operation)) {
      return { isGeoOnly: true, geoProperty: param, geoOperation: operation };
    }
    return { isGeoOnly: false };
  }
  function searchByGeoWhereClause(index2, filters) {
    const indexTyped = index2;
    const geoInfo = isGeosearchOnlyQuery(filters, indexTyped);
    if (!geoInfo.isGeoOnly || !geoInfo.geoProperty || !geoInfo.geoOperation) {
      return null;
    }
    const { node: node2 } = indexTyped.indexes[geoInfo.geoProperty];
    const operation = geoInfo.geoOperation;
    const bkdNode = node2;
    let results;
    if ("radius" in operation) {
      const { value, coordinates, unit = "m", inside = true, highPrecision = false } = operation.radius;
      const centerPoint = coordinates;
      const distanceInMeters = convertDistanceToMeters(value, unit);
      results = bkdNode.searchByRadius(centerPoint, distanceInMeters, inside, "asc", highPrecision);
      return createGeoTokenScores(results, centerPoint, highPrecision);
    } else if ("polygon" in operation) {
      const { coordinates, inside = true, highPrecision = false } = operation.polygon;
      results = bkdNode.searchByPolygon(coordinates, inside, "asc", highPrecision);
      const centroid = BKDTree.calculatePolygonCentroid(coordinates);
      return createGeoTokenScores(results, centroid, highPrecision);
    }
    return null;
  }
  function addFindResult(set, filteredIDsResults) {
    if (!set) {
      set = /* @__PURE__ */ new Set();
    }
    const keys = Object.keys(filteredIDsResults);
    const keysLength = keys.length;
    for (let i = 0; i < keysLength; i++) {
      const ids = filteredIDsResults[keys[i]];
      const idsLength = ids.length;
      for (let j = 0; j < idsLength; j++) {
        set.add(ids[j]);
      }
    }
    return set;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/sorter.js
  function innerCreate(orama, sharedInternalDocumentStore, schema, sortableDeniedProperties, prefix) {
    const sorter = {
      language: orama.tokenizer.language,
      sharedInternalDocumentStore,
      enabled: true,
      isSorted: true,
      sortableProperties: [],
      sortablePropertiesWithTypes: {},
      sorts: {}
    };
    for (const [prop, type] of Object.entries(schema)) {
      const path2 = `${prefix}${prefix ? "." : ""}${prop}`;
      if (sortableDeniedProperties.includes(path2)) {
        continue;
      }
      if (typeof type === "object" && !Array.isArray(type)) {
        const ret = innerCreate(orama, sharedInternalDocumentStore, type, sortableDeniedProperties, path2);
        safeArrayPush(sorter.sortableProperties, ret.sortableProperties);
        sorter.sorts = {
          ...sorter.sorts,
          ...ret.sorts
        };
        sorter.sortablePropertiesWithTypes = {
          ...sorter.sortablePropertiesWithTypes,
          ...ret.sortablePropertiesWithTypes
        };
        continue;
      }
      if (!isVectorType(type)) {
        switch (type) {
          case "boolean":
          case "number":
          case "string":
            sorter.sortableProperties.push(path2);
            sorter.sortablePropertiesWithTypes[path2] = type;
            sorter.sorts[path2] = {
              docs: /* @__PURE__ */ new Map(),
              orderedDocsToRemove: /* @__PURE__ */ new Map(),
              orderedDocs: [],
              type
            };
            break;
          case "geopoint":
          case "enum":
            continue;
          case "enum[]":
          case "boolean[]":
          case "number[]":
          case "string[]":
            continue;
          default:
            throw createError("INVALID_SORT_SCHEMA_TYPE", Array.isArray(type) ? "array" : type, path2);
        }
      }
    }
    return sorter;
  }
  function create3(orama, sharedInternalDocumentStore, schema, config) {
    const isSortEnabled = config?.enabled !== false;
    if (!isSortEnabled) {
      return {
        disabled: true
      };
    }
    return innerCreate(orama, sharedInternalDocumentStore, schema, (config || {}).unsortableProperties || [], "");
  }
  function insert2(sorter, prop, id, value) {
    if (!sorter.enabled) {
      return;
    }
    sorter.isSorted = false;
    const internalId = getInternalDocumentId(sorter.sharedInternalDocumentStore, id);
    const s = sorter.sorts[prop];
    if (s.orderedDocsToRemove.has(internalId)) {
      ensureOrderedDocsAreDeletedByProperty(sorter, prop);
    }
    s.docs.set(internalId, s.orderedDocs.length);
    s.orderedDocs.push([internalId, value]);
  }
  function ensureIsSorted(sorter) {
    if (sorter.isSorted || !sorter.enabled) {
      return;
    }
    const properties = Object.keys(sorter.sorts);
    for (const prop of properties) {
      ensurePropertyIsSorted(sorter, prop);
    }
    sorter.isSorted = true;
  }
  function stringSort(language, value, d) {
    return value[1].localeCompare(d[1], getLocale(language));
  }
  function numberSort(value, d) {
    return value[1] - d[1];
  }
  function booleanSort(value, d) {
    return d[1] ? -1 : 1;
  }
  function ensurePropertyIsSorted(sorter, prop) {
    const s = sorter.sorts[prop];
    let predicate;
    switch (s.type) {
      case "string":
        predicate = stringSort.bind(null, sorter.language);
        break;
      case "number":
        predicate = numberSort.bind(null);
        break;
      case "boolean":
        predicate = booleanSort.bind(null);
        break;
    }
    s.orderedDocs.sort(predicate);
    const orderedDocsLength = s.orderedDocs.length;
    for (let i = 0; i < orderedDocsLength; i++) {
      const docId = s.orderedDocs[i][0];
      s.docs.set(docId, i);
    }
  }
  function ensureOrderedDocsAreDeleted(sorter) {
    const properties = Object.keys(sorter.sorts);
    for (const prop of properties) {
      ensureOrderedDocsAreDeletedByProperty(sorter, prop);
    }
  }
  function ensureOrderedDocsAreDeletedByProperty(sorter, prop) {
    const s = sorter.sorts[prop];
    if (!s.orderedDocsToRemove.size)
      return;
    s.orderedDocs = s.orderedDocs.filter((doc) => !s.orderedDocsToRemove.has(doc[0]));
    s.orderedDocsToRemove.clear();
  }
  function remove3(sorter, prop, id) {
    if (!sorter.enabled) {
      return;
    }
    const s = sorter.sorts[prop];
    const internalId = getInternalDocumentId(sorter.sharedInternalDocumentStore, id);
    const index2 = s.docs.get(internalId);
    if (!index2)
      return;
    s.docs.delete(internalId);
    s.orderedDocsToRemove.set(internalId, true);
  }
  function sortBy(sorter, docIds, by) {
    if (!sorter.enabled) {
      throw createError("SORT_DISABLED");
    }
    const property2 = by.property;
    const isDesc = by.order === "DESC";
    const s = sorter.sorts[property2];
    if (!s) {
      throw createError("UNABLE_TO_SORT_ON_UNKNOWN_FIELD", property2, sorter.sortableProperties.join(", "));
    }
    ensureOrderedDocsAreDeletedByProperty(sorter, property2);
    ensureIsSorted(sorter);
    docIds.sort((a, b) => {
      const indexOfA = s.docs.get(getInternalDocumentId(sorter.sharedInternalDocumentStore, a[0]));
      const indexOfB = s.docs.get(getInternalDocumentId(sorter.sharedInternalDocumentStore, b[0]));
      const isAIndexed = typeof indexOfA !== "undefined";
      const isBIndexed = typeof indexOfB !== "undefined";
      if (!isAIndexed && !isBIndexed) {
        return 0;
      }
      if (!isAIndexed) {
        return 1;
      }
      if (!isBIndexed) {
        return -1;
      }
      return isDesc ? indexOfB - indexOfA : indexOfA - indexOfB;
    });
    return docIds;
  }
  function getSortableProperties(sorter) {
    if (!sorter.enabled) {
      return [];
    }
    return sorter.sortableProperties;
  }
  function getSortablePropertiesWithTypes(sorter) {
    if (!sorter.enabled) {
      return {};
    }
    return sorter.sortablePropertiesWithTypes;
  }
  function load4(sharedInternalDocumentStore, raw) {
    const rawDocument = raw;
    if (!rawDocument.enabled) {
      return {
        enabled: false
      };
    }
    const sorts = Object.keys(rawDocument.sorts).reduce((acc, prop) => {
      const { docs, orderedDocs, type } = rawDocument.sorts[prop];
      acc[prop] = {
        docs: new Map(Object.entries(docs).map(([k, v2]) => [+k, v2])),
        orderedDocsToRemove: /* @__PURE__ */ new Map(),
        orderedDocs,
        type
      };
      return acc;
    }, {});
    return {
      sharedInternalDocumentStore,
      language: rawDocument.language,
      sortableProperties: rawDocument.sortableProperties,
      sortablePropertiesWithTypes: rawDocument.sortablePropertiesWithTypes,
      sorts,
      enabled: true,
      isSorted: rawDocument.isSorted
    };
  }
  function save4(sorter) {
    if (!sorter.enabled) {
      return {
        enabled: false
      };
    }
    ensureOrderedDocsAreDeleted(sorter);
    ensureIsSorted(sorter);
    const sorts = Object.keys(sorter.sorts).reduce((acc, prop) => {
      const { docs, orderedDocs, type } = sorter.sorts[prop];
      acc[prop] = {
        docs: Object.fromEntries(docs.entries()),
        orderedDocs,
        type
      };
      return acc;
    }, {});
    return {
      language: sorter.language,
      sortableProperties: sorter.sortableProperties,
      sortablePropertiesWithTypes: sorter.sortablePropertiesWithTypes,
      sorts,
      enabled: sorter.enabled,
      isSorted: sorter.isSorted
    };
  }
  function createSorter() {
    return {
      create: create3,
      insert: insert2,
      remove: remove3,
      save: save4,
      load: load4,
      sortBy,
      getSortableProperties,
      getSortablePropertiesWithTypes
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/tokenizer/diacritics.js
  var DIACRITICS_CHARCODE_START = 192;
  var DIACRITICS_CHARCODE_END = 383;
  var CHARCODE_REPLACE_MAPPING = [
    65,
    65,
    65,
    65,
    65,
    65,
    65,
    67,
    69,
    69,
    69,
    69,
    73,
    73,
    73,
    73,
    69,
    78,
    79,
    79,
    79,
    79,
    79,
    null,
    79,
    85,
    85,
    85,
    85,
    89,
    80,
    115,
    97,
    97,
    97,
    97,
    97,
    97,
    97,
    99,
    101,
    101,
    101,
    101,
    105,
    105,
    105,
    105,
    101,
    110,
    111,
    111,
    111,
    111,
    111,
    null,
    111,
    117,
    117,
    117,
    117,
    121,
    112,
    121,
    65,
    97,
    65,
    97,
    65,
    97,
    67,
    99,
    67,
    99,
    67,
    99,
    67,
    99,
    68,
    100,
    68,
    100,
    69,
    101,
    69,
    101,
    69,
    101,
    69,
    101,
    69,
    101,
    71,
    103,
    71,
    103,
    71,
    103,
    71,
    103,
    72,
    104,
    72,
    104,
    73,
    105,
    73,
    105,
    73,
    105,
    73,
    105,
    73,
    105,
    73,
    105,
    74,
    106,
    75,
    107,
    107,
    76,
    108,
    76,
    108,
    76,
    108,
    76,
    108,
    76,
    108,
    78,
    110,
    78,
    110,
    78,
    110,
    110,
    78,
    110,
    79,
    111,
    79,
    111,
    79,
    111,
    79,
    111,
    82,
    114,
    82,
    114,
    82,
    114,
    83,
    115,
    83,
    115,
    83,
    115,
    83,
    115,
    84,
    116,
    84,
    116,
    84,
    116,
    85,
    117,
    85,
    117,
    85,
    117,
    85,
    117,
    85,
    117,
    85,
    117,
    87,
    119,
    89,
    121,
    89,
    90,
    122,
    90,
    122,
    90,
    122,
    115
  ];
  function replaceChar(charCode) {
    if (charCode < DIACRITICS_CHARCODE_START || charCode > DIACRITICS_CHARCODE_END)
      return charCode;
    return CHARCODE_REPLACE_MAPPING[charCode - DIACRITICS_CHARCODE_START] || charCode;
  }
  function replaceDiacritics(str) {
    const stringCharCode = [];
    for (let idx = 0; idx < str.length; idx++) {
      stringCharCode[idx] = replaceChar(str.charCodeAt(idx));
    }
    return String.fromCharCode(...stringCharCode);
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/tokenizer/english-stemmer.js
  var step2List = {
    ational: "ate",
    tional: "tion",
    enci: "ence",
    anci: "ance",
    izer: "ize",
    bli: "ble",
    alli: "al",
    entli: "ent",
    eli: "e",
    ousli: "ous",
    ization: "ize",
    ation: "ate",
    ator: "ate",
    alism: "al",
    iveness: "ive",
    fulness: "ful",
    ousness: "ous",
    aliti: "al",
    iviti: "ive",
    biliti: "ble",
    logi: "log"
  };
  var step3List = {
    icate: "ic",
    ative: "",
    alize: "al",
    iciti: "ic",
    ical: "ic",
    ful: "",
    ness: ""
  };
  var c = "[^aeiou]";
  var v = "[aeiouy]";
  var C = c + "[^aeiouy]*";
  var V = v + "[aeiou]*";
  var mgr0 = "^(" + C + ")?" + V + C;
  var meq1 = "^(" + C + ")?" + V + C + "(" + V + ")?$";
  var mgr1 = "^(" + C + ")?" + V + C + V + C;
  var s_v = "^(" + C + ")?" + v;
  function stemmer(w) {
    let stem;
    let suffix;
    let re;
    let re2;
    let re3;
    let re4;
    if (w.length < 3) {
      return w;
    }
    const firstch = w.substring(0, 1);
    if (firstch == "y") {
      w = firstch.toUpperCase() + w.substring(1);
    }
    re = /^(.+?)(ss|i)es$/;
    re2 = /^(.+?)([^s])s$/;
    if (re.test(w)) {
      w = w.replace(re, "$1$2");
    } else if (re2.test(w)) {
      w = w.replace(re2, "$1$2");
    }
    re = /^(.+?)eed$/;
    re2 = /^(.+?)(ed|ing)$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      re = new RegExp(mgr0);
      if (re.test(fp[1])) {
        re = /.$/;
        w = w.replace(re, "");
      }
    } else if (re2.test(w)) {
      const fp = re2.exec(w);
      stem = fp[1];
      re2 = new RegExp(s_v);
      if (re2.test(stem)) {
        w = stem;
        re2 = /(at|bl|iz)$/;
        re3 = new RegExp("([^aeiouylsz])\\1$");
        re4 = new RegExp("^" + C + v + "[^aeiouwxy]$");
        if (re2.test(w)) {
          w = w + "e";
        } else if (re3.test(w)) {
          re = /.$/;
          w = w.replace(re, "");
        } else if (re4.test(w)) {
          w = w + "e";
        }
      }
    }
    re = /^(.+?)y$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      stem = fp?.[1];
      re = new RegExp(s_v);
      if (stem && re.test(stem)) {
        w = stem + "i";
      }
    }
    re = /^(.+?)(ational|tional|enci|anci|izer|bli|alli|entli|eli|ousli|ization|ation|ator|alism|iveness|fulness|ousness|aliti|iviti|biliti|logi)$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      stem = fp?.[1];
      suffix = fp?.[2];
      re = new RegExp(mgr0);
      if (stem && re.test(stem)) {
        w = stem + step2List[suffix];
      }
    }
    re = /^(.+?)(icate|ative|alize|iciti|ical|ful|ness)$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      stem = fp?.[1];
      suffix = fp?.[2];
      re = new RegExp(mgr0);
      if (stem && re.test(stem)) {
        w = stem + step3List[suffix];
      }
    }
    re = /^(.+?)(al|ance|ence|er|ic|able|ible|ant|ement|ment|ent|ou|ism|ate|iti|ous|ive|ize)$/;
    re2 = /^(.+?)(s|t)(ion)$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      stem = fp?.[1];
      re = new RegExp(mgr1);
      if (stem && re.test(stem)) {
        w = stem;
      }
    } else if (re2.test(w)) {
      const fp = re2.exec(w);
      stem = fp?.[1] ?? "" + fp?.[2] ?? "";
      re2 = new RegExp(mgr1);
      if (re2.test(stem)) {
        w = stem;
      }
    }
    re = /^(.+?)e$/;
    if (re.test(w)) {
      const fp = re.exec(w);
      stem = fp?.[1];
      re = new RegExp(mgr1);
      re2 = new RegExp(meq1);
      re3 = new RegExp("^" + C + v + "[^aeiouwxy]$");
      if (stem && (re.test(stem) || re2.test(stem) && !re3.test(stem))) {
        w = stem;
      }
    }
    re = /ll$/;
    re2 = new RegExp(mgr1);
    if (re.test(w) && re2.test(w)) {
      re = /.$/;
      w = w.replace(re, "");
    }
    if (firstch == "y") {
      w = firstch.toLowerCase() + w.substring(1);
    }
    return w;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/tokenizer/index.js
  function normalizeToken(prop, token, withCache = true) {
    const key = `${this.language}:${prop}:${token}`;
    if (withCache && this.normalizationCache.has(key)) {
      return this.normalizationCache.get(key);
    }
    if (this.stopWords?.includes(token)) {
      if (withCache) {
        this.normalizationCache.set(key, "");
      }
      return "";
    }
    if (this.stemmer && !this.stemmerSkipProperties.has(prop)) {
      token = this.stemmer(token);
    }
    token = replaceDiacritics(token);
    if (withCache) {
      this.normalizationCache.set(key, token);
    }
    return token;
  }
  function trim(text5) {
    while (text5[text5.length - 1] === "") {
      text5.pop();
    }
    while (text5[0] === "") {
      text5.shift();
    }
    return text5;
  }
  function tokenize(input, language, prop, withCache = true) {
    if (language && language !== this.language) {
      throw createError("LANGUAGE_NOT_SUPPORTED", language);
    }
    if (typeof input !== "string") {
      return [input];
    }
    const normalizeToken2 = this.normalizeToken.bind(this, prop ?? "");
    let tokens;
    if (prop && this.tokenizeSkipProperties.has(prop)) {
      tokens = [normalizeToken2(input, withCache)];
    } else {
      const splitRule = SPLITTERS[this.language];
      tokens = input.toLowerCase().split(splitRule).map((t) => normalizeToken2(t, withCache)).filter(Boolean);
    }
    const trimTokens = trim(tokens);
    if (!this.allowDuplicates) {
      return Array.from(new Set(trimTokens));
    }
    return trimTokens;
  }
  function createTokenizer2(config = {}) {
    if (!config.language) {
      config.language = "english";
    } else if (!SUPPORTED_LANGUAGES.includes(config.language)) {
      throw createError("LANGUAGE_NOT_SUPPORTED", config.language);
    }
    let stemmer2;
    if (config.stemming || config.stemmer && !("stemming" in config)) {
      if (config.stemmer) {
        if (typeof config.stemmer !== "function") {
          throw createError("INVALID_STEMMER_FUNCTION_TYPE");
        }
        stemmer2 = config.stemmer;
      } else {
        if (config.language === "english") {
          stemmer2 = stemmer;
        } else {
          throw createError("MISSING_STEMMER", config.language);
        }
      }
    }
    let stopWords;
    if (config.stopWords !== false) {
      stopWords = [];
      if (Array.isArray(config.stopWords)) {
        stopWords = config.stopWords;
      } else if (typeof config.stopWords === "function") {
        stopWords = config.stopWords(stopWords);
      } else if (config.stopWords) {
        throw createError("CUSTOM_STOP_WORDS_MUST_BE_FUNCTION_OR_ARRAY");
      }
      if (!Array.isArray(stopWords)) {
        throw createError("CUSTOM_STOP_WORDS_MUST_BE_FUNCTION_OR_ARRAY");
      }
      for (const s of stopWords) {
        if (typeof s !== "string") {
          throw createError("CUSTOM_STOP_WORDS_MUST_BE_FUNCTION_OR_ARRAY");
        }
      }
    }
    const tokenizer3 = {
      tokenize,
      language: config.language,
      stemmer: stemmer2,
      stemmerSkipProperties: new Set(config.stemmerSkipProperties ? [config.stemmerSkipProperties].flat() : []),
      tokenizeSkipProperties: new Set(config.tokenizeSkipProperties ? [config.tokenizeSkipProperties].flat() : []),
      stopWords,
      allowDuplicates: Boolean(config.allowDuplicates),
      normalizeToken,
      normalizationCache: /* @__PURE__ */ new Map()
    };
    tokenizer3.tokenize = tokenize.bind(tokenizer3);
    tokenizer3.normalizeToken = normalizeToken;
    return tokenizer3;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/pinning.js
  function create4(sharedInternalDocumentStore) {
    return {
      sharedInternalDocumentStore,
      rules: /* @__PURE__ */ new Map()
    };
  }
  function addRule(store2, rule) {
    if (store2.rules.has(rule.id)) {
      throw new Error(`PINNING_RULE_ALREADY_EXISTS: A pinning rule with id "${rule.id}" already exists. Use updateRule to modify it.`);
    }
    store2.rules.set(rule.id, rule);
  }
  function updateRule(store2, rule) {
    if (!store2.rules.has(rule.id)) {
      throw new Error(`PINNING_RULE_NOT_FOUND: Cannot update pinning rule with id "${rule.id}" because it does not exist. Use addRule to create it.`);
    }
    store2.rules.set(rule.id, rule);
  }
  function removeRule(store2, ruleId) {
    return store2.rules.delete(ruleId);
  }
  function getRule(store2, ruleId) {
    return store2.rules.get(ruleId);
  }
  function getAllRules(store2) {
    return Array.from(store2.rules.values());
  }
  function matchesCondition(term, condition) {
    const normalizedTerm = term.toLowerCase().trim();
    const normalizedPattern = condition.pattern.toLowerCase().trim();
    switch (condition.anchoring) {
      case "is":
        return normalizedTerm === normalizedPattern;
      case "starts_with":
        return normalizedTerm.startsWith(normalizedPattern);
      case "contains":
        return normalizedTerm.includes(normalizedPattern);
      default:
        return false;
    }
  }
  function matchesRule(term, rule) {
    if (!term) {
      return false;
    }
    return rule.conditions.every((condition) => matchesCondition(term, condition));
  }
  function getMatchingRules(store2, term) {
    if (!term) {
      return [];
    }
    const matchingRules = [];
    for (const rule of store2.rules.values()) {
      if (matchesRule(term, rule)) {
        matchingRules.push(rule);
      }
    }
    return matchingRules;
  }
  function load5(sharedInternalDocumentStore, raw) {
    const rawStore = raw;
    return {
      sharedInternalDocumentStore,
      rules: new Map(rawStore?.rules ?? [])
    };
  }
  function save5(store2) {
    return {
      rules: Array.from(store2.rules.entries())
    };
  }
  function createPinning() {
    return {
      create: create4,
      addRule,
      updateRule,
      removeRule,
      getRule,
      getAllRules,
      getMatchingRules,
      load: load5,
      save: save5
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/create.js
  function validateComponents(components) {
    const defaultComponents = {
      formatElapsedTime,
      getDocumentIndexId,
      getDocumentProperties,
      validateSchema
    };
    for (const rawKey of FUNCTION_COMPONENTS) {
      const key = rawKey;
      if (components[key]) {
        if (typeof components[key] !== "function") {
          throw createError("COMPONENT_MUST_BE_FUNCTION", key);
        }
      } else {
        components[key] = defaultComponents[key];
      }
    }
    for (const rawKey of Object.keys(components)) {
      if (!OBJECT_COMPONENTS.includes(rawKey) && !FUNCTION_COMPONENTS.includes(rawKey)) {
        throw createError("UNSUPPORTED_COMPONENT", rawKey);
      }
    }
  }
  function create5({ schema, sort, language, components, id, plugins }) {
    if (!components) {
      components = {};
    }
    for (const plugin of plugins ?? []) {
      if (!("getComponents" in plugin)) {
        continue;
      }
      if (typeof plugin.getComponents !== "function") {
        continue;
      }
      const pluginComponents = plugin.getComponents(schema);
      const keys = Object.keys(pluginComponents);
      for (const key of keys) {
        if (components[key]) {
          throw createError("PLUGIN_COMPONENT_CONFLICT", key, plugin.name);
        }
      }
      components = {
        ...components,
        ...pluginComponents
      };
    }
    if (!id) {
      id = uniqueId();
    }
    let tokenizer3 = components.tokenizer;
    let index2 = components.index;
    let documentsStore = components.documentsStore;
    let sorter = components.sorter;
    let pinning = components.pinning;
    if (!tokenizer3) {
      tokenizer3 = createTokenizer2({ language: language ?? "english" });
    } else if (!tokenizer3.tokenize) {
      tokenizer3 = createTokenizer2(tokenizer3);
    } else {
      const customTokenizer = tokenizer3;
      tokenizer3 = customTokenizer;
    }
    if (components.tokenizer && language) {
      throw createError("NO_LANGUAGE_WITH_CUSTOM_TOKENIZER");
    }
    const internalDocumentStore = createInternalDocumentIDStore();
    index2 || (index2 = createIndex());
    sorter || (sorter = createSorter());
    documentsStore || (documentsStore = createDocumentsStore());
    pinning || (pinning = createPinning());
    validateComponents(components);
    const { getDocumentProperties: getDocumentProperties2, getDocumentIndexId: getDocumentIndexId2, validateSchema: validateSchema2, formatElapsedTime: formatElapsedTime2 } = components;
    const orama = {
      data: {},
      caches: {},
      schema,
      tokenizer: tokenizer3,
      index: index2,
      sorter,
      documentsStore,
      pinning,
      internalDocumentIDStore: internalDocumentStore,
      getDocumentProperties: getDocumentProperties2,
      getDocumentIndexId: getDocumentIndexId2,
      validateSchema: validateSchema2,
      beforeInsert: [],
      afterInsert: [],
      beforeRemove: [],
      afterRemove: [],
      beforeUpdate: [],
      afterUpdate: [],
      beforeUpsert: [],
      afterUpsert: [],
      beforeSearch: [],
      afterSearch: [],
      beforeInsertMultiple: [],
      afterInsertMultiple: [],
      beforeRemoveMultiple: [],
      afterRemoveMultiple: [],
      beforeUpdateMultiple: [],
      afterUpdateMultiple: [],
      beforeUpsertMultiple: [],
      afterUpsertMultiple: [],
      afterCreate: [],
      formatElapsedTime: formatElapsedTime2,
      id,
      plugins,
      version: getVersion()
    };
    orama.data = {
      index: orama.index.create(orama, internalDocumentStore, schema),
      docs: orama.documentsStore.create(orama, internalDocumentStore),
      sorting: orama.sorter.create(orama, internalDocumentStore, schema, sort),
      pinning: orama.pinning.create(internalDocumentStore)
    };
    for (const hook of AVAILABLE_PLUGIN_HOOKS) {
      orama[hook] = (orama[hook] ?? []).concat(getAllPluginsByHook(orama, hook));
    }
    const afterCreate = orama["afterCreate"];
    if (afterCreate) {
      runAfterCreate(afterCreate, orama);
    }
    return orama;
  }
  function getVersion() {
    return "{{VERSION}}";
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/docs.js
  function getByID(db, id) {
    return db.documentsStore.get(db.data.docs, id);
  }
  function count2(db) {
    return db.documentsStore.count(db.data.docs);
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/insert.js
  function insert3(orama, doc, language, skipHooks, options) {
    const errorProperty = orama.validateSchema(doc, orama.schema);
    if (errorProperty) {
      throw createError("SCHEMA_VALIDATION_FAILURE", errorProperty);
    }
    const asyncNeeded = isAsyncFunction(orama.beforeInsert) || isAsyncFunction(orama.afterInsert) || isAsyncFunction(orama.index.beforeInsert) || isAsyncFunction(orama.index.insert) || isAsyncFunction(orama.index.afterInsert);
    if (asyncNeeded) {
      return innerInsertAsync(orama, doc, language, skipHooks, options);
    }
    return innerInsertSync(orama, doc, language, skipHooks, options);
  }
  var ENUM_TYPE = /* @__PURE__ */ new Set(["enum", "enum[]"]);
  var STRING_NUMBER_TYPE = /* @__PURE__ */ new Set(["string", "number"]);
  async function innerInsertAsync(orama, doc, language, skipHooks, options) {
    const { index: index2, docs } = orama.data;
    const id = orama.getDocumentIndexId(doc);
    if (typeof id !== "string") {
      throw createError("DOCUMENT_ID_MUST_BE_STRING", typeof id);
    }
    const internalId = getInternalDocumentId(orama.internalDocumentIDStore, id);
    if (!skipHooks) {
      await runSingleHook(orama.beforeInsert, orama, id, doc);
    }
    if (!orama.documentsStore.store(docs, id, internalId, doc)) {
      throw createError("DOCUMENT_ALREADY_EXISTS", id);
    }
    const docsCount = orama.documentsStore.count(docs);
    const indexableProperties = orama.index.getSearchableProperties(index2);
    const indexablePropertiesWithTypes = orama.index.getSearchablePropertiesWithTypes(index2);
    const indexableValues = orama.getDocumentProperties(doc, indexableProperties);
    for (const [key, value] of Object.entries(indexableValues)) {
      if (typeof value === "undefined")
        continue;
      const actualType = typeof value;
      const expectedType = indexablePropertiesWithTypes[key];
      validateDocumentProperty(actualType, expectedType, key, value);
    }
    await indexAndSortDocument(orama, id, indexableProperties, indexableValues, docsCount, language, doc, options);
    if (!skipHooks) {
      await runSingleHook(orama.afterInsert, orama, id, doc);
    }
    return id;
  }
  function innerInsertSync(orama, doc, language, skipHooks, options) {
    const { index: index2, docs } = orama.data;
    const id = orama.getDocumentIndexId(doc);
    if (typeof id !== "string") {
      throw createError("DOCUMENT_ID_MUST_BE_STRING", typeof id);
    }
    const internalId = getInternalDocumentId(orama.internalDocumentIDStore, id);
    if (!skipHooks) {
      runSingleHook(orama.beforeInsert, orama, id, doc);
    }
    if (!orama.documentsStore.store(docs, id, internalId, doc)) {
      throw createError("DOCUMENT_ALREADY_EXISTS", id);
    }
    const docsCount = orama.documentsStore.count(docs);
    const indexableProperties = orama.index.getSearchableProperties(index2);
    const indexablePropertiesWithTypes = orama.index.getSearchablePropertiesWithTypes(index2);
    const indexableValues = orama.getDocumentProperties(doc, indexableProperties);
    for (const [key, value] of Object.entries(indexableValues)) {
      if (typeof value === "undefined")
        continue;
      const actualType = typeof value;
      const expectedType = indexablePropertiesWithTypes[key];
      validateDocumentProperty(actualType, expectedType, key, value);
    }
    indexAndSortDocumentSync(orama, id, indexableProperties, indexableValues, docsCount, language, doc, options);
    if (!skipHooks) {
      runSingleHook(orama.afterInsert, orama, id, doc);
    }
    return id;
  }
  function validateDocumentProperty(actualType, expectedType, key, value) {
    if (isGeoPointType(expectedType) && typeof value === "object" && typeof value.lon === "number" && typeof value.lat === "number") {
      return;
    }
    if (isVectorType(expectedType) && Array.isArray(value))
      return;
    if (isArrayType(expectedType) && Array.isArray(value))
      return;
    if (ENUM_TYPE.has(expectedType) && STRING_NUMBER_TYPE.has(actualType))
      return;
    if (actualType !== expectedType) {
      throw createError("INVALID_DOCUMENT_PROPERTY", key, expectedType, actualType);
    }
  }
  async function indexAndSortDocument(orama, id, indexableProperties, indexableValues, docsCount, language, doc, options) {
    for (const prop of indexableProperties) {
      const value = indexableValues[prop];
      if (typeof value === "undefined")
        continue;
      const expectedType = orama.index.getSearchablePropertiesWithTypes(orama.data.index)[prop];
      await orama.index.beforeInsert?.(orama.data.index, prop, id, value, expectedType, language, orama.tokenizer, docsCount);
      const internalId = orama.internalDocumentIDStore.idToInternalId.get(id);
      await orama.index.insert(orama.index, orama.data.index, prop, id, internalId, value, expectedType, language, orama.tokenizer, docsCount, options);
      await orama.index.afterInsert?.(orama.data.index, prop, id, value, expectedType, language, orama.tokenizer, docsCount);
    }
    const sortableProperties = orama.sorter.getSortableProperties(orama.data.sorting);
    const sortableValues = orama.getDocumentProperties(doc, sortableProperties);
    for (const prop of sortableProperties) {
      const value = sortableValues[prop];
      if (typeof value === "undefined")
        continue;
      const expectedType = orama.sorter.getSortablePropertiesWithTypes(orama.data.sorting)[prop];
      orama.sorter.insert(orama.data.sorting, prop, id, value, expectedType, language);
    }
  }
  function indexAndSortDocumentSync(orama, id, indexableProperties, indexableValues, docsCount, language, doc, options) {
    for (const prop of indexableProperties) {
      const value = indexableValues[prop];
      if (typeof value === "undefined")
        continue;
      const expectedType = orama.index.getSearchablePropertiesWithTypes(orama.data.index)[prop];
      const internalDocumentId = getInternalDocumentId(orama.internalDocumentIDStore, id);
      orama.index.beforeInsert?.(orama.data.index, prop, id, value, expectedType, language, orama.tokenizer, docsCount);
      orama.index.insert(orama.index, orama.data.index, prop, id, internalDocumentId, value, expectedType, language, orama.tokenizer, docsCount, options);
      orama.index.afterInsert?.(orama.data.index, prop, id, value, expectedType, language, orama.tokenizer, docsCount);
    }
    const sortableProperties = orama.sorter.getSortableProperties(orama.data.sorting);
    const sortableValues = orama.getDocumentProperties(doc, sortableProperties);
    for (const prop of sortableProperties) {
      const value = sortableValues[prop];
      if (typeof value === "undefined")
        continue;
      const expectedType = orama.sorter.getSortablePropertiesWithTypes(orama.data.sorting)[prop];
      orama.sorter.insert(orama.data.sorting, prop, id, value, expectedType, language);
    }
  }
  function insertMultiple(orama, docs, batchSize, language, skipHooks, timeout) {
    const asyncNeeded = isAsyncFunction(orama.afterInsertMultiple) || isAsyncFunction(orama.beforeInsertMultiple) || isAsyncFunction(orama.index.beforeInsert) || isAsyncFunction(orama.index.insert) || isAsyncFunction(orama.index.afterInsert);
    if (asyncNeeded) {
      return innerInsertMultipleAsync(orama, docs, batchSize, language, skipHooks, timeout);
    }
    return innerInsertMultipleSync(orama, docs, batchSize, language, skipHooks, timeout);
  }
  async function innerInsertMultipleAsync(orama, docs, batchSize = 1e3, language, skipHooks, timeout = 0) {
    const ids = [];
    const processNextBatch = async (startIndex) => {
      const endIndex = Math.min(startIndex + batchSize, docs.length);
      const batch = docs.slice(startIndex, endIndex);
      for (const doc of batch) {
        const options = { avlRebalanceThreshold: batch.length };
        const id = await insert3(orama, doc, language, skipHooks, options);
        ids.push(id);
      }
      return endIndex;
    };
    const processAllBatches = async () => {
      let currentIndex = 0;
      while (currentIndex < docs.length) {
        const startTime = Date.now();
        currentIndex = await processNextBatch(currentIndex);
        if (timeout > 0) {
          const elapsedTime = Date.now() - startTime;
          const waitTime = timeout - elapsedTime;
          if (waitTime > 0) {
            sleep(waitTime);
          }
        }
      }
    };
    await processAllBatches();
    if (!skipHooks) {
      await runMultipleHook(orama.afterInsertMultiple, orama, docs);
    }
    return ids;
  }
  function innerInsertMultipleSync(orama, docs, batchSize = 1e3, language, skipHooks, timeout = 0) {
    const ids = [];
    let i = 0;
    function processNextBatch() {
      const batch = docs.slice(i * batchSize, (i + 1) * batchSize);
      if (batch.length === 0)
        return false;
      for (const doc of batch) {
        const options = { avlRebalanceThreshold: batch.length };
        const id = insert3(orama, doc, language, skipHooks, options);
        ids.push(id);
      }
      i++;
      return true;
    }
    function processAllBatches() {
      const startTime = Date.now();
      while (true) {
        const hasMoreBatches = processNextBatch();
        if (!hasMoreBatches)
          break;
        if (timeout > 0) {
          const elapsedTime = Date.now() - startTime;
          if (elapsedTime >= timeout) {
            const remainingTime = timeout - elapsedTime % timeout;
            if (remainingTime > 0) {
              sleep(remainingTime);
            }
          }
        }
      }
    }
    processAllBatches();
    if (!skipHooks) {
      runMultipleHook(orama.afterInsertMultiple, orama, docs);
    }
    return ids;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/constants.js
  var MODE_FULLTEXT_SEARCH = "fulltext";
  var MODE_HYBRID_SEARCH = "hybrid";
  var MODE_VECTOR_SEARCH = "vector";

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/facets.js
  function sortAsc(a, b) {
    return a[1] - b[1];
  }
  function sortDesc(a, b) {
    return b[1] - a[1];
  }
  function sortingPredicateBuilder(order2 = "desc") {
    return order2.toLowerCase() === "asc" ? sortAsc : sortDesc;
  }
  function getFacets(orama, results, facetsConfig) {
    const facets = {};
    const allIDs = results.map(([id]) => id);
    const allDocs = orama.documentsStore.getMultiple(orama.data.docs, allIDs);
    const facetKeys = Object.keys(facetsConfig);
    const properties = orama.index.getSearchablePropertiesWithTypes(orama.data.index);
    for (const facet of facetKeys) {
      let values;
      if (properties[facet] === "number") {
        const { ranges } = facetsConfig[facet];
        const rangesLength = ranges.length;
        const tmp = Array.from({ length: rangesLength });
        for (let i = 0; i < rangesLength; i++) {
          const range = ranges[i];
          tmp[i] = [`${range.from}-${range.to}`, 0];
        }
        values = Object.fromEntries(tmp);
      }
      facets[facet] = {
        count: 0,
        values: values ?? {}
      };
    }
    const allDocsLength = allDocs.length;
    for (let i = 0; i < allDocsLength; i++) {
      const doc = allDocs[i];
      for (const facet of facetKeys) {
        const facetValue = facet.includes(".") ? getNested(doc, facet) : doc[facet];
        const propertyType = properties[facet];
        const facetValues = facets[facet].values;
        switch (propertyType) {
          case "number": {
            const ranges = facetsConfig[facet].ranges;
            calculateNumberFacetBuilder(ranges, facetValues)(facetValue);
            break;
          }
          case "number[]": {
            const alreadyInsertedValues = /* @__PURE__ */ new Set();
            const ranges = facetsConfig[facet].ranges;
            const calculateNumberFacet = calculateNumberFacetBuilder(ranges, facetValues, alreadyInsertedValues);
            for (const v2 of facetValue) {
              calculateNumberFacet(v2);
            }
            break;
          }
          case "boolean":
          case "enum":
          case "string": {
            calculateBooleanStringOrEnumFacetBuilder(facetValues, propertyType)(facetValue);
            break;
          }
          case "boolean[]":
          case "enum[]":
          case "string[]": {
            const alreadyInsertedValues = /* @__PURE__ */ new Set();
            const innerType = propertyType === "boolean[]" ? "boolean" : "string";
            const calculateBooleanStringOrEnumFacet = calculateBooleanStringOrEnumFacetBuilder(facetValues, innerType, alreadyInsertedValues);
            for (const v2 of facetValue) {
              calculateBooleanStringOrEnumFacet(v2);
            }
            break;
          }
          default:
            throw createError("FACET_NOT_SUPPORTED", propertyType);
        }
      }
    }
    for (const facet of facetKeys) {
      const currentFacet = facets[facet];
      currentFacet.count = Object.keys(currentFacet.values).length;
      if (properties[facet] === "string") {
        const stringFacetDefinition = facetsConfig[facet];
        const sortingPredicate = sortingPredicateBuilder(stringFacetDefinition.sort);
        currentFacet.values = Object.fromEntries(Object.entries(currentFacet.values).sort(sortingPredicate).slice(stringFacetDefinition.offset ?? 0, stringFacetDefinition.limit ?? 10));
      }
    }
    return facets;
  }
  function calculateNumberFacetBuilder(ranges, values, alreadyInsertedValues) {
    return (facetValue) => {
      for (const range of ranges) {
        const value = `${range.from}-${range.to}`;
        if (alreadyInsertedValues?.has(value)) {
          continue;
        }
        if (facetValue >= range.from && facetValue <= range.to) {
          if (values[value] === void 0) {
            values[value] = 1;
          } else {
            values[value]++;
            alreadyInsertedValues?.add(value);
          }
        }
      }
    };
  }
  function calculateBooleanStringOrEnumFacetBuilder(values, propertyType, alreadyInsertedValues) {
    const defaultValue = propertyType === "boolean" ? "false" : "";
    return (facetValue) => {
      const value = facetValue?.toString() ?? defaultValue;
      if (alreadyInsertedValues?.has(value)) {
        return;
      }
      values[value] = (values[value] ?? 0) + 1;
      alreadyInsertedValues?.add(value);
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/groups.js
  var DEFAULT_REDUCE = {
    reducer: (_, acc, res, index2) => {
      acc[index2] = res;
      return acc;
    },
    getInitialValue: (length) => Array.from({ length })
  };
  var ALLOWED_TYPES = ["string", "number", "boolean"];
  function getGroups(orama, results, groupBy) {
    const properties = groupBy.properties;
    const propertiesLength = properties.length;
    const schemaProperties = orama.index.getSearchablePropertiesWithTypes(orama.data.index);
    for (let i = 0; i < propertiesLength; i++) {
      const property2 = properties[i];
      if (typeof schemaProperties[property2] === "undefined") {
        throw createError("UNKNOWN_GROUP_BY_PROPERTY", property2);
      }
      if (!ALLOWED_TYPES.includes(schemaProperties[property2])) {
        throw createError("INVALID_GROUP_BY_PROPERTY", property2, ALLOWED_TYPES.join(", "), schemaProperties[property2]);
      }
    }
    const allIDs = results.map(([id]) => getDocumentIdFromInternalId(orama.internalDocumentIDStore, id));
    const allDocs = orama.documentsStore.getMultiple(orama.data.docs, allIDs);
    const allDocsLength = allDocs.length;
    const returnedCount = groupBy.maxResult || Number.MAX_SAFE_INTEGER;
    const listOfValues = [];
    const g = {};
    for (let i = 0; i < propertiesLength; i++) {
      const groupByKey = properties[i];
      const group = {
        property: groupByKey,
        perValue: {}
      };
      const values = /* @__PURE__ */ new Set();
      for (let j = 0; j < allDocsLength; j++) {
        const doc = allDocs[j];
        const value = getNested(doc, groupByKey);
        if (typeof value === "undefined") {
          continue;
        }
        const keyValue = typeof value !== "boolean" ? value : "" + value;
        const perValue = group.perValue[keyValue] ?? {
          indexes: [],
          count: 0
        };
        if (perValue.count >= returnedCount) {
          continue;
        }
        perValue.indexes.push(j);
        perValue.count++;
        group.perValue[keyValue] = perValue;
        values.add(value);
      }
      listOfValues.push(Array.from(values));
      g[groupByKey] = group;
    }
    const combinations = calculateCombination(listOfValues);
    const combinationsLength = combinations.length;
    const groups = [];
    for (let i = 0; i < combinationsLength; i++) {
      const combination = combinations[i];
      const combinationLength = combination.length;
      const group = {
        values: [],
        indexes: []
      };
      const indexes = [];
      for (let j = 0; j < combinationLength; j++) {
        const value = combination[j];
        const property2 = properties[j];
        indexes.push(g[property2].perValue[typeof value !== "boolean" ? value : "" + value].indexes);
        group.values.push(value);
      }
      group.indexes = intersect(indexes).sort((a, b) => a - b);
      if (group.indexes.length === 0) {
        continue;
      }
      groups.push(group);
    }
    const groupsLength = groups.length;
    const res = Array.from({ length: groupsLength });
    for (let i = 0; i < groupsLength; i++) {
      const group = groups[i];
      const reduce = groupBy.reduce || DEFAULT_REDUCE;
      const docs = group.indexes.map((index2) => {
        return {
          id: allIDs[index2],
          score: results[index2][1],
          document: allDocs[index2]
        };
      });
      const func = reduce.reducer.bind(null, group.values);
      const initialValue = reduce.getInitialValue(group.indexes.length);
      const aggregationValue = docs.reduce(func, initialValue);
      res[i] = {
        values: group.values,
        result: aggregationValue
      };
    }
    return res;
  }
  function calculateCombination(arrs, index2 = 0) {
    if (index2 + 1 === arrs.length)
      return arrs[index2].map((item) => [item]);
    const head = arrs[index2];
    const c2 = calculateCombination(arrs, index2 + 1);
    const combinations = [];
    for (const value of head) {
      for (const combination of c2) {
        const result = [value];
        safeArrayPush(result, combination);
        combinations.push(result);
      }
    }
    return combinations;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/components/pinning-manager.js
  function applyPinningRules(orama, pinningStore, uniqueDocsArray, searchTerm) {
    const matchingRules = getMatchingRules(pinningStore, searchTerm);
    if (matchingRules.length === 0) {
      return uniqueDocsArray;
    }
    const allPromotions = matchingRules.flatMap((rule) => rule.consequence.promote);
    allPromotions.sort((a, b) => a.position - b.position);
    const pinnedInternalIds = /* @__PURE__ */ new Set();
    const promotionsMap = /* @__PURE__ */ new Map();
    const positionsTaken = /* @__PURE__ */ new Set();
    for (const promotion of allPromotions) {
      const internalId = getInternalDocumentId(orama.internalDocumentIDStore, promotion.doc_id);
      if (internalId === void 0) {
        continue;
      }
      if (promotionsMap.has(internalId)) {
        const existingPosition = promotionsMap.get(internalId);
        if (promotion.position < existingPosition) {
          promotionsMap.set(internalId, promotion.position);
        }
        continue;
      }
      if (positionsTaken.has(promotion.position)) {
        continue;
      }
      pinnedInternalIds.add(internalId);
      promotionsMap.set(internalId, promotion.position);
      positionsTaken.add(promotion.position);
    }
    if (promotionsMap.size === 0) {
      return uniqueDocsArray;
    }
    const unpinnedResults = uniqueDocsArray.filter(([id]) => !pinnedInternalIds.has(id));
    const BASE_PIN_SCORE = 1e6;
    const pinnedResults = [];
    for (const [internalId, position2] of promotionsMap.entries()) {
      const existingResult = uniqueDocsArray.find(([id]) => id === internalId);
      if (existingResult) {
        pinnedResults.push([internalId, BASE_PIN_SCORE - position2]);
      } else {
        const doc = orama.documentsStore.get(orama.data.docs, internalId);
        if (doc) {
          pinnedResults.push([internalId, 0]);
        }
      }
    }
    pinnedResults.sort((a, b) => {
      const posA = promotionsMap.get(a[0]) ?? Infinity;
      const posB = promotionsMap.get(b[0]) ?? Infinity;
      return posA - posB;
    });
    const finalResults = [];
    const pinnedByPosition = /* @__PURE__ */ new Map();
    for (const pinnedResult of pinnedResults) {
      const position2 = promotionsMap.get(pinnedResult[0]);
      pinnedByPosition.set(position2, pinnedResult);
    }
    let unpinnedIndex = 0;
    let currentPosition = 0;
    while (currentPosition < unpinnedResults.length + pinnedResults.length) {
      if (pinnedByPosition.has(currentPosition)) {
        finalResults.push(pinnedByPosition.get(currentPosition));
        currentPosition++;
      } else if (unpinnedIndex < unpinnedResults.length) {
        finalResults.push(unpinnedResults[unpinnedIndex]);
        unpinnedIndex++;
        currentPosition++;
      } else {
        break;
      }
    }
    for (const [position2, pinnedResult] of pinnedByPosition.entries()) {
      if (position2 >= finalResults.length) {
        finalResults.push(pinnedResult);
      }
    }
    return finalResults;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/search-fulltext.js
  function innerFullTextSearch(orama, params, language) {
    const { term, properties } = params;
    const index2 = orama.data.index;
    let propertiesToSearch = orama.caches["propertiesToSearch"];
    if (!propertiesToSearch) {
      const propertiesToSearchWithTypes = orama.index.getSearchablePropertiesWithTypes(index2);
      propertiesToSearch = orama.index.getSearchableProperties(index2);
      propertiesToSearch = propertiesToSearch.filter((prop) => propertiesToSearchWithTypes[prop].startsWith("string"));
      orama.caches["propertiesToSearch"] = propertiesToSearch;
    }
    if (properties && properties !== "*") {
      for (const prop of properties) {
        if (!propertiesToSearch.includes(prop)) {
          throw createError("UNKNOWN_INDEX", prop, propertiesToSearch.join(", "));
        }
      }
      propertiesToSearch = propertiesToSearch.filter((prop) => properties.includes(prop));
    }
    const hasFilters = Object.keys(params.where ?? {}).length > 0;
    let whereFiltersIDs;
    if (hasFilters) {
      whereFiltersIDs = orama.index.searchByWhereClause(index2, orama.tokenizer, params.where, language);
    }
    let uniqueDocsIDs;
    const threshold = params.threshold !== void 0 && params.threshold !== null ? params.threshold : 1;
    if (term || properties) {
      const docsCount = count2(orama);
      uniqueDocsIDs = orama.index.search(index2, term || "", orama.tokenizer, language, propertiesToSearch, params.exact || false, params.tolerance || 0, params.boost || {}, applyDefault(params.relevance), docsCount, whereFiltersIDs, threshold);
      if (params.exact && term) {
        const searchTerms = term.trim().split(/\s+/);
        uniqueDocsIDs = uniqueDocsIDs.filter(([docId]) => {
          const doc = orama.documentsStore.get(orama.data.docs, docId);
          if (!doc)
            return false;
          for (const prop of propertiesToSearch) {
            const propValue = getPropValue(doc, prop);
            if (typeof propValue === "string") {
              const hasAllTerms = searchTerms.every((searchTerm) => {
                const regex2 = new RegExp(`\\b${escapeRegex(searchTerm)}\\b`);
                return regex2.test(propValue);
              });
              if (hasAllTerms) {
                return true;
              }
            }
          }
          return false;
        });
      }
    } else {
      if (hasFilters) {
        const geoResults = searchByGeoWhereClause(index2, params.where);
        if (geoResults) {
          uniqueDocsIDs = geoResults;
        } else {
          const docIds = whereFiltersIDs ? Array.from(whereFiltersIDs) : [];
          uniqueDocsIDs = docIds.map((k) => [+k, 0]);
        }
      } else {
        const docIds = Object.keys(orama.documentsStore.getAll(orama.data.docs));
        uniqueDocsIDs = docIds.map((k) => [+k, 0]);
      }
    }
    return uniqueDocsIDs;
  }
  function escapeRegex(str) {
    return str.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
  }
  function getPropValue(obj, path2) {
    const keys = path2.split(".");
    let value = obj;
    for (const key of keys) {
      if (value && typeof value === "object" && key in value) {
        value = value[key];
      } else {
        return void 0;
      }
    }
    return value;
  }
  function fullTextSearch(orama, params, language) {
    const timeStart = getNanosecondsTime();
    function performSearchLogic() {
      const vectorProperties = Object.keys(orama.data.index.vectorIndexes);
      const shouldCalculateFacets = params.facets && Object.keys(params.facets).length > 0;
      const { limit = 10, offset: offset2 = 0, distinctOn, includeVectors = false } = params;
      const isPreflight = params.preflight === true;
      let uniqueDocsArray = innerFullTextSearch(orama, params, language);
      if (params.sortBy) {
        if (typeof params.sortBy === "function") {
          const ids = uniqueDocsArray.map(([id]) => id);
          const docs = orama.documentsStore.getMultiple(orama.data.docs, ids);
          const docsWithIdAndScore = docs.map((d, i) => [
            uniqueDocsArray[i][0],
            uniqueDocsArray[i][1],
            d
          ]);
          docsWithIdAndScore.sort(params.sortBy);
          uniqueDocsArray = docsWithIdAndScore.map(([id, score]) => [id, score]);
        } else {
          uniqueDocsArray = orama.sorter.sortBy(orama.data.sorting, uniqueDocsArray, params.sortBy).map(([id, score]) => [getInternalDocumentId(orama.internalDocumentIDStore, id), score]);
        }
      } else {
        uniqueDocsArray = uniqueDocsArray.sort(sortTokenScorePredicate);
      }
      uniqueDocsArray = applyPinningRules(orama, orama.data.pinning, uniqueDocsArray, params.term);
      let results;
      if (!isPreflight) {
        results = distinctOn ? fetchDocumentsWithDistinct(orama, uniqueDocsArray, offset2, limit, distinctOn) : fetchDocuments(orama, uniqueDocsArray, offset2, limit);
      }
      const searchResult = {
        elapsed: {
          formatted: "",
          raw: 0
        },
        hits: [],
        count: uniqueDocsArray.length
      };
      if (typeof results !== "undefined") {
        searchResult.hits = results.filter(Boolean);
        if (!includeVectors) {
          removeVectorsFromHits(searchResult, vectorProperties);
        }
      }
      if (shouldCalculateFacets) {
        const facets = getFacets(orama, uniqueDocsArray, params.facets);
        searchResult.facets = facets;
      }
      if (params.groupBy) {
        searchResult.groups = getGroups(orama, uniqueDocsArray, params.groupBy);
      }
      searchResult.elapsed = orama.formatElapsedTime(getNanosecondsTime() - timeStart);
      return searchResult;
    }
    async function executeSearchAsync() {
      if (orama.beforeSearch) {
        await runBeforeSearch(orama.beforeSearch, orama, params, language);
      }
      const searchResult = performSearchLogic();
      if (orama.afterSearch) {
        await runAfterSearch(orama.afterSearch, orama, params, language, searchResult);
      }
      return searchResult;
    }
    const asyncNeeded = orama.beforeSearch?.length || orama.afterSearch?.length;
    if (asyncNeeded) {
      return executeSearchAsync();
    }
    return performSearchLogic();
  }
  var defaultBM25Params = {
    k: 1.2,
    b: 0.75,
    d: 0.5
  };
  function applyDefault(bm25Relevance) {
    const r = bm25Relevance ?? {};
    r.k = r.k ?? defaultBM25Params.k;
    r.b = r.b ?? defaultBM25Params.b;
    r.d = r.d ?? defaultBM25Params.d;
    return r;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/search-vector.js
  function innerVectorSearch(orama, params, language) {
    const vector = params.vector;
    if (vector && (!("value" in vector) || !("property" in vector))) {
      throw createError("INVALID_VECTOR_INPUT", Object.keys(vector).join(", "));
    }
    const vectorIndex = orama.data.index.vectorIndexes[vector.property];
    if (!vectorIndex) {
      throw createError("UNKNOWN_VECTOR_PROPERTY", vector.property);
    }
    const vectorSize = vectorIndex.node.size;
    if (vector?.value.length !== vectorSize) {
      if (vector?.property === void 0 || vector?.value.length === void 0) {
        throw createError("INVALID_INPUT_VECTOR", "undefined", vectorSize, "undefined");
      }
      throw createError("INVALID_INPUT_VECTOR", vector.property, vectorSize, vector.value.length);
    }
    const index2 = orama.data.index;
    let whereFiltersIDs;
    const hasFilters = Object.keys(params.where ?? {}).length > 0;
    if (hasFilters) {
      whereFiltersIDs = orama.index.searchByWhereClause(index2, orama.tokenizer, params.where, language);
    }
    return vectorIndex.node.find(vector.value, params.similarity ?? DEFAULT_SIMILARITY, whereFiltersIDs);
  }
  function searchVector(orama, params, language = "english") {
    const timeStart = getNanosecondsTime();
    function performSearchLogic() {
      let results = innerVectorSearch(orama, params, language).sort(sortTokenScorePredicate);
      results = applyPinningRules(orama, orama.data.pinning, results, void 0);
      let facetsResults = [];
      const shouldCalculateFacets = params.facets && Object.keys(params.facets).length > 0;
      if (shouldCalculateFacets) {
        const facets = getFacets(orama, results, params.facets);
        facetsResults = facets;
      }
      const vectorProperty = params.vector.property;
      const includeVectors = params.includeVectors ?? false;
      const limit = params.limit ?? 10;
      const offset2 = params.offset ?? 0;
      const docs = Array.from({ length: limit });
      for (let i = 0; i < limit; i++) {
        const result = results[i + offset2];
        if (!result) {
          break;
        }
        const doc = orama.data.docs.docs[result[0]];
        if (doc) {
          if (!includeVectors) {
            doc[vectorProperty] = null;
          }
          const newDoc = {
            id: getDocumentIdFromInternalId(orama.internalDocumentIDStore, result[0]),
            score: result[1],
            document: doc
          };
          docs[i] = newDoc;
        }
      }
      let groups = [];
      if (params.groupBy) {
        groups = getGroups(orama, results, params.groupBy);
      }
      const timeEnd = getNanosecondsTime();
      const elapsedTime = timeEnd - timeStart;
      return {
        count: results.length,
        hits: docs.filter(Boolean),
        elapsed: {
          raw: Number(elapsedTime),
          formatted: formatNanoseconds(elapsedTime)
        },
        ...facetsResults ? { facets: facetsResults } : {},
        ...groups ? { groups } : {}
      };
    }
    async function executeSearchAsync() {
      if (orama.beforeSearch) {
        await runBeforeSearch(orama.beforeSearch, orama, params, language);
      }
      const results = performSearchLogic();
      if (orama.afterSearch) {
        await runAfterSearch(orama.afterSearch, orama, params, language, results);
      }
      return results;
    }
    const asyncNeeded = orama.beforeSearch?.length || orama.afterSearch?.length;
    if (asyncNeeded) {
      return executeSearchAsync();
    }
    return performSearchLogic();
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/search-hybrid.js
  function innerHybridSearch(orama, params, language) {
    const fullTextIDs = minMaxScoreNormalization(innerFullTextSearch(orama, params, language));
    const vectorIDs = innerVectorSearch(orama, params, language);
    const hybridWeights = params.hybridWeights;
    return mergeAndRankResults(fullTextIDs, vectorIDs, params.term ?? "", hybridWeights);
  }
  function hybridSearch(orama, params, language) {
    const timeStart = getNanosecondsTime();
    function performSearchLogic() {
      let uniqueTokenScores = innerHybridSearch(orama, params, language);
      uniqueTokenScores = applyPinningRules(orama, orama.data.pinning, uniqueTokenScores, params.term);
      let facetsResults;
      const shouldCalculateFacets = params.facets && Object.keys(params.facets).length > 0;
      if (shouldCalculateFacets) {
        facetsResults = getFacets(orama, uniqueTokenScores, params.facets);
      }
      let groups;
      if (params.groupBy) {
        groups = getGroups(orama, uniqueTokenScores, params.groupBy);
      }
      const offset2 = params.offset ?? 0;
      const limit = params.limit ?? 10;
      const results = fetchDocuments(orama, uniqueTokenScores, offset2, limit).filter(Boolean);
      const timeEnd = getNanosecondsTime();
      const returningResults = {
        count: uniqueTokenScores.length,
        elapsed: {
          raw: Number(timeEnd - timeStart),
          formatted: formatNanoseconds(timeEnd - timeStart)
        },
        hits: results,
        ...facetsResults ? { facets: facetsResults } : {},
        ...groups ? { groups } : {}
      };
      const includeVectors = params.includeVectors ?? false;
      if (!includeVectors) {
        const vectorProperties = Object.keys(orama.data.index.vectorIndexes);
        removeVectorsFromHits(returningResults, vectorProperties);
      }
      return returningResults;
    }
    async function executeSearchAsync() {
      if (orama.beforeSearch) {
        await runBeforeSearch(orama.beforeSearch, orama, params, language);
      }
      const results = performSearchLogic();
      if (orama.afterSearch) {
        await runAfterSearch(orama.afterSearch, orama, params, language, results);
      }
      return results;
    }
    const asyncNeeded = orama.beforeSearch?.length || orama.afterSearch?.length;
    if (asyncNeeded) {
      return executeSearchAsync();
    }
    return performSearchLogic();
  }
  function extractScore(token) {
    return token[1];
  }
  function minMaxScoreNormalization(results) {
    const maxScore = Math.max.apply(Math, results.map(extractScore));
    return results.map(([id, score]) => [id, score / maxScore]);
  }
  function normalizeScore(score, maxScore) {
    return score / maxScore;
  }
  function hybridScoreBuilder(textWeight, vectorWeight) {
    return (textScore, vectorScore) => textScore * textWeight + vectorScore * vectorWeight;
  }
  function mergeAndRankResults(textResults, vectorResults, query, hybridWeights) {
    const maxTextScore = Math.max.apply(Math, textResults.map(extractScore));
    const maxVectorScore = Math.max.apply(Math, vectorResults.map(extractScore));
    const hasHybridWeights = hybridWeights && hybridWeights.text && hybridWeights.vector;
    const { text: textWeight, vector: vectorWeight } = hasHybridWeights ? hybridWeights : getQueryWeights(query);
    const mergedResults = /* @__PURE__ */ new Map();
    const textResultsLength = textResults.length;
    const hybridScore = hybridScoreBuilder(textWeight, vectorWeight);
    for (let i = 0; i < textResultsLength; i++) {
      const [id, score] = textResults[i];
      const normalizedScore = normalizeScore(score, maxTextScore);
      const hybridScoreValue = hybridScore(normalizedScore, 0);
      mergedResults.set(id, hybridScoreValue);
    }
    const vectorResultsLength = vectorResults.length;
    for (let i = 0; i < vectorResultsLength; i++) {
      const [resultId, score] = vectorResults[i];
      const normalizedScore = normalizeScore(score, maxVectorScore);
      const existingRes = mergedResults.get(resultId) ?? 0;
      mergedResults.set(resultId, existingRes + hybridScore(0, normalizedScore));
    }
    return [...mergedResults].sort((a, b) => b[1] - a[1]);
  }
  function getQueryWeights(query) {
    return {
      text: 0.5,
      vector: 0.5
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/search.js
  function search3(orama, params, language) {
    const mode = params.mode ?? MODE_FULLTEXT_SEARCH;
    if (mode === MODE_FULLTEXT_SEARCH) {
      return fullTextSearch(orama, params, language);
    }
    if (mode === MODE_VECTOR_SEARCH) {
      return searchVector(orama, params);
    }
    if (mode === MODE_HYBRID_SEARCH) {
      return hybridSearch(orama, params);
    }
    throw createError("INVALID_SEARCH_MODE", mode);
  }
  function fetchDocumentsWithDistinct(orama, uniqueDocsArray, offset2, limit, distinctOn) {
    const docs = orama.data.docs;
    const values = /* @__PURE__ */ new Map();
    const results = [];
    const resultIDs = /* @__PURE__ */ new Set();
    const uniqueDocsArrayLength = uniqueDocsArray.length;
    let count3 = 0;
    for (let i = 0; i < uniqueDocsArrayLength; i++) {
      const idAndScore = uniqueDocsArray[i];
      if (typeof idAndScore === "undefined") {
        continue;
      }
      const [id, score] = idAndScore;
      if (resultIDs.has(id)) {
        continue;
      }
      const doc = orama.documentsStore.get(docs, id);
      const value = getNested(doc, distinctOn);
      if (typeof value === "undefined" || values.has(value)) {
        continue;
      }
      values.set(value, true);
      count3++;
      if (count3 <= offset2) {
        continue;
      }
      results.push({ id: getDocumentIdFromInternalId(orama.internalDocumentIDStore, id), score, document: doc });
      resultIDs.add(id);
      if (count3 >= offset2 + limit) {
        break;
      }
    }
    return results;
  }
  function fetchDocuments(orama, uniqueDocsArray, offset2, limit) {
    const docs = orama.data.docs;
    const results = Array.from({
      length: limit
    });
    const resultIDs = /* @__PURE__ */ new Set();
    for (let i = offset2; i < limit + offset2; i++) {
      const idAndScore = uniqueDocsArray[i];
      if (typeof idAndScore === "undefined") {
        break;
      }
      const [id, score] = idAndScore;
      if (!resultIDs.has(id)) {
        const fullDoc = orama.documentsStore.get(docs, id);
        results[i] = { id: getDocumentIdFromInternalId(orama.internalDocumentIDStore, id), score, document: fullDoc };
        resultIDs.add(id);
      }
    }
    return results;
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/methods/serialization.js
  function save6(orama) {
    return {
      internalDocumentIDStore: orama.internalDocumentIDStore.save(orama.internalDocumentIDStore),
      index: orama.index.save(orama.data.index),
      docs: orama.documentsStore.save(orama.data.docs),
      sorting: orama.sorter.save(orama.data.sorting),
      pinning: orama.pinning.save(orama.data.pinning),
      language: orama.tokenizer.language
    };
  }

  // ../../../ui/node_modules/.pnpm/@orama+orama@3.1.18/node_modules/@orama/orama/dist/esm/types.js
  var kInsertions = Symbol("orama.insertions");
  var kRemovals = Symbol("orama.removals");

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/advanced-BnqbPavQ.js
  async function searchAdvanced(db, query, tag = [], { mode = "fulltext", ...override } = {}) {
    if (typeof tag === "string") tag = [tag];
    const params = {
      limit: 60,
      mode,
      ...override,
      where: removeUndefined({
        tags: tag.length > 0 ? { containsAll: tag } : void 0,
        ...override.where
      }),
      groupBy: {
        properties: ["page_id"],
        maxResult: 8,
        ...override.groupBy
      },
      properties: mode === "fulltext" ? ["content"] : ["content", "embeddings"]
    };
    if (query.length > 0) params.term = query;
    const highlighter = createContentHighlighter(query);
    const result = await search3(db, params);
    const list4 = [];
    for (const item of result.groups ?? []) {
      const pageId = item.values[0];
      const page = getByID(db, pageId);
      if (!page) continue;
      list4.push({
        id: pageId,
        type: "page",
        content: highlighter.highlightMarkdown(page.content),
        breadcrumbs: page.breadcrumbs,
        url: page.url
      });
      for (const hit of item.result) {
        if (hit.document.type === "page") continue;
        list4.push({
          id: hit.document.id.toString(),
          content: highlighter.highlightMarkdown(hit.document.content),
          breadcrumbs: hit.document.breadcrumbs,
          type: hit.document.type,
          url: hit.document.url
        });
      }
    }
    if (typeof params.limit === "number" && list4.length > params.limit) return list4.slice(0, params.limit);
    return list4;
  }

  // ../../../ui/node_modules/.pnpm/fumadocs-core@16.10.5_@mdx-js+mdx@3.1.1_@types+estree-jsx@1.0.5_@types+hast@3.0.4_@type_45158c13ab337377714ea1553413c41d/node_modules/fumadocs-core/dist/search/server.js
  var advancedSchema = {
    content: "string",
    page_id: "string",
    type: "string",
    breadcrumbs: "string[]",
    tags: "enum[]",
    url: "string",
    embeddings: "vector[512]"
  };
  async function createDB({ indexes, tokenizer: tokenizer3, search: _, ...rest }) {
    const items = typeof indexes === "function" ? await indexes() : indexes;
    const db = create5({
      schema: advancedSchema,
      ...rest,
      components: {
        ...rest.components,
        tokenizer: tokenizer3 ?? rest.components?.tokenizer
      }
    });
    await insertMultiple(db, buildDocuments(items));
    return db;
  }
  function initAdvancedSearch(options) {
    const get2 = createDB(options);
    return {
      async export() {
        return {
          type: "advanced",
          ...save6(await get2)
        };
      },
      async search(query, searchOptions = {}) {
        const db = await get2;
        const { limit, tag, mode } = searchOptions;
        return searchAdvanced(db, query, tag, {
          ...options.search,
          limit,
          mode: mode === "vector" ? "vector" : "fulltext"
        }).catch((err) => {
          if (mode === "vector") throw new Error("failed to search, make sure you have installed `@orama/plugin-embeddings` according to their docs.", { cause: err });
          throw err;
        });
      }
    };
  }

  // entry.mjs
  var processor2 = unified().use(remarkMdx).use(remarkGfm).use(remarkStructure).freeze();
  var server = null;
  globalThis.DocsSearch = {
    // structure(mdast JSON) -> structuredData JSON.
    structure(json) {
      const tree = JSON.parse(json);
      const file = new VFile();
      processor2.runSync(tree, file);
      return JSON.stringify(file.data.structuredData);
    },
    // index(indexes JSON): [{ id, title, description, url, breadcrumbs, structuredData }].
    index(json) {
      server = initAdvancedSearch({ indexes: JSON.parse(json) });
    },
    // search(query) -> Promise<JSON>: the /api/search response body.
    async search(query) {
      if (!server) return "[]";
      return JSON.stringify(await server.search(query, {}));
    }
  };
})();
