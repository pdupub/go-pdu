//
//  PDUMApp.swift
//  PDUM
//
//  Created by Liu Peng on 2021/8/16.
//

import SwiftUI

@main
struct PDUMApp: App {
    @StateObject private var modelData = ModelData()
    
    var body: some Scene {
        WindowGroup {
            ContentView().environmentObject(modelData)
        }
    }
}
