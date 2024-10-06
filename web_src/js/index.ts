// bootstrap module must be the first one to be imported, it handles webpack lazy-loading and global errors
import './bootstrap.ts';
import './htmx.ts';

import {initDashboardRepoList} from './components/DashboardRepoList.vue';

import {initGlobalCopyToClipboardListener} from './features/clipboard.ts';
import {initContextPopups} from './features/contextpopup.ts';
import {initHeatmap} from './features/heatmap.ts';
import {initImageDiff} from './features/imagediff.ts';
import {initTableSort} from './features/tablesort.ts';
import {initAutoFocusEnd} from './features/autofocus-end.ts';
import {initAdminUserListSearchForm} from './features/admin/users.ts';
import {initAdminConfigs} from './features/admin/config.ts';
import {initNotificationCount, initNotificationsTable} from './features/notification.ts';
import {initStopwatch} from './features/stopwatch.ts';
import {initPdfViewer} from './render/pdf.ts';

import {initUserAuthOauth2} from './features/user-auth.ts';
import {initAdminEmails} from './features/admin/emails.ts';
import {initAdminCommon} from './features/admin/common.ts';
import {initUserSettings} from './features/user-settings.ts';
import {initUserAuthWebAuthn, initUserAuthWebAuthnRegister} from './features/user-auth-webauthn.ts';
import {initCompSearchUserBox} from './features/comp/SearchUserBox.ts';
import {initInstall} from './features/install.ts';
import {initCompWebHookEditor} from './features/comp/WebHookEditor.ts';
import {initCopyContent} from './features/copycontent.ts';
import {initCaptcha} from './features/captcha.ts';
import {initGlobalTooltips} from './modules/tippy.ts';
import {initGiteaFomantic} from './modules/fomantic.ts';
import {initSubmitEventPolyfill, onDomReady} from './utils/dom.ts';
import {initDirAuto} from './modules/dirauto.ts';
import {initColorPickers} from './features/colorpicker.ts';
import {initAdminSelfCheck} from './features/admin/selfcheck.ts';
import {initOAuth2SettingsDisableCheckbox} from './features/oauth2-settings.ts';
import {initGlobalFetchAction} from './features/common-fetch-action.ts';
import {initScopedAccessTokenCategories} from './features/scoped-access-token.ts';
import {
  initFootLanguageMenu,
  initGlobalDropdown,
  initGlobalTabularMenu,
  initHeadNavbarContentToggle,
} from './features/common-page.ts';
import {
  initGlobalButtonClickOnEnter,
  initGlobalButtons,
  initGlobalDeleteButton,
  initGlobalShowModal,
} from './features/common-button.ts';
import {initGlobalEnterQuickSubmit, initGlobalFormDirtyLeaveConfirm} from './features/common-form.ts';

initGiteaFomantic();
initDirAuto();
initSubmitEventPolyfill();

function callInitFunctions(functions: (() => any)[]) {
  // Start performance trace by accessing a URL by "https://localhost/?_ui_performance_trace=1" or "https://localhost/?key=value&_ui_performance_trace=1"
  // It is a quick check, no side effect so no need to do slow URL parsing.
  const initStart = performance.now();
  if (window.location.search.includes('_ui_performance_trace=1')) {
    let results: {name: string, dur: number}[] = [];
    for (const func of functions) {
      const start = performance.now();
      func();
      results.push({name: func.name, dur: performance.now() - start});
    }
    results = results.sort((a, b) => b.dur - a.dur);
    for (let i = 0; i < 20 && i < results.length; i++) {
      // eslint-disable-next-line no-console
      console.log(`performance trace: ${results[i].name} ${results[i].dur.toFixed(3)}`);
    }
  } else {
    for (const func of functions) {
      func();
    }
  }
  const initDur = performance.now() - initStart;
  if (initDur > 500) {
    console.error(`slow init functions took ${initDur.toFixed(3)}ms`);
  }
}

onDomReady(() => {
  callInitFunctions([
    initGlobalDropdown,
    initGlobalTabularMenu,
    initGlobalShowModal,
    initGlobalFetchAction,
    initGlobalTooltips,
    initGlobalButtonClickOnEnter,
    initGlobalButtons,
    initGlobalCopyToClipboardListener,
    initGlobalEnterQuickSubmit,
    initGlobalFormDirtyLeaveConfirm,
    initGlobalDeleteButton,

    initCompSearchUserBox,
    initCompWebHookEditor,

    initInstall,

    initHeadNavbarContentToggle,
    initFootLanguageMenu,

    initContextPopups,
    initHeatmap,
    initImageDiff,
    initStopwatch,
    initTableSort,
    initAutoFocusEnd,
    initCopyContent,

    initAdminCommon,
    initAdminEmails,
    initAdminUserListSearchForm,
    initAdminConfigs,
    initAdminSelfCheck,

    initDashboardRepoList,

    initNotificationCount,
    initNotificationsTable,

    initCaptcha,

    initUserAuthOauth2,
    initUserAuthWebAuthn,
    initUserAuthWebAuthnRegister,
    initUserSettings,
    initPdfViewer,
    initScopedAccessTokenCategories,
    initColorPickers,

    initOAuth2SettingsDisableCheckbox,
  ]);
});
