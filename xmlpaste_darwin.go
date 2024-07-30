package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework AppKit -framework Foundation

#include <string.h>
#import <AppKit/NSPasteboard.h>
#import <Foundation/Foundation.h>

const char *paste() {
	const char *str = "";
	NSPasteboard* myPasteboard = [NSPasteboard generalPasteboard];
	for (NSString* type in [myPasteboard types]) {
		if ([type hasPrefix:@"dyn.ah62d4rv4gk8zuxn"] || [type isEqualToString:@"dyn.agk8u"] || [type isEqualToString:@"public.utf16-plain-text"]) {
			// dyn.ah62d4rv4gk8zuxnykk: Table
			// dyn.ah62d4rv4gk8zuxngku: Field
			// dyn.ah62d4rv4gk8zuxnxkq: Script
			// dyn.ah62d4rv4gk8zuxnxnq: Script Step
			// dyn.ah62d4rv4gk8zuxngm2: Custom Function
			// dyn.ah62d4rv4gk8zuxnqm6: Layout Object (.fp7)
			// dyn.ah62d4rv4gk8zuxnqgk: Layout Object (.fmp12)
			// dyn.ah62d4rv4gk8zuxn0mu: Value List
			// dyn.agk8u              : Theme (FileMaker 17, 18, 19 and 2023)
			// dyn.ah62d4rv4gk8zuxnyma: Theme (FileMaker 2024)
			// public.utf16-plain-text: Custom Menu
			NSData* data = [myPasteboard dataForType:type];
			if (data) {
				NSString *string = @"";
				if ([type isEqualToString:@"public.utf16-plain-text"]) {
					// for custom menus
					string = [[NSString alloc] initWithData:data encoding:NSUTF16LittleEndianStringEncoding];
					if (!([string hasPrefix:@"<?xml version=\"1.0\" encoding=\"utf-16\"?><FMObjectTransfer "])) {
						break;
					}
				} else {
					string = [[NSString alloc] initWithData:data encoding:NSUTF8StringEncoding];
				}
				str = [string UTF8String];
				break;
			}
		}
	}

	return str;
}
*/
import "C"

func getClipboard() (string, error) {
	content := C.GoString(C.paste())
	return content, nil
}
