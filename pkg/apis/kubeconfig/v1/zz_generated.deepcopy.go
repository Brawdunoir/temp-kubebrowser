//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright Yann Lacroix.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kubeconfig) DeepCopyInto(out *Kubeconfig) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kubeconfig.
func (in *Kubeconfig) DeepCopy() *Kubeconfig {
	if in == nil {
		return nil
	}
	out := new(Kubeconfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeconfigList) DeepCopyInto(out *KubeconfigList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Kubeconfig, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeconfigList.
func (in *KubeconfigList) DeepCopy() *KubeconfigList {
	if in == nil {
		return nil
	}
	out := new(KubeconfigList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeconfigSpec) DeepCopyInto(out *KubeconfigSpec) {
	*out = *in
	in.Kubeconfig.DeepCopyInto(&out.Kubeconfig)
	if in.Whitelist != nil {
		in, out := &in.Whitelist, &out.Whitelist
		*out = new(Whitelist)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeconfigSpec.
func (in *KubeconfigSpec) DeepCopy() *KubeconfigSpec {
	if in == nil {
		return nil
	}
	out := new(KubeconfigSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Whitelist) DeepCopyInto(out *Whitelist) {
	*out = *in
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Whitelist.
func (in *Whitelist) DeepCopy() *Whitelist {
	if in == nil {
		return nil
	}
	out := new(Whitelist)
	in.DeepCopyInto(out)
	return out
}
