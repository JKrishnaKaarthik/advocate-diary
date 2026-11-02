package com.advocatediary.tests;

import com.advocatediary.pages.InboxPage;
import org.testng.Assert;
import org.testng.annotations.Test;

public class InboxShowsHearingTest extends BaseTest {

  @Test
  public void nextHearingShowsInInbox() {
    signInToDiary();
    InboxPage inbox = completeClientCaseOnDiary().nav().openInbox();
    inbox.waitForNextHearing(CLIENT_NAME, CASE_TITLE, CASE_NUMBER);
    Assert.assertTrue(inboxShowsMatter(inbox), "next hearing should show in Inbox");
  }
}
