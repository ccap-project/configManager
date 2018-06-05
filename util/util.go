/*
 *
 * Copyright (c) 2017, 2018 Alexandre Biancalana <ale@biancalanas.net>.
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *     * Redistributions of source code must retain the above copyright
 *       notice, this list of conditions and the following disclaimer.
 *     * Redistributions in binary form must reproduce the above copyright
 *       notice, this list of conditions and the following disclaimer in the
 *       documentation and/or other materials provided with the distribution.
 *     * Neither the name of the <organization> nor the
 *       names of its contributors may be used to endorse or promote products
 *       derived from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL <COPYRIGHT HOLDER> BE LIABLE FOR ANY
 * DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package util

import (
	"fmt"
	"reflect"
	"strings"

	"configManager/models"

	"github.com/go-openapi/swag"

	"github.com/fatih/structs"
)

func BuildQuery(Struct interface{}, prefix string, queryType string, skip []string) string {

	var res []string

	f := structs.Fields(Struct)

	for _, fi := range f {
		if fi.IsZero() == false {
			if grep(fi.Name(), skip) {
				continue
			}
			name := swag.ToFileName(fi.Name())

			switch queryType {
			case "update":
				if len(prefix) == 0 {
					res = append(res, fmt.Sprintf("%s={%s} ", name, name))
				} else {
					res = append(res, fmt.Sprintf("%s.%s={%s} ", prefix, name, name))
				}
			default:
				if len(prefix) == 0 {
					res = append(res, fmt.Sprintf("%s:{%s}", name, name))
				} else {
					res = append(res, fmt.Sprintf("%s:{%s_%s}", name, prefix, name))
				}
			}
		}
	}

	return strings.Join(res, ", ")
}

func BuildParams(Struct interface{}, prefix string, Vals map[string]interface{}, skip []string) map[string]interface{} {

	f := structs.Fields(Struct)

	for _, fi := range f {
		if fi.IsZero() == false {

			if grep(fi.Name(), skip) {
				continue
			}
			var k string

			// Add prefix string if needed, convert field name to snake case
			if len(prefix) == 0 {
				k = swag.ToFileName(fi.Name())
			} else {
				k = fmt.Sprintf("%s_%s", prefix, swag.ToFileName(fi.Name()))
			}

			v := fi.Value()

			if strings.HasPrefix(reflect.TypeOf(v).String(), "*") {
				val := reflect.ValueOf(v)
				Vals[k] = val.Elem().Interface()
			} else {
				Vals[k] = v
			}
		}
	}
	return Vals
}

func grep(s string, list []string) bool {

	for _, k := range list {
		if s == k {
			return true
		}
	}
	return false

}

func SetField(obj interface{}, name string, value interface{}) error {

	structValue := reflect.ValueOf(obj).Elem()
	structFieldValue := structValue.FieldByName(name)

	if !structFieldValue.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !structFieldValue.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	structFieldType := structFieldValue.Type()

	var val reflect.Value

	if strings.HasPrefix(structFieldType.String(), "*") {
		val = reflect.New(reflect.TypeOf(value))
		val.Elem().Set(reflect.ValueOf(value))

	} else if structFieldType.String() == "models.ULID" {
		val = reflect.ValueOf(models.ULID(value.(string)))

	} else {
		val = reflect.ValueOf(value)
	}

	structFieldValue.Set(val)
	return nil
}

func FillStruct(s interface{}, m map[string]interface{}) error {
	for k, v := range m {
		err := SetField(s, swag.ToGoName(k), v)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
	return nil
}
